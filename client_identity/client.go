package client_identity

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Client wraps calls to your central identity service.
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

// ============== Health & Heartbeat ==============

func (c *Client) Health(ctx context.Context) (HealthResponse, int, error) {
	var resp HealthResponse
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/api/health", nil)
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

func (c *Client) Heartbeat(ctx context.Context) (HeartbeatResponse, int, error) {
	var resp HeartbeatResponse
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/heartbeat", nil)
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

// ============== Plans (Public API) ==============

func (c *Client) ListPlans(ctx context.Context) ([]Plan, int, error) {
	var resp []Plan
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/api/plans", nil)
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

func (c *Client) GetPlan(ctx context.Context, planID string) (Plan, int, error) {
	var resp Plan
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/api/plans/%s", c.BaseURL, planID), nil)
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

// ============== Users (Service API - require API key) ==============

func (c *Client) GetUserByEmail(ctx context.Context, apiKey, email string) (User, int, error) {
	var resp User
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/api/users/email/%s", c.BaseURL, url.QueryEscape(email)), nil)
	req.Header.Set("X-API-Key", apiKey)
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

func (c *Client) GetUserByID(ctx context.Context, apiKey, userID string) (User, int, error) {
	var resp User
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/api/users/%s", c.BaseURL, userID), nil)
	req.Header.Set("X-API-Key", apiKey)
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

func (c *Client) CreateUser(ctx context.Context, apiKey string, req CreateUserRequest) (User, int, error) {
	var resp User
	body, err := json.Marshal(req)
	if err != nil {
		return resp, 0, err
	}
	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/api/users", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-API-Key", apiKey)
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

// ============== Events (Service API - require API key) ==============

func (c *Client) CreateEvent(ctx context.Context, apiKey string, req CreateEventRequest) (Event, int, error) {
	var resp Event
	body, err := json.Marshal(req)
	if err != nil {
		return resp, 0, err
	}
	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/api/events", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-API-Key", apiKey)
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

// ============== Users (Protected API - require session) ==============

func (c *Client) ListUsers(ctx context.Context, cookies []*http.Cookie) ([]User, int, error) {
	var resp []User
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/api/users", nil)
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

func (c *Client) GetUser(ctx context.Context, cookies []*http.Cookie, userID string) (User, int, error) {
	var resp User
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/api/users/%s", c.BaseURL, userID), nil)
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

func (c *Client) UpdateUser(ctx context.Context, cookies []*http.Cookie, userID string, req UpdateUserRequest) (User, int, error) {
	var resp User
	body, err := json.Marshal(req)
	if err != nil {
		return resp, 0, err
	}
	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodPut, fmt.Sprintf("%s/api/users/%s", c.BaseURL, userID), bytes.NewReader(body))
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

func (c *Client) DeleteUser(ctx context.Context, cookies []*http.Cookie, userID string) (DeleteResponse, int, error) {
	var resp DeleteResponse
	req, _ := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s/api/users/%s", c.BaseURL, userID), nil)
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

// ============== Organisations (Protected API) ==============

func (c *Client) ListOrganisations(ctx context.Context, cookies []*http.Cookie) ([]Organisation, int, error) {
	var resp []Organisation
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/api/organisations", nil)
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

func (c *Client) CreateOrganisation(ctx context.Context, cookies []*http.Cookie, req CreateOrganisationRequest) (Organisation, int, error) {
	var resp Organisation
	body, err := json.Marshal(req)
	if err != nil {
		return resp, 0, err
	}
	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/api/organisations", bytes.NewReader(body))
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

func (c *Client) GetOrganisation(ctx context.Context, cookies []*http.Cookie, orgID string) (Organisation, int, error) {
	var resp Organisation
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/api/organisations/%s", c.BaseURL, orgID), nil)
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

func (c *Client) UpdateOrganisation(ctx context.Context, cookies []*http.Cookie, orgID string, req UpdateOrganisationRequest) (Organisation, int, error) {
	var resp Organisation
	body, err := json.Marshal(req)
	if err != nil {
		return resp, 0, err
	}
	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodPut, fmt.Sprintf("%s/api/organisations/%s", c.BaseURL, orgID), bytes.NewReader(body))
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

func (c *Client) DeleteOrganisation(ctx context.Context, cookies []*http.Cookie, orgID string) (DeleteResponse, int, error) {
	var resp DeleteResponse
	req, _ := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s/api/organisations/%s", c.BaseURL, orgID), nil)
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

// ============== Organisation Members (Protected API) ==============

func (c *Client) ListMembers(ctx context.Context, cookies []*http.Cookie, orgID string) ([]Member, int, error) {
	var resp []Member
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/api/organisations/%s/members", c.BaseURL, orgID), nil)
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

func (c *Client) AddMember(ctx context.Context, cookies []*http.Cookie, orgID string, req AddMemberRequest) (Member, int, error) {
	var resp Member
	body, err := json.Marshal(req)
	if err != nil {
		return resp, 0, err
	}
	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/api/organisations/%s/members", c.BaseURL, orgID), bytes.NewReader(body))
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

func (c *Client) UpdateMemberStatus(ctx context.Context, cookies []*http.Cookie, orgID, userID string, req UpdateMemberStatusRequest) (Member, int, error) {
	var resp Member
	body, err := json.Marshal(req)
	if err != nil {
		return resp, 0, err
	}
	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodPut, fmt.Sprintf("%s/api/organisations/%s/members/%s", c.BaseURL, orgID, userID), bytes.NewReader(body))
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

func (c *Client) RemoveMember(ctx context.Context, cookies []*http.Cookie, orgID, userID string) (DeleteResponse, int, error) {
	var resp DeleteResponse
	req, _ := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s/api/organisations/%s/members/%s", c.BaseURL, orgID, userID), nil)
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

// ============== User Organisations (Protected API) ==============

func (c *Client) GetUserOrganisations(ctx context.Context, cookies []*http.Cookie, userID string) ([]Organisation, int, error) {
	var resp []Organisation
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/api/users/%s/organisations", c.BaseURL, userID), nil)
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

// ============== Subscriptions (Protected API) ==============

func (c *Client) CreateSubscription(ctx context.Context, cookies []*http.Cookie, req CreateSubscriptionRequest) (Subscription, int, error) {
	var resp Subscription
	body, err := json.Marshal(req)
	if err != nil {
		return resp, 0, err
	}
	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/api/subscriptions", bytes.NewReader(body))
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

func (c *Client) GetSubscription(ctx context.Context, cookies []*http.Cookie, subscriptionID string) (Subscription, int, error) {
	var resp Subscription
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/api/subscriptions/%s", c.BaseURL, subscriptionID), nil)
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

func (c *Client) UpdateSubscription(ctx context.Context, cookies []*http.Cookie, subscriptionID string, req UpdateSubscriptionRequest) (Subscription, int, error) {
	var resp Subscription
	body, err := json.Marshal(req)
	if err != nil {
		return resp, 0, err
	}
	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodPut, fmt.Sprintf("%s/api/subscriptions/%s", c.BaseURL, subscriptionID), bytes.NewReader(body))
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

func (c *Client) DeleteSubscription(ctx context.Context, cookies []*http.Cookie, subscriptionID string) (DeleteResponse, int, error) {
	var resp DeleteResponse
	req, _ := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s/api/subscriptions/%s", c.BaseURL, subscriptionID), nil)
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

func (c *Client) CancelSubscription(ctx context.Context, cookies []*http.Cookie, subscriptionID string) (Subscription, int, error) {
	var resp Subscription
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/api/subscriptions/%s/cancel", c.BaseURL, subscriptionID), nil)
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

// ============== User Subscriptions (Protected API) ==============

func (c *Client) GetUserSubscriptions(ctx context.Context, cookies []*http.Cookie, userID string) ([]Subscription, int, error) {
	var resp []Subscription
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/api/users/%s/subscriptions", c.BaseURL, userID), nil)
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

func (c *Client) GetActiveSubscription(ctx context.Context, cookies []*http.Cookie, userID string) (Subscription, int, error) {
	var resp Subscription
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/api/users/%s/subscription/active", c.BaseURL, userID), nil)
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

// ============== Events (Protected API) ==============

func (c *Client) ListEvents(ctx context.Context, cookies []*http.Cookie) ([]Event, int, error) {
	var resp []Event
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/api/events", nil)
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

func (c *Client) GetUserEvents(ctx context.Context, cookies []*http.Cookie, userID string) ([]Event, int, error) {
	var resp []Event
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/api/users/%s/events", c.BaseURL, userID), nil)
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
