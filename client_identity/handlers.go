package client_identity

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// ============== Health & Heartbeat Handlers ==============

// HealthHandler proxies GET /api/health.
func HealthHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, code, err := c.Health(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// HeartbeatHandler proxies GET /heartbeat.
func HeartbeatHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, code, err := c.Heartbeat(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// ============== Plan Handlers ==============

// ListPlansHandler proxies GET /api/plans.
func ListPlansHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, code, err := c.ListPlans(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// GetPlanHandler proxies GET /api/plans/{id}.
func GetPlanHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		planID := vars["id"]
		resp, code, err := c.GetPlan(r.Context(), planID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// ============== User Handlers (Service API) ==============

// GetUserByEmailHandler proxies GET /api/users/email/{email} with API key.
func GetUserByEmailHandler(c *Client, apiKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		email := vars["email"]
		resp, code, err := c.GetUserByEmail(r.Context(), apiKey, email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// CreateUserHandler proxies POST /api/users with API key.
func CreateUserHandler(c *Client, apiKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		resp, code, err := c.CreateUser(r.Context(), apiKey, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// GetUserByIDHandler proxies GET /api/users/{id} with API key.
func GetUserByIDHandler(c *Client, apiKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["id"]
		resp, code, err := c.GetUserByID(r.Context(), apiKey, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// ============== Event Handlers (Service API) ==============

// CreateEventHandler proxies POST /api/events with API key.
func CreateEventHandler(c *Client, apiKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateEventRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		resp, code, err := c.CreateEvent(r.Context(), apiKey, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// ============== User Handlers (Protected API) ==============

// ListUsersHandler proxies GET /api/users.
func ListUsersHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, code, err := c.ListUsers(r.Context(), r.Cookies())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// GetUserHandler proxies GET /api/users/{id}.
func GetUserHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["id"]
		resp, code, err := c.GetUser(r.Context(), r.Cookies(), userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// UpdateUserHandler proxies PUT /api/users/{id}.
func UpdateUserHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["id"]
		var req UpdateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		resp, code, err := c.UpdateUser(r.Context(), r.Cookies(), userID, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// DeleteUserHandler proxies DELETE /api/users/{id}.
func DeleteUserHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["id"]
		resp, code, err := c.DeleteUser(r.Context(), r.Cookies(), userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// ============== Organisation Handlers ==============

// ListOrganisationsHandler proxies GET /api/organisations.
func ListOrganisationsHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, code, err := c.ListOrganisations(r.Context(), r.Cookies())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// CreateOrganisationHandler proxies POST /api/organisations.
func CreateOrganisationHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateOrganisationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		resp, code, err := c.CreateOrganisation(r.Context(), r.Cookies(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// GetOrganisationHandler proxies GET /api/organisations/{id}.
func GetOrganisationHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		orgID := vars["id"]
		resp, code, err := c.GetOrganisation(r.Context(), r.Cookies(), orgID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// UpdateOrganisationHandler proxies PUT /api/organisations/{id}.
func UpdateOrganisationHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		orgID := vars["id"]
		var req UpdateOrganisationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		resp, code, err := c.UpdateOrganisation(r.Context(), r.Cookies(), orgID, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// DeleteOrganisationHandler proxies DELETE /api/organisations/{id}.
func DeleteOrganisationHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		orgID := vars["id"]
		resp, code, err := c.DeleteOrganisation(r.Context(), r.Cookies(), orgID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// ============== Organisation Member Handlers ==============

// ListMembersHandler proxies GET /api/organisations/{id}/members.
func ListMembersHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		orgID := vars["id"]
		resp, code, err := c.ListMembers(r.Context(), r.Cookies(), orgID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// AddMemberHandler proxies POST /api/organisations/{id}/members.
func AddMemberHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		orgID := vars["id"]
		var req AddMemberRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		resp, code, err := c.AddMember(r.Context(), r.Cookies(), orgID, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// UpdateMemberStatusHandler proxies PUT /api/organisations/{id}/members/{user_id}.
func UpdateMemberStatusHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		orgID := vars["id"]
		userID := vars["user_id"]
		var req UpdateMemberStatusRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		resp, code, err := c.UpdateMemberStatus(r.Context(), r.Cookies(), orgID, userID, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// RemoveMemberHandler proxies DELETE /api/organisations/{id}/members/{user_id}.
func RemoveMemberHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		orgID := vars["id"]
		userID := vars["user_id"]
		resp, code, err := c.RemoveMember(r.Context(), r.Cookies(), orgID, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// ============== User Organisation Handlers ==============

// GetUserOrganisationsHandler proxies GET /api/users/{user_id}/organisations.
func GetUserOrganisationsHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["user_id"]
		resp, code, err := c.GetUserOrganisations(r.Context(), r.Cookies(), userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// ============== Subscription Handlers ==============

// CreateSubscriptionHandler proxies POST /api/subscriptions.
func CreateSubscriptionHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateSubscriptionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		resp, code, err := c.CreateSubscription(r.Context(), r.Cookies(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// GetSubscriptionHandler proxies GET /api/subscriptions/{id}.
func GetSubscriptionHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		subscriptionID := vars["id"]
		resp, code, err := c.GetSubscription(r.Context(), r.Cookies(), subscriptionID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// UpdateSubscriptionHandler proxies PUT /api/subscriptions/{id}.
func UpdateSubscriptionHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		subscriptionID := vars["id"]
		var req UpdateSubscriptionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		resp, code, err := c.UpdateSubscription(r.Context(), r.Cookies(), subscriptionID, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// DeleteSubscriptionHandler proxies DELETE /api/subscriptions/{id}.
func DeleteSubscriptionHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		subscriptionID := vars["id"]
		resp, code, err := c.DeleteSubscription(r.Context(), r.Cookies(), subscriptionID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// CancelSubscriptionHandler proxies POST /api/subscriptions/{id}/cancel.
func CancelSubscriptionHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		subscriptionID := vars["id"]
		resp, code, err := c.CancelSubscription(r.Context(), r.Cookies(), subscriptionID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// ============== User Subscription Handlers ==============

// GetUserSubscriptionsHandler proxies GET /api/users/{user_id}/subscriptions.
func GetUserSubscriptionsHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["user_id"]
		resp, code, err := c.GetUserSubscriptions(r.Context(), r.Cookies(), userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// GetActiveSubscriptionHandler proxies GET /api/users/{user_id}/subscription/active.
func GetActiveSubscriptionHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["user_id"]
		resp, code, err := c.GetActiveSubscription(r.Context(), r.Cookies(), userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// ============== Event Handlers (Protected API) ==============

// ListEventsHandler proxies GET /api/events.
func ListEventsHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, code, err := c.ListEvents(r.Context(), r.Cookies())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

// GetUserEventsHandler proxies GET /api/users/{user_id}/events.
func GetUserEventsHandler(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["user_id"]
		resp, code, err := c.GetUserEvents(r.Context(), r.Cookies(), userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}
