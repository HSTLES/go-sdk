package client_identity

import "time"

// ============== Health & Heartbeat ==============

// HealthResponse is returned by GET /api/health
type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// HeartbeatResponse is returned by GET /heartbeat
type HeartbeatResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

// ============== Plans ==============

// Plan represents a subscription plan
type Plan struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Currency    string    `json:"currency"`
	Interval    string    `json:"interval"` // monthly, yearly, etc.
	Features    []string  `json:"features,omitempty"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ============== Users ==============

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateUserRequest is the request body for POST /api/users
type CreateUserRequest struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

// UpdateUserRequest is the request body for PUT /api/users/{id}
type UpdateUserRequest struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Active    *bool   `json:"active,omitempty"`
}

// ============== Organisations ==============

// Organisation represents an organisation
type Organisation struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateOrganisationRequest is the request body for POST /api/organisations
type CreateOrganisationRequest struct {
	Name string `json:"name"`
}

// UpdateOrganisationRequest is the request body for PUT /api/organisations/{id}
type UpdateOrganisationRequest struct {
	Name   *string `json:"name,omitempty"`
	Active *bool   `json:"active,omitempty"`
}

// ============== Organisation Members ==============

// Member represents an organisation member
type Member struct {
	UserID   string    `json:"user_id"`
	User     *User     `json:"user,omitempty"`
	Status   string    `json:"status"` // active, inactive, pending
	Role     string    `json:"role"`   // admin, member, etc.
	JoinedAt time.Time `json:"joined_at"`
}

// AddMemberRequest is the request body for POST /api/organisations/{id}/members
type AddMemberRequest struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

// UpdateMemberStatusRequest is the request body for PUT /api/organisations/{id}/members/{user_id}
type UpdateMemberStatusRequest struct {
	Status string `json:"status"`
	Role   string `json:"role,omitempty"`
}

// ============== Subscriptions ==============

// Subscription represents a user subscription
type Subscription struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	PlanID      string     `json:"plan_id"`
	Plan        *Plan      `json:"plan,omitempty"`
	Status      string     `json:"status"` // active, inactive, cancelled, expired
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	CancelledAt *time.Time `json:"cancelled_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// CreateSubscriptionRequest is the request body for POST /api/subscriptions
type CreateSubscriptionRequest struct {
	UserID string `json:"user_id"`
	PlanID string `json:"plan_id"`
}

// UpdateSubscriptionRequest is the request body for PUT /api/subscriptions/{id}
type UpdateSubscriptionRequest struct {
	Status  *string    `json:"status,omitempty"`
	EndDate *time.Time `json:"end_date,omitempty"`
}

// ============== Events ==============

// Event represents a system event
type Event struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id,omitempty"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Metadata    string    `json:"metadata,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateEventRequest is the request body for POST /api/events
type CreateEventRequest struct {
	UserID      string `json:"user_id,omitempty"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Metadata    string `json:"metadata,omitempty"`
}

// ============== Generic Responses ==============

// DeleteResponse is a generic response for delete operations
type DeleteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ErrorResponse is a generic error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
