package shared_utilities

import (
	"net/http"
	"strings"
)

// NormalizeRoutesMiddleware ensures all routes are normalized to lowercase
func NormalizeRoutesMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		normalizedPath := strings.ToLower(r.URL.Path)
		r.URL.Path = normalizedPath
		next.ServeHTTP(w, r)
	})
}
