package shared_utilities

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/hstles/go-sdk/client_auth"
)

// SecurityConfig holds configuration for route-based security
type SecurityConfig struct {
	AuthServiceURL     string
	IdentityServiceURL string
	NotifyServiceURL   string
	AccountServiceURL  string
	LoginServiceURL    string
	APIKeys            map[string]string // service_name -> api_key
}

// LoadSecurityConfig loads security configuration from environment
func LoadSecurityConfig() *SecurityConfig {
	apiKeys := make(map[string]string)

	// Load API keys from environment variables
	if key := os.Getenv("AUTH_SERVICE_API_KEY"); key != "" {
		apiKeys["auth"] = key
	}
	if key := os.Getenv("IDENTITY_SERVICE_API_KEY"); key != "" {
		apiKeys["identity"] = key
	}
	if key := os.Getenv("NOTIFY_SERVICE_API_KEY"); key != "" {
		apiKeys["notify"] = key
	}
	if key := os.Getenv("ACCOUNT_SERVICE_API_KEY"); key != "" {
		apiKeys["account"] = key
	}

	// Load service URLs from environment variables with defaults
	authURL := os.Getenv("AUTH_SERVICE_URL")
	if authURL == "" {
		authURL = "https://auth.hstles.com"
	}

	identityURL := os.Getenv("IDENTITY_SERVICE_URL")
	if identityURL == "" {
		identityURL = "https://identity.hstles.com"
	}

	notifyURL := os.Getenv("NOTIFY_SERVICE_URL")
	if notifyURL == "" {
		notifyURL = "https://notify.hstles.com"
	}

	accountURL := os.Getenv("ACCOUNT_SERVICE_URL")
	if accountURL == "" {
		accountURL = "https://account.hstles.com"
	}

	loginURL := os.Getenv("LOGIN_SERVICE_URL")
	if loginURL == "" {
		loginURL = "https://login.hstles.com"
	}

	return &SecurityConfig{
		AuthServiceURL:     authURL,
		IdentityServiceURL: identityURL,
		NotifyServiceURL:   notifyURL,
		AccountServiceURL:  accountURL,
		LoginServiceURL:    loginURL,
		APIKeys:            apiKeys,
	}
}

// SecurityRoutes helps organize routes by security level
type SecurityRoutes struct {
	Router *mux.Router
	config *SecurityConfig

	// Route groups
	Public    *mux.Router // No authentication required
	Protected *mux.Router // Requires user session
	Service   *mux.Router // Requires API key (service-to-service)
	Mixed     *mux.Router // Custom security handling per endpoint
}

// NewSecurityRoutes creates a new security-aware router setup
func NewSecurityRoutes(parentDomain string, config *SecurityConfig) *SecurityRoutes {
	r := mux.NewRouter()
	r.Use(SecurityMiddleware(parentDomain))

	return &SecurityRoutes{
		Router:    r,
		config:    config,
		Public:    r.PathPrefix("/").Subrouter(),    // No middleware
		Protected: createProtectedRoutes(r, config), // Session middleware
		Service:   createServiceRoutes(r, config),   // API key middleware
		Mixed:     r.PathPrefix("/").Subrouter(),    // Custom per-endpoint
	}
}

// createProtectedRoutes creates a subrouter with session validation
func createProtectedRoutes(parent *mux.Router, config *SecurityConfig) *mux.Router {
	protected := parent.PathPrefix("/").Subrouter()
	protected.Use(SessionValidationMiddleware(config.AuthServiceURL))
	return protected
}

// createServiceRoutes creates a subrouter with API key validation
func createServiceRoutes(parent *mux.Router, config *SecurityConfig) *mux.Router {
	service := parent.PathPrefix("/").Subrouter()
	service.Use(ServiceAuthMiddleware(config.APIKeys))
	return service
}

