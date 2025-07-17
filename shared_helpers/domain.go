package shared_helpers

import (
	"fmt"
	"net"
	"os"
)

// GetServerDomain returns the domain or hostname of the server
func GetServerDomain() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("failed to get hostname: %v", err)
	}

	// Try to resolve the full domain
	addrs, err := net.LookupHost(hostname)
	if err != nil {
		return hostname, nil // Return just the hostname if full domain can't be resolved
	}

	// Return the first resolved address (which could be the full domain)
	return addrs[0], nil
}
