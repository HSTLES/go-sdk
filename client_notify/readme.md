# client\_notify

A Go client library for sending emails via the HSTLES Notification Service.
Supports two integration styles:

1. **Instance-based**: create an `*EmailClient` with `NewClient(...)` and call its methods directly.
2. **Global wrappers**: call `Init(...)` once at startup, then use top-level functions everywhere without passing a client instance around.

---

## Installation

```bash
go get github.com/hstles/go-sdk/client_notify
```

---

## Package Contents

* **`types.go`**

  * Request and response structs:

    * `EmailResponse` – standard response for all endpoints
    * `WelcomeEmailRequest`, `SecurityCodeEmailRequest`, `ServiceAlertEmailRequest`, `LoginLinkEmailRequest`, `GenericEmailRequest`

* **`client.go`**

  * `type EmailClient` – holds `baseURL` and `http.Client`
  * `func NewClient(baseURL string) *EmailClient` – constructor
  * Internal helper:

    * `post(ctx, endpoint, payload) (*EmailResponse, error)` – marshals payload and handles POST to `/api/email/{endpoint}`
  * Instance methods:

    * `GetStatus(ctx) (*EmailResponse, error)`
    * `SendWelcomeEmail(ctx, to, userName) (*EmailResponse, error)`
    * `SendSecurityCodeEmail(ctx, to, code) (*EmailResponse, error)`
    * `SendServiceAlertEmail(ctx, to, title, message) (*EmailResponse, error)`
    * `SendLoginLinkEmail(ctx, to, link, name) (*EmailResponse, error)`
    * `SendGenericEmail(ctx, to, subject, message) (*EmailResponse, error)`

* **`wrappers.go`**

  * Package-level state:

    * `var Default *EmailClient`
  * Initialization:

    * `func Init(baseURL string)` – constructs the default client, prepending `https://` if needed
  * Helper to guard initialization:

    * `func ensure() error`
  * Top-level functions:

    * `GetStatus(ctx) (*EmailResponse, error)`
    * `SendWelcomeEmail(ctx, to, userName) (*EmailResponse, error)`
    * `SendSecurityCodeEmail(ctx, to, code) (*EmailResponse, error)`
    * `SendServiceAlertEmail(ctx, to, title, message) (*EmailResponse, error)`
    * `SendLoginLinkEmail(ctx, to, link, name) (*EmailResponse, error)`
    * `SendGenericEmail(ctx, to, subject, message) (*EmailResponse, error)`

---

## Usage

### 1) Instance-Based Client

Use when you want explicit control over multiple clients or dependency injection.

```go
package main

import (
    "context"
    "log"

    "github.com/hstles/go-sdk/client_notify"
)

func main() {
    client := client_notify.NewClient("https://notify.hstles.com")

    resp, err := client.SendLoginLinkEmail(
        context.Background(),
        "user@example.com",
        "https://app.hstles.com/login?token=abc123",
        "Alice",
    )
    if err != nil {
        log.Fatalf("failed to send login email: %v", err)
    }
    log.Printf("email sent successfully: %s", resp.Message)
}
```

---

### 2) Global Wrappers

Use for brevity when you only ever need one notification endpoint in your process.

```go
package worker

import (
    "context"
    "log"

    "github.com/hstles/go-sdk/client_notify"
)

func init() {
    client_notify.Init("notify.hstles.com")
}

func sendAlert(ctx context.Context, email, title, msg string) error {
    resp, err := client_notify.SendServiceAlertEmail(ctx, email, title, msg)
    if err != nil {
        return fmt.Errorf("alert email failed: %w", err)
    }
    log.Printf("alert sent: %s", resp.Message)
    return nil
}
```

Call any wrapper directly once `Init` is done:

```go
client_notify.GetStatus(ctx)
client_notify.SendWelcomeEmail(ctx, ...)
client_notify.SendSecurityCodeEmail(ctx, ...)
client_notify.SendLoginLinkEmail(ctx, ...)
client_notify.SendGenericEmail(ctx, ...)
```

---

## Error Handling

* `EmailResponse.Success == false` returns an error with server message
* Wrapper functions return an error if `Init` wasn’t called first

---

## License

MIT. See `LICENSE` for details.
