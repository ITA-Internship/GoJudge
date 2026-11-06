//go:build linux

package trusted

import (
	"net"
	"os"
	"testing"
	"time"
)

func TestNetworkEgress(t *testing.T) {
	targets := []string{
		"1.1.1.1:443",    // public IP, no DNS needed
		"8.8.8.8:53",     // public DNS server over TCP
		"google.com:443", // needs DNS resolution first
	}
	// Kubernetes API server inside the cluster (injected into every Pod)
	if host := os.Getenv("KUBERNETES_SERVICE_HOST"); host != "" {
		targets = append(targets, net.JoinHostPort(host, os.Getenv("KUBERNETES_SERVICE_PORT")))
	}

	for _, addr := range targets {
		start := time.Now()
		conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
		elapsed := time.Since(start).Round(time.Millisecond)
		if err == nil {
			t.Errorf("outbound TCP to %s succeeded in %v (remote %s)", addr, elapsed, conn.RemoteAddr())
			_ = conn.Close()
			continue
		}
		t.Logf("outbound TCP to %s blocked after %v: %v", addr, elapsed, err)
	}
}
