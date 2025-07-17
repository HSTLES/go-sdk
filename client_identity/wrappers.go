package client_identity

import (
	"context"
	"errors"
	"net/http"
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
		return errors.New("client_identity: not initialized; call Init(baseURL) first")
	}
	return nil
}

// ============== Health & Heartbeat Wrappers ==============

// Health wraps Client.Health on the default client.
func Health(ctx context.Context) (HealthResponse, int, error) {
	if err := ensure(); err != nil {
		return HealthResponse{}, 0, err
	}
	return Default.Health(ctx)
}

// Heartbeat wraps Client.Heartbeat on the default client.
func Heartbeat(ctx context.Context) (HeartbeatResponse, int, error) {
	if err := ensure(); err != nil {
		return HeartbeatResponse{}, 0, err
	}
	return Default.Heartbeat(ctx)
}

// ============== Plan Wrappers ==============

// ListPlans wraps Client.ListPlans on the default client.
func ListPlans(ctx context.Context) ([]Plan, int, error) {
	if err := ensure(); err != nil {
		return nil, 0, err
	}
	return Default.ListPlans(ctx)
}

// GetPlan wraps Client.GetPlan on the default client.
func GetPlan(ctx context.Context, planID string) (Plan, int, error) {
	if err := ensure(); err != nil {
		return Plan{}, 0, err
	}
	return Default.GetPlan(ctx, planID)
}

// ============== User Wrappers (Service API) ==============

// GetUserByEmail wraps Client.GetUserByEmail on the default client.
func GetUserByEmail(ctx context.Context, apiKey, email string) (User, int, error) {
	if err := ensure(); err != nil {
		return User{}, 0, err
	}
	return Default.GetUserByEmail(ctx, apiKey, email)
}

// GetUserByID wraps Client.GetUserByID on the default client.
func GetUserByID(ctx context.Context, apiKey, userID string) (User, int, error) {
	if err := ensure(); err != nil {
		return User{}, 0, err
	}
	return Default.GetUserByID(ctx, apiKey, userID)
}

// CreateUser wraps Client.CreateUser on the default client.
func CreateUser(ctx context.Context, apiKey string, req CreateUserRequest) (User, int, error) {
	if err := ensure(); err != nil {
		return User{}, 0, err
	}
	return Default.CreateUser(ctx, apiKey, req)
}

// ============== Event Wrappers (Service API) ==============

// CreateEvent wraps Client.CreateEvent on the default client.
func CreateEvent(ctx context.Context, apiKey string, req CreateEventRequest) (Event, int, error) {
	if err := ensure(); err != nil {
		return Event{}, 0, err
	}
	return Default.CreateEvent(ctx, apiKey, req)
}

// ============== User Wrappers (Protected API) ==============

// ListUsers wraps Client.ListUsers on the default client.
func ListUsers(ctx context.Context, cookies []*http.Cookie) ([]User, int, error) {
	if err := ensure(); err != nil {
		return nil, 0, err
	}
	return Default.ListUsers(ctx, cookies)
}

// GetUser wraps Client.GetUser on the default client.
func GetUser(ctx context.Context, cookies []*http.Cookie, userID string) (User, int, error) {
	if err := ensure(); err != nil {
		return User{}, 0, err
	}
	return Default.GetUser(ctx, cookies, userID)
}

// UpdateUser wraps Client.UpdateUser on the default client.
func UpdateUser(ctx context.Context, cookies []*http.Cookie, userID string, req UpdateUserRequest) (User, int, error) {
	if err := ensure(); err != nil {
		return User{}, 0, err
	}
	return Default.UpdateUser(ctx, cookies, userID, req)
}

// DeleteUser wraps Client.DeleteUser on the default client.
func DeleteUser(ctx context.Context, cookies []*http.Cookie, userID string) (DeleteResponse, int, error) {
	if err := ensure(); err != nil {
		return DeleteResponse{}, 0, err
	}
	return Default.DeleteUser(ctx, cookies, userID)
}

// ============== Organisation Wrappers ==============

// ListOrganisations wraps Client.ListOrganisations on the default client.
func ListOrganisations(ctx context.Context, cookies []*http.Cookie) ([]Organisation, int, error) {
	if err := ensure(); err != nil {
		return nil, 0, err
	}
	return Default.ListOrganisations(ctx, cookies)
}

// CreateOrganisation wraps Client.CreateOrganisation on the default client.
func CreateOrganisation(ctx context.Context, cookies []*http.Cookie, req CreateOrganisationRequest) (Organisation, int, error) {
	if err := ensure(); err != nil {
		return Organisation{}, 0, err
	}
	return Default.CreateOrganisation(ctx, cookies, req)
}

// GetOrganisation wraps Client.GetOrganisation on the default client.
func GetOrganisation(ctx context.Context, cookies []*http.Cookie, orgID string) (Organisation, int, error) {
	if err := ensure(); err != nil {
		return Organisation{}, 0, err
	}
	return Default.GetOrganisation(ctx, cookies, orgID)
}

// UpdateOrganisation wraps Client.UpdateOrganisation on the default client.
func UpdateOrganisation(ctx context.Context, cookies []*http.Cookie, orgID string, req UpdateOrganisationRequest) (Organisation, int, error) {
	if err := ensure(); err != nil {
		return Organisation{}, 0, err
	}
	return Default.UpdateOrganisation(ctx, cookies, orgID, req)
}

