package client_auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client wraps calls to your central auth service.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient constructs a new client.
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) ValidateSession(ctx context.Context, cookies []*http.Cookie) (SessionResponse, int, error) {
	var resp SessionResponse
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/api/session", nil)
	for _, ck := range cookies {
		req.AddCookie(ck)
	}
	r, err := c.HTTPClient.Do(req)
	if err != nil {
		return resp, 0, err
	}
	defer r.Body.Close()
	status := r.StatusCode
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return resp, status, err
	}
	return resp, status, nil
}

func (c *Client) DeleteSession(ctx context.Context, cookies []*http.Cookie, sessionID string) (DeleteSessionResponse, int, error) {
	var resp DeleteSessionResponse
	body, _ := json.Marshal(map[string]string{"session_id": sessionID})
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/api/session", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	for _, ck := range cookies {
		req.AddCookie(ck)
	}
	r, err := c.HTTPClient.Do(req)
	if err != nil {
		return resp, 0, err
	}
	defer r.Body.Close()
	status := r.StatusCode
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return resp, status, err
	}
	return resp, status, nil
}

func (c *Client) DeleteAllSessions(ctx context.Context, cookies []*http.Cookie) (DeleteSessionResponse, int, error) {
	var resp DeleteSessionResponse
	req, _ := http.NewRequestWithContext(ctx, http.MethodDelete, c.BaseURL+"/api/session", nil)
	for _, ck := range cookies {
		req.AddCookie(ck)
	}
	r, err := c.HTTPClient.Do(req)
	if err != nil {
		return resp, 0, err
	}
	defer r.Body.Close()
	status := r.StatusCode
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return resp, status, err
	}
	return resp, status, nil
}

func (c *Client) Get2FAStatus(ctx context.Context, cookies []*http.Cookie) (TwoFAStatusResponse, int, error) {
	var resp TwoFAStatusResponse
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/api/2fa", nil)
	for _, ck := range cookies {
		req.AddCookie(ck)
	}
	r, err := c.HTTPClient.Do(req)
	if err != nil {
		return resp, 0, err
	}
	defer r.Body.Close()
	status := r.StatusCode
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return resp, status, err
	}
	return resp, status, nil
}

func (c *Client) CheckPendingSession(ctx context.Context, cookies []*http.Cookie) (PendingSessionResponse, int, error) {
	var resp PendingSessionResponse
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/api/2fa", nil)
	for _, ck := range cookies {
		req.AddCookie(ck)
	}
	r, err := c.HTTPClient.Do(req)
	if err != nil {
		return resp, 0, err
	}
	defer r.Body.Close()
	status := r.StatusCode
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return resp, status, err
	}
	return resp, status, nil
}

func (c *Client) Delete2FA(ctx context.Context, cookies []*http.Cookie) (DeleteSessionResponse, int, error) {
	var resp DeleteSessionResponse
	req, _ := http.NewRequestWithContext(ctx, http.MethodDelete, c.BaseURL+"/api/2fa", nil)
	for _, ck := range cookies {
		req.AddCookie(ck)
	}
	r, err := c.HTTPClient.Do(req)
	if err != nil {
		return resp, 0, err
	}
	defer r.Body.Close()
	status := r.StatusCode
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return resp, status, err
	}
	return resp, status, nil
}

func (c *Client) CheckLockout(ctx context.Context, cookies []*http.Cookie, req LockoutRequest) (LockoutResponse, int, error) {
	var resp LockoutResponse

	body, err := json.Marshal(req)
	if err != nil {
		return resp, 0, err
	}
	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/api/2fa/lockout", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	for _, ck := range cookies {
		httpReq.AddCookie(ck)
	}

	r, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return resp, 0, err
	}
	defer r.Body.Close()

	status := r.StatusCode
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return resp, status, err
	}
	return resp, status, nil
}

func (c *Client) LockoutUser(ctx context.Context, cookies []*http.Cookie, userID string) (LockoutResponse, int, error) {
	var resp LockoutResponse
	body, _ := json.Marshal(LockoutRequest{UserID: userID})
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/api/2fa/lockout", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	for _, ck := range cookies {
		req.AddCookie(ck)
	}
	r, err := c.HTTPClient.Do(req)
	if err != nil {
		return resp, 0, err
	}
	defer r.Body.Close()
	status := r.StatusCode
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return resp, status, err
	}
	return resp, status, nil
}

func (c *Client) AuthFlow(ctx context.Context, cookies []*http.Cookie, provider, next string) (string, int, error) {
	q := url.Values{"provider": {provider}, "next": {next}}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/auth?"+q.Encode(), nil)
	for _, ck := range cookies {
		req.AddCookie(ck)
	}
	r, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer r.Body.Close()
	if loc := r.Header.Get("hx-redirect"); loc != "" {
		return loc, r.StatusCode, nil
	}
	if loc, err := r.Location(); err == nil {
		return loc.String(), r.StatusCode, nil
	}
	body, _ := io.ReadAll(r.Body)
	return string(body), r.StatusCode, nil
}

func (c *Client) Auth(ctx context.Context, cookies []*http.Cookie, provider, next string, form url.Values) (string, int, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("%s/auth/%s?next=%s", c.BaseURL, provider, url.QueryEscape(next)),
		strings.NewReader(form.Encode()),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for _, ck := range cookies {
		req.AddCookie(ck)
	}
	r, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer r.Body.Close()
	b, _ := io.ReadAll(r.Body)
	return string(b), r.StatusCode, nil
}

