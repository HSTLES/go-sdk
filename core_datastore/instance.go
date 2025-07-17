package core_datastore

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TursoClient wraps Turso Platform API HTTP calls.
type TursoClient struct {
	baseURL   string
	authToken string
	http      *http.Client
}

// NewTursoClient creates a new TursoClient.
// If httpClient is nil, http.DefaultClient will be used.
func NewTursoClient(baseURL, authToken string, httpClient *http.Client) *TursoClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &TursoClient{
		baseURL:   baseURL,
		authToken: authToken,
		http:      httpClient,
	}
}

// CreateInstanceRequest holds the parameters for creating a new Turso instance.
type CreateInstanceRequest struct {
	Name   string `json:"name"`
	Region string `json:"region,omitempty"`
	// Add other fields here (e.g. "config", "replica_of", etc.)
}

// InstanceInfo represents metadata about a Turso instance.
type InstanceInfo struct {
	Name      string `json:"name"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at,omitempty"`
	// Add other fields returned by the API as needed
}

// CreateInstance provisions a new Turso database instance.
func (c *TursoClient) CreateInstance(ctx context.Context, req CreateInstanceRequest) (*InstanceInfo, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal CreateInstanceRequest: %w", err)
	}

	endpoint := fmt.Sprintf("%s/instances", c.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("new request POST %s: %w", endpoint, err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+c.authToken)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("do POST %s: %w", endpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %s: %s", resp.Status, string(body))
	}

	var info InstanceInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("decode CreateInstance response: %w", err)
	}
	return &info, nil
}

// ListInstances fetches all Turso instances under the current account.
func (c *TursoClient) ListInstances(ctx context.Context) ([]InstanceInfo, error) {
	endpoint := fmt.Sprintf("%s/instances", c.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("new request GET %s: %w", endpoint, err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+c.authToken)

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("do GET %s: %w", endpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %s: %s", resp.Status, string(body))
	}

	var instances []InstanceInfo
	if err := json.NewDecoder(resp.Body).Decode(&instances); err != nil {
		return nil, fmt.Errorf("decode ListInstances response: %w", err)
	}
	return instances, nil
}

// GetInstance retrieves details for a single Turso instance by name.
func (c *TursoClient) GetInstance(ctx context.Context, name string) (*InstanceInfo, error) {
	endpoint := fmt.Sprintf("%s/instances/%s", c.baseURL, name)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("new request GET %s: %w", endpoint, err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+c.authToken)

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("do GET %s: %w", endpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %s: %s", resp.Status, string(body))
	}

	var info InstanceInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("decode GetInstance response: %w", err)
	}
	return &info, nil
}

// DeleteInstance deprovisions a Turso instance permanently.
func (c *TursoClient) DeleteInstance(ctx context.Context, name string) error {
	endpoint := fmt.Sprintf("%s/instances/%s", c.baseURL, name)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return fmt.Errorf("new request DELETE %s: %w", endpoint, err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+c.authToken)

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return fmt.Errorf("do DELETE %s: %w", endpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status %s: %s", resp.Status, string(body))
	}
	return nil
}
