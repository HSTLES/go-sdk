package client_auth

// SessionResponse is returned by GET /api/session
type SessionResponse struct {
	Valid    bool   `json:"valid"`
	UserID   string `json:"user_id,omitempty"`
	Provider string `json:"provider,omitempty"`
	Message  string `json:"message,omitempty"`
	Error    string `json:"error,omitempty"`
}

// DeleteSessionResponse is returned by POST /api/session
// or DELETE /api/session
type DeleteSessionResponse struct {
	Message string `json:"message"`
}

// TwoFAStatusResponse is returned by GET /api/2fa
type TwoFAStatusResponse struct {
	TwoFactorEnabled bool `json:"two_factor_enabled"`
}

// PendingSessionResponse is returned by POST /api/2fa
type PendingSessionResponse struct {
	Valid    bool   `json:"valid"`
	UserID   string `json:"user_id,omitempty"`
	Provider string `json:"provider,omitempty"`
	Next     string `json:"next,omitempty"`
	Error    string `json:"error,omitempty"`
}

// LockoutRequest is the POST body for /api/2fa/lockout
type LockoutRequest struct {
	UserID string `json:"user_id"`
}

// LockoutResponse is returned by POST /api/2fa/lockout
type LockoutResponse struct {
	Duration int `json:"duration"`
}

// ============== 2FA Configure ==============

// Configure2FARequest represents the request body for 2FA configuration
type Configure2FARequest struct {
	Secret      string `json:"secret"`
	Code        string `json:"code"`
	BackupCodes string `json:"backup_codes"`
}

// Configure2FAResponse represents the response for 2FA configuration
type Configure2FAResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// ============== 2FA Verify ==============

// Verify2FARequest represents the request body for 2FA verification
type Verify2FARequest struct {
	Code           string `json:"code"`
	NextURL        string `json:"next_url,omitempty"`
	RememberDevice bool   `json:"remember_device,omitempty"`
}

// Verify2FAResponse represents the response for 2FA verification
type Verify2FAResponse struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	RedirectURL  string `json:"redirect_url,omitempty"`
	Error        string `json:"error,omitempty"`
	Locked       bool   `json:"locked,omitempty"`
	LockDuration int    `json:"lock_duration,omitempty"`
}

// ============== 2FA Reset ==============

// Reset2FARequest represents the request body for 2FA reset
type Reset2FARequest struct {
	BackupCode string `json:"backup_code,omitempty"`
	Password   string `json:"password,omitempty"`
}

// Reset2FAResponse represents the response for 2FA reset
type Reset2FAResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// ============== 2FA Backup Codes ==============

// GenerateBackupCodesRequest represents the request body for generating backup codes
type GenerateBackupCodesRequest struct {
	Code string `json:"code"` // Current TOTP code to verify before generating new backup codes
}

// GenerateBackupCodesResponse represents the response for backup code generation
type GenerateBackupCodesResponse struct {
	Success     bool     `json:"success"`
	BackupCodes []string `json:"backup_codes,omitempty"`
	Message     string   `json:"message"`
	Error       string   `json:"error,omitempty"`
}

// ============== 2FA Trusted Device ==============

// CheckTrustedDeviceRequest represents the request body for checking trusted device status
type CheckTrustedDeviceRequest struct {
	// This endpoint will use the device fingerprint from headers
}

// CheckTrustedDeviceResponse represents the response for trusted device check
type CheckTrustedDeviceResponse struct {
	Success      bool   `json:"success"`
	IsTrusted    bool   `json:"is_trusted"`
	CanBypass2FA bool   `json:"can_bypass_2fa"`
	Error        string `json:"error,omitempty"`
}

// ============== 2FA Lockout ==============

// LockoutStatusResponse represents the response for lockout status check
type LockoutStatusResponse struct {
	Success         bool   `json:"success"`
	IsLocked        bool   `json:"is_locked"`
	LockMessage     string `json:"lock_message,omitempty"`
	RemainingTime   int    `json:"remaining_time,omitempty"`
	AttemptCount    int    `json:"attempt_count"`
	MaxAttempts     int    `json:"max_attempts"`
	LastAttemptTime string `json:"last_attempt_time,omitempty"`
	Error           string `json:"error,omitempty"`
}

// ClearLockoutResponse represents the response for clearing lockout
type ClearLockoutResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// ============== 2FA Recovery ==============

// InitiateRecoveryRequest represents the request body for initiating 2FA recovery
type InitiateRecoveryRequest struct {
	Email string `json:"email"`
}

// InitiateRecoveryResponse represents the response for initiating recovery
type InitiateRecoveryResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// VerifyRecoveryCodeRequest represents the request body for verifying recovery code
type VerifyRecoveryCodeRequest struct {
	Email        string `json:"email"`
	RecoveryCode string `json:"recovery_code"`
}

// VerifyRecoveryCodeResponse represents the response for verifying recovery code
type VerifyRecoveryCodeResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	AccessToken string `json:"access_token,omitempty"` // Temporary token for 2FA reset
	Error       string `json:"error,omitempty"`
}
