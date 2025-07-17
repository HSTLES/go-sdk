# Core Logging Package

This package provides centralized logging functionality for the HSTLES ecosystem.

## Functions

### EnsureSystemUser(db *sql.DB) error
Makes sure a "system-logging" user exists in the database. If not, it generates a condensed UUID and inserts the record.

### LogEvent(db *sql.DB, app, event, desc, userID string) error
Writes into your core logs table. This is a convenience function that calls LogEventWithDetails with empty IP, user agent, and details.

### LogEventWithDetails(db *sql.DB, app, event, desc, userID, ipAddress, userAgent, details string) error
Writes a detailed event into the events table. Includes IP address, user agent, and additional JSON details for security auditing.

## Usage

```go
import "github.com/hstles/go-sdk/core_logging"

// Simple logging
err := core_logging.LogEvent(db, "app.hstles.com", "user_login", "User logged in", userID)

// Detailed logging with security context
err := core_logging.LogEventWithDetails(
    db, 
    "app.hstles.com", 
    "user_login", 
    "User logged in", 
    userID,
    ipAddress,
    userAgent,
    `{"source": "web", "method": "oauth"}`
)
```

## Migration Notice

These functions were previously located in the `core_datastore` package and have been moved here for better separation of concerns.
