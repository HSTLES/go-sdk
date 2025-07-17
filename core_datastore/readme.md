# core\_datastore

A Go package for managing your “Core” database and provisioning Turso SQLite instances via the Turso Platform API. It bundles together:

* **Connection helpers** for opening libsql (Turso) databases with PRAGMA-enforced foreign keys
* **CoreDB utilities** for ensuring a system-logging user and writing structured log events
* A **TursoClient** for Create/List/Get/Delete operations on Turso instances
* A **Manager** struct that ties your CoreDB and Turso API client under a single handle

---

## Features

* **ConnectToInstance**
  Opens a `*sql.DB` to any Turso instance URL + auth token (via `libsql`). Enables foreign keys and verifies connectivity.

* **EnsureSystemUser**
  Verifies that a `users` table contains a “system-logging” user. If missing, generates a condensed UUID and inserts it.

* **LogEvent**
  Inserts structured events (`app_name`, `event_type`, `description`, timestamp, user ID) into a `logs` table.

* **TursoClient**
  A simple HTTP client for Turso’s REST API. Supports:

  * `CreateInstance`
  * `ListInstances`
  * `GetInstance`
  * `DeleteInstance`

* **Manager**
  Combines a `CoreDB` (`*sql.DB`) connection and a `TursoClient` for easy bootstrap and teardown.

---

## Installation

```bash
go get github.com/hstles/go-sdk/core_datastore
```

## Quick Start

```go
package main

import (
    "context"
    "log"
    "net/http"

    "github.com/hstles/go-sdk/core_config"
    "github.com/hstles/go-sdk/core_datastore"
)

func main() {
    // Load configuration
    cfg, err := core_config.LoadCoreConfig()
    if err != nil {
        log.Fatalf("load config: %v", err)
    }

    // Initialize Manager (CoreDB + Turso API client)
    manager, err := core_datastore.NewManager(cfg)
    if err != nil {
        log.Fatalf("new manager: %v", err)
    }
    defer manager.Close()

    // Ensure system-logging user
    if err := core_datastore.EnsureSystemUser(manager.CoreDB); err != nil {
        log.Fatalf("ensure system user: %v", err)
    }

    // Log startup event
    if err := core_datastore.LogEvent(
        manager.CoreDB,
        "myApp",
        "startup",
        "Application has started",
        "",
    ); err != nil {
        log.Printf("log event error: %v", err)
    }

    // Connect to a Turso instance
    db, err := core_datastore.ConnectToInstance(
        cfg.TursoCoreDBURL,
        cfg.TursoAuthToken,
    )
    if err != nil {
        log.Fatalf("connect to instance: %v", err)
    }
    defer db.Close()

    // Provision a new instance
    ctx := context.Background()
    info, err := manager.Turso().CreateInstance(ctx, core_datastore.CreateInstanceRequest{
        Name:   "example-instance",
        Region: "us-east-1",
    })
    if err != nil {
        log.Fatalf("create instance: %v", err)
    }
    log.Printf("Instance created: %#v", info)
}
```

### Overriding the Turso API Endpoint

By default, the `TursoClient` uses `https://api.turso.sh/v1`. To target a different endpoint, set `TURSO_API_URL`:

```bash
export TURSO_API_URL="https://staging-api.turso.sh/v1"
```

Then initialize:

```go
apiURL := core_config.RequireEnvDefault(
    "TURSO_API_URL",
    "https://api.turso.sh/v1",
)
turso := core_datastore.NewTursoClient(
    apiURL,
    cfg.TursoPlatformToken,
    &http.Client{},
)
```

You can then call:

```go
_, _ = turso.CreateInstance(ctx, req)
_, _ = turso.ListInstances(ctx)
_, _ = turso.GetInstance(ctx, "name")
_    = turso.DeleteInstance(ctx, "name")
```

---

## Initializing Additional Turso Databases

To open multiple Turso databases directly, call `ConnectToInstance` with each instance URL:

```go
dbA, err := core_datastore.ConnectToInstance("https://tenant-a.turso.sh/db.sqlite", "tokenA")
if err != nil {
    log.Fatalf("open A: %v", err)
}
defer dbA.Close()

dbB, err := core_datastore.ConnectToInstance("https://tenant-b.turso.sh/db.sqlite", "tokenB")
if err != nil {
    log.Fatalf("open B: %v", err)
}
defer dbB.Close()
```

After provisioning via the API:

```go
inst, err := manager.Turso().CreateInstance(ctx, core_datastore.CreateInstanceRequest{Name: "tenant-c", Region: "us-east-2"})
tenDB, _ := core_datastore.ConnectToInstance(inst.URL, cfg.TursoPlatformToken)
defer tenDB.Close()
```

---

## API Reference

### func ConnectToInstance(instanceURL, authToken string) (\*sql.DB, error)

* Opens a libsql database at `instanceURL?authToken=…`
* Runs `PRAGMA foreign_keys = ON;`
* Pings to verify connectivity

### func EnsureSystemUser(db \*sql.DB) error

* Looks up a user named `system-logging`
* If missing, inserts one with a condensed UUID

### func LogEvent(db \*sql.DB, app, event, desc, userID string) error

* Inserts into `logs` (app\_name, event\_type, description, timestamp, hstles\_user\_id)
* Defaults to system-logging user if `userID` is empty

### type TursoClient

```go
type TursoClient struct {
    baseURL   string
    authToken string
    http      *http.Client
}
```

Provides methods:

* `CreateInstance(ctx, req)`
* `ListInstances(ctx)`
* `GetInstance(ctx, name)`
* `DeleteInstance(ctx, name)`

### type Manager

```go
type Manager struct {
    CoreDB *sql.DB
    turso  *TursoClient
}
```

* `NewManager(cfg)` opens CoreDB and creates TursoClient
* `Close()` closes CoreDB
* `Turso()` returns the client

---

## Configuration

Your `core_config.CoreConfig` must provide:

* `TursoCoreDBURL`
* `TursoAuthToken`
* `TursoPlatformToken`

Loaded via environment variables or config file by `core_config.LoadCoreConfig()`.

---

## License

MIT © [hstles.com](https://hstles.com)
