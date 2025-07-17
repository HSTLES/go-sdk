package client_auth

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// ValidateSessionHandler proxies GET /api/session.
func ValidateSessionHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, code, err := c.ValidateSession(r.Context(), r.Cookies())
		if err != nil {
			http.Error(w, err.Error(), code)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// DeleteSessionHandler proxies POST /api/session.
func DeleteSessionHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			SessionID string `json:"session_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		resp, code, err := c.DeleteSession(r.Context(), r.Cookies(), req.SessionID)
		if err != nil {
			http.Error(w, err.Error(), code)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// DeleteAllSessionsHandler proxies DELETE /api/session.
func DeleteAllSessionsHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, code, err := c.DeleteAllSessions(r.Context(), r.Cookies())
		if err != nil {
			http.Error(w, err.Error(), code)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// Get2FAStatusHandler proxies GET /api/2fa.
func Get2FAStatusHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, code, err := c.Get2FAStatus(r.Context(), r.Cookies())
		if err != nil {
			http.Error(w, err.Error(), code)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// CheckPendingSessionHandler proxies POST /api/2fa.
func CheckPendingSessionHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, code, err := c.CheckPendingSession(r.Context(), r.Cookies())
		if err != nil {
			http.Error(w, err.Error(), code)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// Delete2FAHandler proxies DELETE /api/2fa.
func Delete2FAHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, code, err := c.Delete2FA(r.Context(), r.Cookies())
		if err != nil {
			http.Error(w, err.Error(), code)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// LockoutUserHandler proxies POST /api/2fa/lockout.
func LockoutUserHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LockoutRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		resp, code, err := c.LockoutUser(r.Context(), r.Cookies(), req.UserID)
		if err != nil {
			http.Error(w, err.Error(), code)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// AuthFlowHandler proxies GET /auth?provider={provider}&next={next}.
func AuthFlowHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		provider := r.URL.Query().Get("provider")
		next := r.URL.Query().Get("next")
		redirect, code, err := c.AuthFlow(r.Context(), r.Cookies(), provider, next)
		if err != nil {
			http.Error(w, err.Error(), code)
			return
		}
		w.Header().Set("hx-redirect", redirect)
		w.WriteHeader(code)
	}
}

// AuthHandler proxies GET/POST /auth/{provider}.
func AuthHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		provider := mux.Vars(r)["provider"]
		next := r.URL.Query().Get("next")

		if r.Method == http.MethodGet {
			redirect, code, err := c.AuthFlow(r.Context(), r.Cookies(), provider, next)
			if err != nil {
				http.Error(w, err.Error(), code)
				return
			}
			w.Header().Set("hx-redirect", redirect)
			w.WriteHeader(code)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}
		message, code, err := c.Auth(r.Context(), r.Cookies(), provider, next, r.Form)
		if err != nil {
			http.Error(w, err.Error(), code)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(code)
		w.Write([]byte(message))
	}
}

// AuthCallbackHandler proxies GET /auth/{provider}/callback.
func AuthCallbackHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		provider := mux.Vars(r)["provider"]
		redirect, code, err := c.AuthCallback(r.Context(), r.Cookies(), provider, r.URL.Query())
		if err != nil {
			http.Error(w, err.Error(), code)
			return
		}
		http.Redirect(w, r, redirect, code)
	}
}

// ============== 2FA Configure ==============

// Configure2FAHandler proxies POST /api/2fa/configure.
func Configure2FAHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Configure2FARequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		resp, code, err := c.Configure2FA(r.Context(), r.Cookies(), req)
		if err != nil {
			http.Error(w, err.Error(), code)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// ============== 2FA Verify ==============

// Verify2FAHandler proxies POST /api/2fa/verify.
func Verify2FAHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Verify2FARequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		resp, code, err := c.Verify2FA(r.Context(), r.Cookies(), req)
		if err != nil {
			http.Error(w, err.Error(), code)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// ============== 2FA Reset ==============

// Reset2FAHandler proxies POST /api/2fa/reset.
func Reset2FAHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Reset2FARequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		resp, code, err := c.Reset2FA(r.Context(), r.Cookies(), req)
		if err != nil {
			http.Error(w, err.Error(), code)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// ============== 2FA Backup Codes ==============

// GenerateBackupCodesHandler proxies POST /api/2fa/backup-codes.
func GenerateBackupCodesHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req GenerateBackupCodesRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		resp, code, err := c.GenerateBackupCodes(r.Context(), r.Cookies(), req)
		if err != nil {
			http.Error(w, err.Error(), code)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// ============== 2FA Trusted Device ==============

// CheckTrustedDeviceHandler proxies GET /api/2fa/trusted-device.
func CheckTrustedDeviceHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, code, err := c.CheckTrustedDevice(r.Context(), r.Cookies())
		if err != nil {
			http.Error(w, err.Error(), code)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// ============== 2FA Lockout ==============

// GetLockoutStatusHandler proxies GET /api/2fa/lockout.
func GetLockoutStatusHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, code, err := c.GetLockoutStatus(r.Context(), r.Cookies())
		if err != nil {
			http.Error(w, err.Error(), code)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// ClearLockoutHandler proxies DELETE /api/2fa/lockout.
func ClearLockoutHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, code, err := c.ClearLockout(r.Context(), r.Cookies())
		if err != nil {
			http.Error(w, err.Error(), code)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// ============== 2FA Recovery ==============

// InitiateRecoveryHandler proxies POST /api/2fa/recovery.
func InitiateRecoveryHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req InitiateRecoveryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		resp, code, err := c.InitiateRecovery(r.Context(), req)
		if err != nil {
			http.Error(w, err.Error(), code)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// VerifyRecoveryCodeHandler proxies POST /api/2fa/recovery/verify.
func VerifyRecoveryCodeHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req VerifyRecoveryCodeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		resp, code, err := c.VerifyRecoveryCode(r.Context(), req)
		if err != nil {
			http.Error(w, err.Error(), code)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}
