package shared_utilities

import (
	"context"
	"net/http"
)

// UserSessionData stores user session information
type UserSessionData struct {
	UserID   string
	Provider string
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
	data, ok := r.Context().Value(sessionContextKey).(UserSessionData)
	if !ok || data.UserID == "" {
		return "", false
	}
	return data.UserID, true
}

// GetProviderFromContext retrieves the provider from the request context
func GetProviderFromContext(ctx context.Context) string {
	data, ok := ctx.Value(sessionContextKey).(UserSessionData)
	if !ok {
		return ""
	}
	return data.Provider
}

// GetSessionDataFromContext retrieves the full session data from context
func GetSessionDataFromContext(ctx context.Context) (UserSessionData, bool) {
	data, ok := ctx.Value(sessionContextKey).(UserSessionData)
	return data, ok
}

// SetSessionDataInContext stores session data in context
func SetSessionDataInContext(ctx context.Context, sessionData UserSessionData) context.Context {
	return context.WithValue(ctx, sessionContextKey, sessionData)
}
