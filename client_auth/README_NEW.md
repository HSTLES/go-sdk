# client\_auth

A Go client library and HTTP handler wrappers for interacting with the HSTLES Auth Service.
Supports two integration styles:

1. **Instance-based**: create a `*Client` with `NewClient(...)`, inject it into your handlers and call its methods.
2. **Global wrappers**: call `Init(...)` once at startup, then use top-level functions everywhere without passing a client around.

---

## Installation

```bash
go get github.com/hstles/go-sdk/client_auth
```

---

## Package Contents

* **`client.go`**

  * Defines `type Client`
  * Constructor: `NewClient(baseURL string) *Client`
  * Instance methods:

    **Session Management:**
    * `ValidateSession(ctx, cookies) (SessionResponse, int, error)`
    * `DeleteSession(ctx, cookies, sessionID) (DeleteSessionResponse, int, error)`
    * `DeleteAllSessions(ctx, cookies) (DeleteSessionResponse, int, error)`

    **Legacy 2FA:**
    * `Get2FAStatus(ctx, cookies) (TwoFAStatusResponse, int, error)`
    * `CheckPendingSession(ctx, cookies) (PendingSessionResponse, int, error)`
    * `Delete2FA(ctx, cookies) (DeleteSessionResponse, int, error)`
    * `LockoutUser(ctx, cookies, userID) (LockoutResponse, int, error)`

    **New 2FA API:**
    * `Configure2FA(ctx, cookies, req) (Configure2FAResponse, int, error)`
    * `Verify2FA(ctx, cookies, req) (Verify2FAResponse, int, error)`
    * `Reset2FA(ctx, cookies, req) (Reset2FAResponse, int, error)`
    * `GenerateBackupCodes(ctx, cookies, req) (GenerateBackupCodesResponse, int, error)`
    * `CheckTrustedDevice(ctx, cookies) (CheckTrustedDeviceResponse, int, error)`
    * `GetLockoutStatus(ctx, cookies) (LockoutStatusResponse, int, error)`
    * `ClearLockout(ctx, cookies) (ClearLockoutResponse, int, error)`
    * `InitiateRecovery(ctx, req) (InitiateRecoveryResponse, int, error)`
    * `VerifyRecoveryCode(ctx, req) (VerifyRecoveryCodeResponse, int, error)`

    **Auth Flow:**
    * `AuthFlow(ctx, cookies, provider, next) (string, int, error)`
    * `Auth(ctx, cookies, provider, next, form) (string, int, error)`
    * `AuthCallback(ctx, cookies, provider, query) (string, int, error)`

* **`wrappers.go`**

  * `var Default *Client`
  * `func Init(baseURL string)`
  * Top-level helper functions that forward to `Default` after a nil-check.

* **`handlers.go`**

  * HTTP handlers (for Gorilla Mux) that proxy incoming requests to an injected `*Client`, including:

    **Session Management:**
    * `ValidateSessionHandler(*Client)`
    * `DeleteSessionHandler(*Client)`
    * `DeleteAllSessionsHandler(*Client)`

    **Legacy 2FA:**
    * `Get2FAStatusHandler(*Client)`
    * `CheckPendingSessionHandler(*Client)`
    * `Delete2FAHandler(*Client)`
    * `LockoutUserHandler(*Client)`

    **New 2FA API:**
    * `Configure2FAHandler(*Client)`
    * `Verify2FAHandler(*Client)`
    * `Reset2FAHandler(*Client)`
    * `GenerateBackupCodesHandler(*Client)`
    * `CheckTrustedDeviceHandler(*Client)`
    * `GetLockoutStatusHandler(*Client)`
    * `ClearLockoutHandler(*Client)`
    * `InitiateRecoveryHandler(*Client)`
    * `VerifyRecoveryCodeHandler(*Client)`

    **Auth Flow:**
    * `AuthFlowHandler(*Client)`
    * `AuthHandler(*Client)`
    * `AuthCallbackHandler(*Client)`

