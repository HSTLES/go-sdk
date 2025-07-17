package core_utilities

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
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
