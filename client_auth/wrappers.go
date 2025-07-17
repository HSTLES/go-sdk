package client_auth

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// Default is the package-level client. Call Init once before using.
var Default *Client

// Init sets up the Default client. Call this at application startup.
func Init(baseURL string) {
	Default = NewClient(baseURL)
}

// ensure checks that Default has been initialized.
func ensure() error {
	if Default == nil {
		return errors.New("client_auth: not initialized; call Init(baseURL) first")
	}
	return nil
}

// ValidateSession wraps Client.ValidateSession on the default client.
func ValidateSession(ctx context.Context, cookies []*http.Cookie) (SessionResponse, int, error) {
	if err := ensure(); err != nil {
		return SessionResponse{}, 0, err
	}
	return Default.ValidateSession(ctx, cookies)
}

// DeleteSession wraps Client.DeleteSession on the default client.
func DeleteSession(ctx context.Context, cookies []*http.Cookie, sessionID string) (DeleteSessionResponse, int, error) {
	if err := ensure(); err != nil {
		return DeleteSessionResponse{}, 0, err
	}
	return Default.DeleteSession(ctx, cookies, sessionID)
}

// DeleteAllSessions wraps Client.DeleteAllSessions on the default client.
func DeleteAllSessions(ctx context.Context, cookies []*http.Cookie) (DeleteSessionResponse, int, error) {
	if err := ensure(); err != nil {
		return DeleteSessionResponse{}, 0, err
	}
	return Default.DeleteAllSessions(ctx, cookies)
}

// Get2FAStatus wraps Client.Get2FAStatus on the default client.
func Get2FAStatus(ctx context.Context, cookies []*http.Cookie) (TwoFAStatusResponse, int, error) {
	if err := ensure(); err != nil {
		return TwoFAStatusResponse{}, 0, err
	}
	return Default.Get2FAStatus(ctx, cookies)
}

// CheckPendingSession wraps Client.CheckPendingSession on the default client.
func CheckPendingSession(ctx context.Context, cookies []*http.Cookie) (PendingSessionResponse, int, error) {
	if err := ensure(); err != nil {
		return PendingSessionResponse{}, 0, err
	}
	return Default.CheckPendingSession(ctx, cookies)
}

// Delete2FA wraps Client.Delete2FA on the default client.
func Delete2FA(ctx context.Context, cookies []*http.Cookie) (DeleteSessionResponse, int, error) {
	if err := ensure(); err != nil {
		return DeleteSessionResponse{}, 0, err
	}
	return Default.Delete2FA(ctx, cookies)
}

// LockoutUser wraps Client.LockoutUser on the default client.
func LockoutUser(ctx context.Context, cookies []*http.Cookie, userID string) (LockoutResponse, int, error) {
	if err := ensure(); err != nil {
		return LockoutResponse{}, 0, err
	}
	return Default.LockoutUser(ctx, cookies, userID)
}

// AuthFlow wraps Client.AuthFlow on the default client.
func AuthFlow(ctx context.Context, cookies []*http.Cookie, provider, next string) (string, int, error) {
	if err := ensure(); err != nil {
		return "", 0, err
	}
	return Default.AuthFlow(ctx, cookies, provider, next)
}

// Auth wraps Client.Auth on the default client.
func Auth(ctx context.Context, cookies []*http.Cookie, provider, next string, form url.Values) (string, int, error) {
	if err := ensure(); err != nil {
		return "", 0, err
	}
	return Default.Auth(ctx, cookies, provider, next, form)
}

// AuthCallback wraps Client.AuthCallback on the default client.
func AuthCallback(ctx context.Context, cookies []*http.Cookie, provider string, query url.Values) (string, int, error) {
	if err := ensure(); err != nil {
		return "", 0, err
	}
	return Default.AuthCallback(ctx, cookies, provider, query)
}

// ============== 2FA Configure ==============

// Configure2FA wraps Client.Configure2FA on the default client.
func Configure2FA(ctx context.Context, cookies []*http.Cookie, req Configure2FARequest) (Configure2FAResponse, int, error) {
	if err := ensure(); err != nil {
		return Configure2FAResponse{}, 0, err
	}
	return Default.Configure2FA(ctx, cookies, req)
}

// ============== 2FA Verify ==============

// Verify2FA wraps Client.Verify2FA on the default client.
func Verify2FA(ctx context.Context, cookies []*http.Cookie, req Verify2FARequest) (Verify2FAResponse, int, error) {
	if err := ensure(); err != nil {
		return Verify2FAResponse{}, 0, err
	}
	return Default.Verify2FA(ctx, cookies, req)
}

// ============== 2FA Reset ==============

// Reset2FA wraps Client.Reset2FA on the default client.
func Reset2FA(ctx context.Context, cookies []*http.Cookie, req Reset2FARequest) (Reset2FAResponse, int, error) {
	if err := ensure(); err != nil {
		return Reset2FAResponse{}, 0, err
	}
	return Default.Reset2FA(ctx, cookies, req)
}

// ============== 2FA Backup Codes ==============

// GenerateBackupCodes wraps Client.GenerateBackupCodes on the default client.
func GenerateBackupCodes(ctx context.Context, cookies []*http.Cookie, req GenerateBackupCodesRequest) (GenerateBackupCodesResponse, int, error) {
	if err := ensure(); err != nil {
		return GenerateBackupCodesResponse{}, 0, err
	}
	return Default.GenerateBackupCodes(ctx, cookies, req)
}

// ============== 2FA Trusted Device ==============

// CheckTrustedDevice wraps Client.CheckTrustedDevice on the default client.
func CheckTrustedDevice(ctx context.Context, cookies []*http.Cookie) (CheckTrustedDeviceResponse, int, error) {
	if err := ensure(); err != nil {
		return CheckTrustedDeviceResponse{}, 0, err
	}
	return Default.CheckTrustedDevice(ctx, cookies)
}

// ============== 2FA Lockout ==============

// GetLockoutStatus wraps Client.GetLockoutStatus on the default client.
func GetLockoutStatus(ctx context.Context, cookies []*http.Cookie) (LockoutStatusResponse, int, error) {
	if err := ensure(); err != nil {
		return LockoutStatusResponse{}, 0, err
	}
	return Default.GetLockoutStatus(ctx, cookies)
}

// ClearLockout wraps Client.ClearLockout on the default client.
func ClearLockout(ctx context.Context, cookies []*http.Cookie) (ClearLockoutResponse, int, error) {
	if err := ensure(); err != nil {
		return ClearLockoutResponse{}, 0, err
	}
	return Default.ClearLockout(ctx, cookies)
}

// ============== 2FA Recovery ==============

// InitiateRecovery wraps Client.InitiateRecovery on the default client.
func InitiateRecovery(ctx context.Context, req InitiateRecoveryRequest) (InitiateRecoveryResponse, int, error) {
	if err := ensure(); err != nil {
		return InitiateRecoveryResponse{}, 0, err
	}
	return Default.InitiateRecovery(ctx, req)
}

// VerifyRecoveryCode wraps Client.VerifyRecoveryCode on the default client.
func VerifyRecoveryCode(ctx context.Context, req VerifyRecoveryCodeRequest) (VerifyRecoveryCodeResponse, int, error) {
	if err := ensure(); err != nil {
		return VerifyRecoveryCodeResponse{}, 0, err
	}
	return Default.VerifyRecoveryCode(ctx, req)
}
