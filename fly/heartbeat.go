package fly

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hstles/go-sdk/core_datastore"
	"github.com/hstles/go-sdk/core_logging"
	"github.com/hstles/go-sdk/shared_helpers"
)

type HeartbeatResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Services  map[string]string `json:"services"`
}

// HeartbeatHandler returns an http.HandlerFunc that checks CoreDB health
// and reports JSON status, logging any failures.
func HeartbeatHandler(mgr *core_datastore.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := HeartbeatResponse{
			Timestamp: time.Now().UTC(),
			Services:  make(map[string]string),
		}

		// 1) Ping CoreDB
		if err := mgr.CoreDB.Ping(); err != nil {
			response.Status = "unhealthy"
			response.Services["database"] = fmt.Sprintf("ping failed: %v", err)
			_ = core_logging.LogEventWithDetails(
				mgr.CoreDB,
				"auth.hstles.com",
				"heartbeat_db_unhealthy",
				fmt.Sprintf("CoreDB ping failed: %v", err),
				"",
				shared_helpers.GetClientIP(r),
				shared_helpers.GetUserAgent(r),
				"",
			)
		} else {
			// 2) Verify simple query
			var count int
			if err := mgr.CoreDB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count); err != nil {
				response.Status = "unhealthy"
				response.Services["database"] = fmt.Sprintf("query failed: %v", err)
				_ = core_logging.LogEventWithDetails(
					mgr.CoreDB,
					"auth.hstles.com",
					"heartbeat_db_query_failed",
					fmt.Sprintf("CoreDB query failed: %v", err),
					"",
					shared_helpers.GetClientIP(r),
					shared_helpers.GetUserAgent(r),
					"",
				)
			} else {
				response.Services["database"] = "healthy"
			}
		}

		// 3) Default to healthy if nothing flagged
		if response.Status == "" {
			response.Status = "healthy"
		}

		// 4) Write JSON response
		w.Header().Set("Content-Type", "application/json")
		if response.Status != "healthy" {
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("HeartbeatHandler: encode error: %v", err)
		}
	}
}
