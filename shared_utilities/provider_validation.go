package shared_utilities

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hstles/go-sdk/core_config"
)

// ProviderValidationError represents an error when provider validation fails
type ProviderValidationError struct {
	Message          string
	SessionProvider  string
	AllowedProviders []string
	AppName          string
}

func (e ProviderValidationError) Error() string {
	return e.Message
}

// ValidateSessionProvider checks if the session provider is allowed for the given app
// It compares the session provider against the allowed providers in core_config for the app
func ValidateSessionProvider(sessionProvider, appName string, appDomains string) error {
	if sessionProvider == "" {
		return &ProviderValidationError{
			Message:          "Session provider is empty",
			SessionProvider:  sessionProvider,
			AllowedProviders: []string{},
			AppName:          appName,
		}
	}

	// Get app config by name first
	appConfig := core_config.GetAppByName(appName)

	// If not found by exact name, try alternative approaches
	if appConfig.AppName == "" {
		// Try to find by route that matches the app name
		appConfig = core_config.GetAppByRoute("/" + appName)

		// If still not found and appDomains is provided, try to extract from domains
		if appConfig.AppName == "" && appDomains != "" {
			domains := strings.Split(appDomains, ",")
			for _, domain := range domains {
				domain = strings.TrimSpace(domain)
				// Extract app name from domain (e.g., "account.hstles.com" -> "account")
				if strings.Contains(domain, ".") {
					potentialAppName := strings.Split(domain, ".")[0]
					appConfig = core_config.GetAppByRoute("/" + potentialAppName)
					if appConfig.AppName != "" {
						break
					}
				}
			}
		}

		// If still not found, try a more flexible search
		if appConfig.AppName == "" {
			for _, config := range core_config.AppConfigs {
				// Check if the app name is contained in the config name (case insensitive)
				if strings.Contains(strings.ToLower(config.AppName), strings.ToLower(appName)) ||
					strings.Contains(strings.ToLower(appName), strings.ToLower(config.AppName)) {
					appConfig = config
					break
				}
				// Check if Domain contains the app name
				if strings.Contains(config.Domain, appName+".") {
					appConfig = config
					break
				}
			}
		}

		// If still not found, default to a permissive config or return error
		if appConfig.AppName == "" {
			return &ProviderValidationError{
				Message:          fmt.Sprintf("App configuration not found for app: %s", appName),
				SessionProvider:  sessionProvider,
				AllowedProviders: []string{},
				AppName:          appName,
			}
		}
	}

	// Check if the session provider is in the allowed list
	for _, allowedProvider := range appConfig.AuthMethods {
		if sessionProvider == allowedProvider {
			return nil // Provider is allowed
		}
	}

	// Provider not found in allowed list
	return &ProviderValidationError{
		Message: fmt.Sprintf(
			"Provider '%s' is not allowed for app '%s' (%s). Allowed providers: %s",
			sessionProvider,
			appName,
			appConfig.DisplayName,
			strings.Join(appConfig.AuthMethods, ", "),
		),
		SessionProvider:  sessionProvider,
		AllowedProviders: appConfig.AuthMethods,
		AppName:          appName,
	}
}

// ProviderValidationMiddleware creates middleware that validates session providers
// This should be used after session validation but before authorization
func ProviderValidationMiddleware(appName, appDomains string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get session data from context (should be set by previous middleware)
			sessionData, ok := r.Context().Value(GetSessionContextKey()).(UserSessionData)
			if !ok {
				http.Error(w, "Session data not found in context", http.StatusInternalServerError)
				return
			}

			// Validate provider
			provider := sessionData.Provider
			if provider == "" {
				http.Error(w, "Provider information not available", http.StatusUnauthorized)
				return
			}

			if err := ValidateSessionProvider(provider, appName, appDomains); err != nil {
				if pvErr, ok := err.(*ProviderValidationError); ok {
					http.Error(w, pvErr.Message, http.StatusForbidden)
				} else {
					http.Error(w, "Provider validation failed", http.StatusInternalServerError)
				}
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// IsProviderAllowed is a simple utility function to check if a provider is allowed for an app
// Returns true if allowed, false otherwise
func IsProviderAllowed(provider, appName string) bool {
	return ValidateSessionProvider(provider, appName, "") == nil
}

// GetAllowedProvidersForApp returns the list of allowed providers for a given app
func GetAllowedProvidersForApp(appName string) []string {
	// Try multiple approaches to find the app config
	appConfig := core_config.GetAppByName(appName)
	if appConfig.AppName == "" {
		appConfig = core_config.GetAppByRoute("/" + appName)
	}
	if appConfig.AppName == "" {
		// Try flexible search
		for _, config := range core_config.AppConfigs {
			if strings.Contains(strings.ToLower(config.AppName), strings.ToLower(appName)) ||
				strings.Contains(strings.ToLower(appName), strings.ToLower(config.AppName)) ||
				strings.Contains(config.Domain, appName+".") {
				appConfig = config
				break
			}
		}
	}

	return appConfig.AuthMethods
}

// ValidateProviderFromContext validates the provider from the request context
// This is a convenience function for handlers that have session data in context
func ValidateProviderFromContext(r *http.Request, appName, appDomains string) error {
	sessionData, ok := GetSessionDataFromContext(r.Context())
	if !ok {
		return &ProviderValidationError{
			Message:          "Session data not found in context",
			SessionProvider:  "",
			AllowedProviders: []string{},
			AppName:          appName,
		}
	}

	return ValidateSessionProvider(sessionData.Provider, appName, appDomains)
}
