//go:build linux

package network

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
)

type LinuxManager struct {
	mu       sync.Mutex
	vpnState string
}

func getPlatformManager() Manager {
	return &LinuxManager{
		vpnState: "Disconnected",
	}
}

func (m *LinuxManager) CheckConnectivity() (bool, error) {
	// Execute real ping to 8.8.8.8
	// -c 1: count 1
	// -W 2: timeout 2 seconds
	cmd := exec.Command("ping", "-c", "1", "-W", "2", "8.8.8.8")
	if err := cmd.Run(); err != nil {
		return false, nil
	}
	return true, nil
}

func (m *LinuxManager) GetLatency(target string) (int64, error) {
	// Simple latency check via ping output parsing is complex in Go without a library.
	// For MVP, we will measure the time it takes to execute the command.

	// Retry up to 2 times to prevent flaky logs in WSL/WiFi
	for i := 0; i < 2; i++ {
		start := time.Now()
		// -c 1: count 1, -W 2: timeout 2s
		cmd := exec.Command("ping", "-c", "1", "-W", "2", target)
		err := cmd.Run()
		duration := time.Since(start).Milliseconds()

		if err == nil {
			return duration, nil
		}
		// Brief pause before retry
		time.Sleep(200 * time.Millisecond)
	}

	return 0, fmt.Errorf("ping failed")
}

func (m *LinuxManager) ListInterfaces() ([]string, error) {
	// Read from /sys/class/net
	files, err := os.ReadDir("/sys/class/net")
	if err != nil {
		return nil, err
	}
	var ifaces []string
	for _, f := range files {
		ifaces = append(ifaces, f.Name())
	}
	return ifaces, nil
}

func (m *LinuxManager) RestartInterface(name string) error {
	// 'ip link set <name> down'
	cmdDown := exec.Command("ip", "link", "set", name, "down")
	if err := cmdDown.Run(); err != nil {
		return err
	}

	time.Sleep(1 * time.Second)

	// 'ip link set <name> up'
	cmdUp := exec.Command("ip", "link", "set", name, "up")
	return cmdUp.Run()
}

func (m *LinuxManager) GetActiveWAN() (string, error) {
	// Real implementation: check routing table 'ip route show default'
	return "eth0", nil
}

func (m *LinuxManager) SwitchWAN(wanInterface string) error {
	// Real implementation: update routing table metric or use mwan3
	return nil
}

func (m *LinuxManager) EnableVPN() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	// Real implementation: 'wg-quick up wg0'
	// For now, we simulate the state change so the UI works
	m.vpnState = "Connected"
	return nil
}

func (m *LinuxManager) DisableVPN() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	// Real implementation: 'wg-quick down wg0'
	m.vpnState = "Disconnected"
	return nil
}

func (m *LinuxManager) GetVPNStatus() (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.vpnState == "" {
		return "Disconnected", nil
	}
	return m.vpnState, nil
}

func (m *LinuxManager) GetWifiInfo() (int, int, error) {
	// Real: Parse 'iw dev wlan0 link'
	return 6, 80, nil
}

func (m *LinuxManager) ScanWifiChannels() ([]WifiChannel, error) {
	// Real: 'iw dev wlan0 scan'
	return []WifiChannel{}, nil
}

func (m *LinuxManager) SetWifiChannel(channel int) error {
	// Real: 'hostapd_cli chan_switch' or 'iw config'
	return nil
}