* **`types.go`**

  * Struct definitions for all request and response payloads:
    * `SessionResponse`, `DeleteSessionResponse`, `TwoFAStatusResponse`, `PendingSessionResponse`
    * `LockoutRequest`, `LockoutResponse`, `LockoutStatusResponse`, `ClearLockoutResponse`
    * `Configure2FARequest`, `Configure2FAResponse`
    * `Verify2FARequest`, `Verify2FAResponse`
    * `Reset2FARequest`, `Reset2FAResponse`
    * `GenerateBackupCodesRequest`, `GenerateBackupCodesResponse`
    * `CheckTrustedDeviceRequest`, `CheckTrustedDeviceResponse`
    * `InitiateRecoveryRequest`, `InitiateRecoveryResponse`
    * `VerifyRecoveryCodeRequest`, `VerifyRecoveryCodeResponse`

---

## Usage

### 1) Instance-Based Client & Handlers

Ideal for server-side proxies (e.g. HTMX), explicit wiring, and testability.

```go
package main

import (
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/hstles/go-sdk/client_auth"
)

func main() {
    // 1) Construct a client pointing at your auth service
    authClient := client_auth.NewClient("https://auth.hstles.com")

    // 2) Set up Gorilla Mux routes, injecting authClient
    r := mux.NewRouter()
    
    // Session management
    r.Handle("/api/session", client_auth.ValidateSessionHandler(authClient)).Methods("GET")
    r.Handle("/api/session", client_auth.DeleteSessionHandler(authClient)).Methods("POST")
    r.Handle("/api/session", client_auth.DeleteAllSessionsHandler(authClient)).Methods("DELETE")

    // Legacy 2FA endpoints
    r.Handle("/api/2fa", client_auth.Get2FAStatusHandler(authClient)).Methods("GET")
    r.Handle("/api/2fa", client_auth.CheckPendingSessionHandler(authClient)).Methods("POST")
    r.Handle("/api/2fa", client_auth.Delete2FAHandler(authClient)).Methods("DELETE")

    // New 2FA API endpoints
    r.Handle("/api/2fa/configure", client_auth.Configure2FAHandler(authClient)).Methods("POST")
    r.Handle("/api/2fa/verify", client_auth.Verify2FAHandler(authClient)).Methods("POST")
    r.Handle("/api/2fa/reset", client_auth.Reset2FAHandler(authClient)).Methods("POST")
    r.Handle("/api/2fa/backup-codes", client_auth.GenerateBackupCodesHandler(authClient)).Methods("POST")
    r.Handle("/api/2fa/trusted-device", client_auth.CheckTrustedDeviceHandler(authClient)).Methods("GET")
    r.Handle("/api/2fa/lockout", client_auth.GetLockoutStatusHandler(authClient)).Methods("GET")
    r.Handle("/api/2fa/lockout", client_auth.ClearLockoutHandler(authClient)).Methods("DELETE")
    r.Handle("/api/2fa/recovery", client_auth.InitiateRecoveryHandler(authClient)).Methods("POST")
    r.Handle("/api/2fa/recovery/verify", client_auth.VerifyRecoveryCodeHandler(authClient)).Methods("POST")

    // Auth flow
    r.Handle("/auth", client_auth.AuthFlowHandler(authClient)).Methods("GET")
    r.Handle("/auth/{provider}", client_auth.AuthHandler(authClient)).Methods("POST")
    r.Handle("/auth/{provider}/callback", client_auth.AuthCallbackHandler(authClient)).Methods("GET")

    log.Fatal(http.ListenAndServe(":8080", r))
}
```

This **server-side proxy** approach keeps all HTMX/AJAX calls on the same origin (`login.hstles.com`), avoiding CORS and SameSite cookie issues.

---

### 2) Global Wrappers

Great for non-HTTP code (background jobs, CLI tools, business logic) when you'd rather not pass a client instance everywhere.

