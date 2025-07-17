package shared_utilities

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	"github.com/hstles/go-sdk/core_config"
)

// SecurityMiddleware returns a CORS + CSP middleware that applies to any https://*.{parentDomain}.
func SecurityMiddleware(parentDomain string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// CORS: allow any https://*.parentDomain
			if origin := r.Header.Get("Origin"); origin != "" {
				if u, err := url.Parse(origin); err == nil && u.Scheme == "https" {
					host := u.Hostname()
					if host == parentDomain || strings.HasSuffix(host, "."+parentDomain) {
						w.Header().Set("Access-Control-Allow-Origin", origin)
						w.Header().Set("Access-Control-Allow-Credentials", "true")
						w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
						w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
					}
				}
			}

			// CSP: subdomains, CDNs, fonts, inline & eval
			csp := fmt.Sprintf(
				"default-src 'self' https://*.%s; "+
					"img-src 'self' data:; "+
					"font-src 'self' https://cdnjs.cloudflare.com https://fonts.gstatic.com; "+
					"script-src 'self' https://*.%s https://unpkg.com https://fonts.googleapis.com https://cdnjs.cloudflare.com 'unsafe-inline' 'unsafe-eval'; "+
					"style-src 'self' https://*.%s https://fonts.googleapis.com https://cdnjs.cloudflare.com 'unsafe-inline';",
				parentDomain, parentDomain, parentDomain,
			)
			w.Header().Set("Content-Security-Policy", csp)

			// Preflight
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// SecurityWithProviderValidationMiddleware combines security headers with provider validation
// This is intended to be used when you want both security headers and provider validation in one middleware
func SecurityWithProviderValidationMiddleware(parentDomain string, coreCfg *core_config.CoreConfig) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Apply security headers first
			securityMiddleware := SecurityMiddleware(parentDomain)
			securityMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Then apply provider validation if session data is available
				if sessionData, ok := GetSessionDataFromContext(r.Context()); ok && sessionData.Provider != "" {
					if err := ValidateSessionProvider(sessionData.Provider, coreCfg.AppName, coreCfg.AppDomains); err != nil {
						http.Error(w, fmt.Sprintf("Access denied: %v", err), http.StatusForbidden)
						return
					}
				}
				next.ServeHTTP(w, r)
			})).ServeHTTP(w, r)
		})
	}
}
