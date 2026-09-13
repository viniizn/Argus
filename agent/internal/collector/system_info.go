package collector

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"time"

	"github.com/viniizn/Argus/agent/internal/model"
)

const agentVersion = "0.1.0"

func SystemInfo() (model.SystemInfo, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return model.SystemInfo{}, fmt.Errorf("collector: read hostname: %w", err)
	}

	ip, err := localIpv4()
	if err != nil {
		return model.SystemInfo{}, fmt.Errorf("collector: resolve local IP: %w", err)
	}

	return model.SystemInfo{
		Hostname: hostname,
		OS: runtime.GOOS,
		Architecture: runtime.GOARCH,
		IPAddress: ip,
		AgentVersion: agentVersion,
		CollectedAt: time.Now().UTC(),
	}, nil
}

func localIpv4() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if (err != nil) {
		return "", fmt.Errorf("list interface addresses: %w", err)
	}

	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if !ok || ipNet.IP.IsLoopback() {
			continue
		}
		if ipv4 := ipNet.IP.To4(); ipv4 != nil {
			return ipv4.String(), nil
		}
	}

	return "", fmt.Errorf("no non-loopback IPv4 addres found")
}