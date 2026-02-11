//go:build linux

package network

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

type LinuxManager struct{}

func getPlatformManager() Manager {
	return &LinuxManager{}
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
	start := time.Now()
	cmd := exec.Command("ping", "-c", "1", "-W", "1", target)
	err := cmd.Run()
	duration := time.Since(start).Milliseconds()

	if err != nil {
		return 0, fmt.Errorf("ping failed")
	}
	return duration, nil
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
	// Real implementation: 'wg-quick up wg0'
	return nil
}

func (m *LinuxManager) DisableVPN() error {
	// Real implementation: 'wg-quick down wg0'
	return nil
}

func (m *LinuxManager) GetVPNStatus() (string, error) {
	return "Active", nil
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
