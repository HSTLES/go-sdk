package core_utilities

import "net/http"

// UserSessionData stores user session information
type UserSessionData struct {
	UserID string
}

// Context key for session
type contextKey string

const sessionContextKey = contextKey("session")

// GetSessionContextKey returns the context key for session data
func GetSessionContextKey() contextKey {
	return sessionContextKey
}

// GetUserIDFromContext retrieves the user ID from the request context
// This should be used in handlers that are protected by session validation middleware
func GetUserIDFromContext(r *http.Request) (userID string, ok bool) {
	data, ok := r.Context().Value(sessionContextKey).(*UserSessionData)
	if !ok || data.UserID == "" {
		return "", false
	}
	return data.UserID, true
}