// SessionValidationMiddleware validates sessions using the auth service
func SessionValidationMiddleware(authServiceURL string) mux.MiddlewareFunc {
	authClient := client_auth.NewClient(authServiceURL)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Validate session with auth service
			resp, code, err := authClient.ValidateSession(r.Context(), r.Cookies())
			if err != nil {
				log.Printf("Session validation error: %v", err)
				http.Error(w, "Session validation failed", code)
				return
			}

			if !resp.Valid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Store session data in context
			sessionData := UserSessionData{
				UserID:   resp.UserID,
				Provider: resp.Provider,
			}
			ctx := SetSessionDataInContext(r.Context(), sessionData)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ServiceAuthMiddleware validates API keys for service-to-service communication
func ServiceAuthMiddleware(validAPIKeys map[string]string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("X-API-Key")
			if apiKey == "" {
				http.Error(w, "API key required", http.StatusUnauthorized)
				return
			}

			// Validate API key
			var serviceName string
			for name, key := range validAPIKeys {
				if key == apiKey {
					serviceName = name
					break
				}
			}

			if serviceName == "" {
				log.Printf("Invalid API key attempt from %s", r.RemoteAddr)
				http.Error(w, "Invalid API key", http.StatusUnauthorized)
				return
			}

			// Store service info in context
			ctx := context.WithValue(r.Context(), "service_name", serviceName)
			log.Printf("Service-to-service call from: %s", serviceName)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// MixedAuthMiddleware allows both session and API key authentication
func MixedAuthMiddleware(authServiceURL string, validAPIKeys map[string]string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Try API key first
			if apiKey := r.Header.Get("X-API-Key"); apiKey != "" {
				for name, key := range validAPIKeys {
					if key == apiKey {
						ctx := context.WithValue(r.Context(), "service_name", name)
						ctx = context.WithValue(ctx, "auth_type", "service")
						next.ServeHTTP(w, r.WithContext(ctx))
						return
					}
				}
				http.Error(w, "Invalid API key", http.StatusUnauthorized)
				return
			}

			// Try session authentication
			authClient := client_auth.NewClient(authServiceURL)
			resp, _, err := authClient.ValidateSession(r.Context(), r.Cookies())
			if err == nil && resp.Valid {
				sessionData := UserSessionData{
					UserID:   resp.UserID,
					Provider: resp.Provider,
				}
				ctx := SetSessionDataInContext(r.Context(), sessionData)
				ctx = context.WithValue(ctx, "auth_type", "session")
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		})
	}
}

// RequireSessionUser is a helper to get user ID from validated session
func RequireSessionUser(r *http.Request) (string, error) {
	sessionData, ok := GetSessionDataFromContext(r.Context())
	if !ok {
		return "", fmt.Errorf("no session data in context")
	}
	if sessionData.UserID == "" {
		return "", fmt.Errorf("no user ID in session")
	}
	return sessionData.UserID, nil
}

// GetServiceName gets the service name from API key authentication
func GetServiceName(r *http.Request) (string, bool) {
	serviceName, ok := r.Context().Value("service_name").(string)
	return serviceName, ok
}

// GetAuthType returns the authentication type used (session or service)
func GetAuthType(r *http.Request) string {
	authType, ok := r.Context().Value("auth_type").(string)
	if !ok {
		return "unknown"
	}
	return authType
}

// GetAuthClient returns a configured auth service client
func (c *SecurityConfig) GetAuthClient() *client_auth.Client {
	return client_auth.NewClient(c.AuthServiceURL)
}

// GetServiceURL returns the URL for a specific service
func (c *SecurityConfig) GetServiceURL(serviceName string) string {
	switch serviceName {
	case "auth":
		return c.AuthServiceURL
	case "identity":
		return c.IdentityServiceURL
	case "notify":
		return c.NotifyServiceURL
	case "account":
		return c.AccountServiceURL
	case "login":
		return c.LoginServiceURL
	default:
		return ""
	}
}

// GetServiceAPIKey returns the API key for a specific service
func (c *SecurityConfig) GetServiceAPIKey(serviceName string) string {
	return c.APIKeys[serviceName]
}
