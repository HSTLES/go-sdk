package client_notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"
)

var (
	// Default is the package‐level client.  It starts as nil.
	Default *EmailClient
)

// Init must be called exactly once (or whenever you want to point at a new URL).
// baseURL must include scheme (“https://” or “http://”) or we’ll prepend “https://”.
func Init(baseURL string) {
	if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
		baseURL = "https://" + baseURL
	}
	Default = NewClient(baseURL)
}

// EmailClient knows how to call your notify.hstles.com API.
type EmailClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *EmailClient {
	return &EmailClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *EmailClient) post(ctx context.Context, endpoint string, payload interface{}) (*EmailResponse, error) {
	fullURL := c.baseURL + path.Join("/api/email/", endpoint)
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal %s: %w", endpoint, err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("new %s request: %w", endpoint, err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s request failed: %w", endpoint, err)
	}
	defer resp.Body.Close()

	var er EmailResponse
	if err := json.NewDecoder(resp.Body).Decode(&er); err != nil {
		return nil, fmt.Errorf("%s decode: %w", endpoint, err)
	}
	if !er.Success {
		return &er, fmt.Errorf("%s backend error: %s", endpoint, er.Error)
	}
	return &er, nil
}

// GetStatus checks GET /api/email/status
func (c *EmailClient) GetStatus(ctx context.Context) (*EmailResponse, error) {
	url := c.baseURL + "/api/email/status"
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var er EmailResponse
	if err := json.NewDecoder(resp.Body).Decode(&er); err != nil {
		return nil, err
	}
	return &er, nil
}

// SendWelcomeEmail wraps POST /api/email/welcome
func (c *EmailClient) SendWelcomeEmail(ctx context.Context, to, userName string) (*EmailResponse, error) {
	req := WelcomeEmailRequest{To: to, UserName: userName}
	return c.post(ctx, "welcome", &req)
}

// SendSecurityCodeEmail wraps POST /api/email/security-code
func (c *EmailClient) SendSecurityCodeEmail(ctx context.Context, to, code string) (*EmailResponse, error) {
	req := SecurityCodeEmailRequest{To: to, Code: code}
	return c.post(ctx, "security-code", &req)
}

// SendRecoveryCodeEmail wraps POST /api/email/recovery-code
func (c *EmailClient) SendRecoveryCodeEmail(ctx context.Context, to, userName, code string) (*EmailResponse, error) {
	req := RecoveryCodeEmailRequest{To: to, UserName: userName, Code: code}
	return c.post(ctx, "recovery-code", &req)
}

// SendServiceAlertEmail wraps POST /api/email/service-alert
func (c *EmailClient) SendServiceAlertEmail(ctx context.Context, to, title, message string) (*EmailResponse, error) {
	req := ServiceAlertEmailRequest{To: to, AlertTitle: title, AlertMessage: message}
	return c.post(ctx, "service-alert", &req)
}

// SendLoginLinkEmail wraps POST /api/email/login-link
func (c *EmailClient) SendLoginLinkEmail(ctx context.Context, to, link, name string) (*EmailResponse, error) {
	req := LoginLinkEmailRequest{To: to, LoginLink: link, UserName: name}
	return c.post(ctx, "login-link", &req)
}

// SendGenericEmail wraps POST /api/email/generic
func (c *EmailClient) SendGenericEmail(ctx context.Context, to, subject, message string) (*EmailResponse, error) {
	req := GenericEmailRequest{To: to, Subject: subject, Message: message}
	return c.post(ctx, "generic", &req)
}
