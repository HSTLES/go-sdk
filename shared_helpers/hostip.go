package shared_helpers

import (
	"log"
	"net"
)

// GetHostIP retrieves the server's IP address
func GetHostIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatalf("Error getting network interfaces: %v", err)
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String()
		}
	}
	return "Unknown IP"
}