// DeleteOrganisation wraps Client.DeleteOrganisation on the default client.
func DeleteOrganisation(ctx context.Context, cookies []*http.Cookie, orgID string) (DeleteResponse, int, error) {
	if err := ensure(); err != nil {
		return DeleteResponse{}, 0, err
	}
	return Default.DeleteOrganisation(ctx, cookies, orgID)
}

// ============== Organisation Member Wrappers ==============

// ListMembers wraps Client.ListMembers on the default client.
func ListMembers(ctx context.Context, cookies []*http.Cookie, orgID string) ([]Member, int, error) {
	if err := ensure(); err != nil {
		return nil, 0, err
	}
	return Default.ListMembers(ctx, cookies, orgID)
}

// AddMember wraps Client.AddMember on the default client.
func AddMember(ctx context.Context, cookies []*http.Cookie, orgID string, req AddMemberRequest) (Member, int, error) {
	if err := ensure(); err != nil {
		return Member{}, 0, err
	}
	return Default.AddMember(ctx, cookies, orgID, req)
}

// UpdateMemberStatus wraps Client.UpdateMemberStatus on the default client.
func UpdateMemberStatus(ctx context.Context, cookies []*http.Cookie, orgID, userID string, req UpdateMemberStatusRequest) (Member, int, error) {
	if err := ensure(); err != nil {
		return Member{}, 0, err
	}
	return Default.UpdateMemberStatus(ctx, cookies, orgID, userID, req)
}

// RemoveMember wraps Client.RemoveMember on the default client.
func RemoveMember(ctx context.Context, cookies []*http.Cookie, orgID, userID string) (DeleteResponse, int, error) {
	if err := ensure(); err != nil {
		return DeleteResponse{}, 0, err
	}
	return Default.RemoveMember(ctx, cookies, orgID, userID)
}

// ============== User Organisation Wrappers ==============

// GetUserOrganisations wraps Client.GetUserOrganisations on the default client.
func GetUserOrganisations(ctx context.Context, cookies []*http.Cookie, userID string) ([]Organisation, int, error) {
	if err := ensure(); err != nil {
		return nil, 0, err
	}
	return Default.GetUserOrganisations(ctx, cookies, userID)
}

// ============== Subscription Wrappers ==============

// CreateSubscription wraps Client.CreateSubscription on the default client.
func CreateSubscription(ctx context.Context, cookies []*http.Cookie, req CreateSubscriptionRequest) (Subscription, int, error) {
	if err := ensure(); err != nil {
		return Subscription{}, 0, err
	}
	return Default.CreateSubscription(ctx, cookies, req)
}

// GetSubscription wraps Client.GetSubscription on the default client.
func GetSubscription(ctx context.Context, cookies []*http.Cookie, subscriptionID string) (Subscription, int, error) {
	if err := ensure(); err != nil {
		return Subscription{}, 0, err
	}
	return Default.GetSubscription(ctx, cookies, subscriptionID)
}

// UpdateSubscription wraps Client.UpdateSubscription on the default client.
func UpdateSubscription(ctx context.Context, cookies []*http.Cookie, subscriptionID string, req UpdateSubscriptionRequest) (Subscription, int, error) {
	if err := ensure(); err != nil {
		return Subscription{}, 0, err
	}
	return Default.UpdateSubscription(ctx, cookies, subscriptionID, req)
}

// DeleteSubscription wraps Client.DeleteSubscription on the default client.
func DeleteSubscription(ctx context.Context, cookies []*http.Cookie, subscriptionID string) (DeleteResponse, int, error) {
	if err := ensure(); err != nil {
		return DeleteResponse{}, 0, err
	}
	return Default.DeleteSubscription(ctx, cookies, subscriptionID)
}

// CancelSubscription wraps Client.CancelSubscription on the default client.
func CancelSubscription(ctx context.Context, cookies []*http.Cookie, subscriptionID string) (Subscription, int, error) {
	if err := ensure(); err != nil {
		return Subscription{}, 0, err
	}
	return Default.CancelSubscription(ctx, cookies, subscriptionID)
}

// ============== User Subscription Wrappers ==============

// GetUserSubscriptions wraps Client.GetUserSubscriptions on the default client.
func GetUserSubscriptions(ctx context.Context, cookies []*http.Cookie, userID string) ([]Subscription, int, error) {
	if err := ensure(); err != nil {
		return nil, 0, err
	}
	return Default.GetUserSubscriptions(ctx, cookies, userID)
}

// GetActiveSubscription wraps Client.GetActiveSubscription on the default client.
func GetActiveSubscription(ctx context.Context, cookies []*http.Cookie, userID string) (Subscription, int, error) {
	if err := ensure(); err != nil {
		return Subscription{}, 0, err
	}
	return Default.GetActiveSubscription(ctx, cookies, userID)
}

// ============== Event Wrappers (Protected API) ==============

// ListEvents wraps Client.ListEvents on the default client.
func ListEvents(ctx context.Context, cookies []*http.Cookie) ([]Event, int, error) {
	if err := ensure(); err != nil {
		return nil, 0, err
	}
	return Default.ListEvents(ctx, cookies)
}

// GetUserEvents wraps Client.GetUserEvents on the default client.
func GetUserEvents(ctx context.Context, cookies []*http.Cookie, userID string) ([]Event, int, error) {
	if err := ensure(); err != nil {
		return nil, 0, err
	}
	return Default.GetUserEvents(ctx, cookies, userID)
}