```go
package worker

import (
    "context"

    "github.com/hstles/go-sdk/client_auth"
)

func init() {
    client_auth.Init("https://auth.hstles.com")
}

func configure2FAForUser(ctx context.Context, cookies []*http.Cookie, secret, code, backupCodes string) error {
    req := client_auth.Configure2FARequest{
        Secret:      secret,
        Code:        code,
        BackupCodes: backupCodes,
    }
    
    resp, status, err := client_auth.Configure2FA(ctx, cookies, req)
    if err != nil {
        return fmt.Errorf("failed to configure 2FA (status %d): %w", status, err)
    }
    
    if !resp.Success {
        return fmt.Errorf("2FA configuration failed: %s", resp.Error)
    }
    
    return nil
}

func verify2FACode(ctx context.Context, cookies []*http.Cookie, code string, rememberDevice bool) error {
    req := client_auth.Verify2FARequest{
        Code:           code,
        RememberDevice: rememberDevice,
    }
    
    resp, status, err := client_auth.Verify2FA(ctx, cookies, req)
    if err != nil {
        return fmt.Errorf("failed to verify 2FA (status %d): %w", status, err)
    }
    
    if !resp.Success {
        return fmt.Errorf("2FA verification failed: %s", resp.Error)
    }
    
    return nil
}
```

---

## 2FA API Examples

### Configure 2FA

```go
req := client_auth.Configure2FARequest{
    Secret:      "JBSWY3DPEHPK3PXP",
    Code:        "123456",
    BackupCodes: "abc123,def456,ghi789",
}

resp, status, err := authClient.Configure2FA(ctx, cookies, req)
if err != nil {
    log.Printf("Error: %v", err)
    return
}

if resp.Success {
    log.Printf("2FA configured: %s", resp.Message)
} else {
    log.Printf("Failed: %s", resp.Error)
}
```

### Verify 2FA

```go
req := client_auth.Verify2FARequest{
    Code:           "654321",
    NextURL:        "https://account.hstles.com/dashboard",
    RememberDevice: true,
}

resp, status, err := authClient.Verify2FA(ctx, cookies, req)
if err != nil {
    log.Printf("Error: %v", err)
    return
}

if resp.Success {
    log.Printf("2FA verified, redirect to: %s", resp.RedirectURL)
} else if resp.Locked {
    log.Printf("Account locked for %d seconds", resp.LockDuration)
} else {
    log.Printf("Verification failed: %s", resp.Error)
}
```

### Generate New Backup Codes

```go
req := client_auth.GenerateBackupCodesRequest{
    Code: "123456", // Current TOTP code for verification
}

resp, status, err := authClient.GenerateBackupCodes(ctx, cookies, req)
if err != nil {
    log.Printf("Error: %v", err)
    return
}

if resp.Success {
    log.Printf("New backup codes: %v", resp.BackupCodes)
} else {
    log.Printf("Failed: %s", resp.Error)
}
```

### Check Lockout Status

```go
resp, status, err := authClient.GetLockoutStatus(ctx, cookies)
if err != nil {
    log.Printf("Error: %v", err)
    return
}

if resp.IsLocked {
    log.Printf("Account locked: %s (remaining: %d seconds)", resp.LockMessage, resp.RemainingTime)
} else {
    log.Printf("Account not locked. Attempts: %d/%d", resp.AttemptCount, resp.MaxAttempts)
}
```

### Initiate Account Recovery

```go
req := client_auth.InitiateRecoveryRequest{
    Email: "user@example.com",
}

resp, status, err := authClient.InitiateRecovery(ctx, req)
if err != nil {
    log.Printf("Error: %v", err)
    return
}

log.Printf("Recovery initiated: %s", resp.Message)
```

---

## HTMX Proxy Example

```html
<!-- On login.hstles.com -->
<button
  hx-get="/api/session"
  hx-swap="outerHTML"
>
  Check Session
</button>

<form hx-post="/api/2fa/verify" hx-swap="outerHTML">
  <input type="text" name="code" placeholder="Enter 2FA code">
  <input type="checkbox" name="remember_device" value="true"> Remember this device
  <button type="submit">Verify</button>
</form>
```

1. HTMX issues a same-origin request to your API endpoints.
2. Your Go handler proxies the request to `auth.hstles.com` using the `client_auth` library.
3. Auth responses are returned directly—no CORS or cross-site cookies needed.

---

## License

Distributed under the MIT License. See `LICENSE` for details.
