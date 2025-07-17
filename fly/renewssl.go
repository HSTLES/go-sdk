// fly/renewssl.go
package fly

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

const graphqlEndpoint = "https://api.fly.io/graphql"

// httpClient with dial/TLS/header timeouts so nothing hangs indefinitely.
var httpClient = &http.Client{
	Transport: &http.Transport{
		DialContext:           (&net.Dialer{Timeout: 5 * time.Second}).DialContext,
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
}

type gqlRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type certQuery struct {
	Data struct {
		App struct {
			Certificate struct {
				Issued struct {
					Nodes []struct {
						ExpiresAt time.Time `json:"expiresAt"`
					} `json:"nodes"`
				} `json:"issued"`
			} `json:"certificate"`
		} `json:"app"`
	} `json:"data"`
}

func doGraphQL(ctx context.Context, req gqlRequest, out interface{}) error {
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}
	httpReq, err := http.NewRequestWithContext(ctx, "POST", graphqlEndpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("http do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}

func getExpiry(ctx context.Context, app, hostname string) (time.Time, error) {
	// per-call timeout
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req := gqlRequest{
		Query: `
query($app:String!,$hostname:String!){
  app(name:$app){
    certificate(hostname:$hostname){
      issued { nodes { expiresAt } }
    }
  }
}`,
		Variables: map[string]interface{}{"app": app, "hostname": hostname},
	}

	var out certQuery
	if err := doGraphQL(ctx, req, &out); err != nil {
		return time.Time{}, fmt.Errorf("fetch expiry %s: %w", hostname, err)
	}
	nodes := out.Data.App.Certificate.Issued.Nodes
	if len(nodes) == 0 {
		return time.Time{}, fmt.Errorf("no certificate for %q", hostname)
	}
	return nodes[0].ExpiresAt, nil
}

func renewCert(ctx context.Context, app, hostname string) error {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	req := gqlRequest{
		Query: `
mutation($app:String!,$hostname:String!){
  renewCertificate(app:$app,hostname:$hostname){ success }
}`,
		Variables: map[string]interface{}{"app": app, "hostname": hostname},
	}
	var sink struct{ Data json.RawMessage }
	if err := doGraphQL(ctx, req, &sink); err != nil {
		return fmt.Errorf("renewal for %s: %w", hostname, err)
	}
	return nil
}

// StartAutoRenew loops forever. For each host:
//   - on expiry-fetch or renew error it calls onError(err) and CONTINUES to the next host,
//   - otherwise logs success/failure on stdout.
//
// After all hosts are handled it waits checkPeriod, then repeats.
// Only exits if ctx is cancelled.
func StartAutoRenew(ctx context.Context, onError func(error)) {
	log.Println("🔄 Starting SSL auto-renew…")
	app := strings.TrimSpace(os.Getenv("FLY_APP_NAME"))
	hosts := strings.Split(os.Getenv("FLY_HOSTNAMES"), ",")

	const (
		threshold   = 7 * 24 * time.Hour // renew if <7d
		checkPeriod = 72 * time.Hour     // run every 3d
	)
	ticker := time.NewTicker(checkPeriod)
	defer ticker.Stop()

	log.Printf("App %q hosts: %v (threshold %v, interval %v)",
		app, hosts, threshold, checkPeriod,
	)

	for {
		log.Println("• Beginning certificate check cycle")
		for _, raw := range hosts {
			h := strings.TrimSpace(raw)
			log.Printf("  – Checking %s", h)

			exp, err := getExpiry(ctx, app, h)
			if err != nil {
				onError(fmt.Errorf("expiry error for %s: %w", h, err))
				continue
			}

			until := time.Until(exp)
			if until < threshold {
				log.Printf("  ❗ %s expires in %v → renewing", h, until)
				if err := renewCert(ctx, app, h); err != nil {
					onError(fmt.Errorf("renewal error for %s: %w", h, err))
				} else {
					log.Printf("  ✅ renewed %s", h)
				}
			} else {
				log.Printf("  ✅ %s valid for %v", h, until)
			}
		}

		select {
		case <-ctx.Done():
			log.Println("🔚 Context done; stopping auto-renew")
			return
		case <-ticker.C:
			// next cycle
		}
	}
}
