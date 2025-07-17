package shared_helpers

import (
	"net/http"
	"strings"
)

// GetClientIP extracts the client IP address from the HTTP request
// It checks X-Forwarded-For, X-Real-IP headers and falls back to RemoteAddr
func GetClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (most common with load balancers/proxies)
	clientIP := r.Header.Get("X-Forwarded-For")
	if clientIP != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		clientIP = strings.TrimSpace(strings.Split(clientIP, ",")[0])
		if clientIP != "" {
			return clientIP
		}
	}

	// Check X-Real-IP header (used by some proxies)
	clientIP = r.Header.Get("X-Real-IP")
	if clientIP != "" {
		return strings.TrimSpace(clientIP)
	}

	// Fall back to RemoteAddr (remove port if present)
	if r.RemoteAddr != "" {
		return strings.Split(r.RemoteAddr, ":")[0]
	}

	return ""
}

// GetUserAgent extracts the User-Agent header from the HTTP request
func GetUserAgent(r *http.Request) string {
	return r.Header.Get("User-Agent")
}