func (c *Client) AuthCallback(ctx context.Context, cookies []*http.Cookie, provider string, query url.Values) (string, int, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%s/auth/%s/callback?%s", c.BaseURL, provider, query.Encode()), nil,
	)
	for _, ck := range cookies {
		req.AddCookie(ck)
	}
	r, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer r.Body.Close()
	if loc := r.Header.Get("Location"); loc != "" {
		return loc, r.StatusCode, nil
	}
	if loc := r.Header.Get("hx-redirect"); loc != "" {
		return loc, r.StatusCode, nil
	}
	return "", r.StatusCode, fmt.Errorf("no redirect in callback")
}

// ============== 2FA Configure ==============

func (c *Client) Configure2FA(ctx context.Context, cookies []*http.Cookie, req Configure2FARequest) (Configure2FAResponse, int, error) {
	var resp Configure2FAResponse
	body, err := json.Marshal(req)
	if err != nil {
		return resp, 0, err
	}
	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/api/2fa/configure", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	for _, ck := range cookies {
		httpReq.AddCookie(ck)
	}
	r, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return resp, 0, err
	}
	defer r.Body.Close()
	status := r.StatusCode
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return resp, status, err
	}
	return resp, status, nil
}

// ============== 2FA Verify ==============

func (c *Client) Verify2FA(ctx context.Context, cookies []*http.Cookie, req Verify2FARequest) (Verify2FAResponse, int, error) {
	var resp Verify2FAResponse
	body, err := json.Marshal(req)
	if err != nil {
		return resp, 0, err
	}
	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/api/2fa/verify", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	for _, ck := range cookies {
		httpReq.AddCookie(ck)
	}
	r, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return resp, 0, err
	}
	defer r.Body.Close()
	status := r.StatusCode
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return resp, status, err
	}
	return resp, status, nil
}

// ============== 2FA Reset ==============

func (c *Client) Reset2FA(ctx context.Context, cookies []*http.Cookie, req Reset2FARequest) (Reset2FAResponse, int, error) {
	var resp Reset2FAResponse
	body, err := json.Marshal(req)
	if err != nil {
		return resp, 0, err
	}
	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/api/2fa/reset", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	for _, ck := range cookies {
		httpReq.AddCookie(ck)
	}
	r, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return resp, 0, err
	}
	defer r.Body.Close()
	status := r.StatusCode
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return resp, status, err
	}
	return resp, status, nil
}

// ============== 2FA Backup Codes ==============

func (c *Client) GenerateBackupCodes(ctx context.Context, cookies []*http.Cookie, req GenerateBackupCodesRequest) (GenerateBackupCodesResponse, int, error) {
	var resp GenerateBackupCodesResponse
	body, err := json.Marshal(req)
	if err != nil {
		return resp, 0, err
	}
	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/api/2fa/backup-codes", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	for _, ck := range cookies {
		httpReq.AddCookie(ck)
	}
	r, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return resp, 0, err
	}
	defer r.Body.Close()
	status := r.StatusCode
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return resp, status, err
	}
	return resp, status, nil
}

// ============== 2FA Trusted Device ==============

func (c *Client) CheckTrustedDevice(ctx context.Context, cookies []*http.Cookie) (CheckTrustedDeviceResponse, int, error) {
	var resp CheckTrustedDeviceResponse
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/api/2fa/trusted-device", nil)
	for _, ck := range cookies {
		req.AddCookie(ck)
	}
	r, err := c.HTTPClient.Do(req)
	if err != nil {
		return resp, 0, err
	}
	defer r.Body.Close()
	status := r.StatusCode
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return resp, status, err
	}
	return resp, status, nil
}

// ============== 2FA Lockout ==============

func (c *Client) GetLockoutStatus(ctx context.Context, cookies []*http.Cookie) (LockoutStatusResponse, int, error) {
	var resp LockoutStatusResponse
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/api/2fa/lockout", nil)
	for _, ck := range cookies {
		req.AddCookie(ck)
	}
	r, err := c.HTTPClient.Do(req)
	if err != nil {
		return resp, 0, err
	}
	defer r.Body.Close()
	status := r.StatusCode
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return resp, status, err
	}
	return resp, status, nil
}

func (c *Client) ClearLockout(ctx context.Context, cookies []*http.Cookie) (ClearLockoutResponse, int, error) {
	var resp ClearLockoutResponse
	req, _ := http.NewRequestWithContext(ctx, http.MethodDelete, c.BaseURL+"/api/2fa/lockout", nil)
	for _, ck := range cookies {
		req.AddCookie(ck)
	}
	r, err := c.HTTPClient.Do(req)
	if err != nil {
		return resp, 0, err
	}
	defer r.Body.Close()
	status := r.StatusCode
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return resp, status, err
	}
	return resp, status, nil
}

// ============== 2FA Recovery ==============

func (c *Client) InitiateRecovery(ctx context.Context, req InitiateRecoveryRequest) (InitiateRecoveryResponse, int, error) {
	var resp InitiateRecoveryResponse
	body, err := json.Marshal(req)
	if err != nil {
		return resp, 0, err
	}
	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/api/2fa/recovery", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	// Note: No cookies for recovery initiation - this is for when user is locked out
	r, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return resp, 0, err
	}
	defer r.Body.Close()
	status := r.StatusCode
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return resp, status, err
	}
	return resp, status, nil
}

func (c *Client) VerifyRecoveryCode(ctx context.Context, req VerifyRecoveryCodeRequest) (VerifyRecoveryCodeResponse, int, error) {
	var resp VerifyRecoveryCodeResponse
	body, err := json.Marshal(req)
	if err != nil {
		return resp, 0, err
	}
	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/api/2fa/recovery/verify", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	// Note: No cookies for recovery verification - this is for when user is locked out
	r, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return resp, 0, err
	}
	defer r.Body.Close()
	status := r.StatusCode
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return resp, status, err
	}
	return resp, status, nil
}
