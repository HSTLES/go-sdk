package core_logging

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/hstles/go-sdk/shared_helpers"
)

// EnsureSystemUser makes sure a "system-logging" user exists.
// If not, it generates a condensed UUID and inserts the record.
func EnsureSystemUser(db *sql.DB) error {
	var id string
	err := db.QueryRow(
		`SELECT hstles_user_id
           FROM users
          WHERE name = ?`,
		"system-logging",
	).Scan(&id)

	if err == sql.ErrNoRows {
		// Need to create the system-logging user
		newID, genErr := shared_helpers.GenerateCondensedUUID()
		if genErr != nil {
			return fmt.Errorf("GenerateCondensedUUID: %w", genErr)
		}

		_, insertErr := db.Exec(
			`INSERT INTO users (hstles_user_id, name, email)
                     VALUES (?,             ?,    ?)`,
			newID, "system-logging", "system@logging.com",
		)
		if insertErr != nil {
			return fmt.Errorf("insert system-logging user: %w", insertErr)
		}

		return nil
	}
	if err != nil {
		// Some other query error
		return fmt.Errorf("query system-logging user: %w", err)
	}

	// Already exists
	return nil
}

// LogEvent writes into your core logs table.
// This is a convenience function that calls LogEventWithDetails with empty IP, user agent, and details.
// For HTTP request contexts, use LogEventWithDetails instead to capture IP addresses and user agents.
func LogEvent(db *sql.DB, app, event, desc, userID string) error {
	return LogEventWithDetails(db, app, event, desc, userID, "", "", "")
}

// LogEventWithDetails writes a detailed event into the events table.
// Includes IP address, user agent, and additional JSON details for security auditing.
func LogEventWithDetails(db *sql.DB, app, event, desc, userID, ipAddress, userAgent, details string) error {
	// If no userID provided, look up the "system-logging" user
	if userID == "" {
		row := db.QueryRow(
			`SELECT hstles_user_id FROM users WHERE name = ?`,
			"system-logging",
		)
		if err := row.Scan(&userID); err != nil {
			return fmt.Errorf("lookup system-logging user: %w", err)
		}
	}

	// Insert the event; only wrap and return on error
	if _, err := db.Exec(
		`INSERT INTO events (app_name, event_type, description, timestamp, user_id, ip_address, user_agent, details)
         VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		app,
		event,
		desc,
		time.Now().UTC(),
		userID,
		ipAddress,
		userAgent,
		details,
	); err != nil {
		return fmt.Errorf("log event: %w", err)
	}

	return nil
}
