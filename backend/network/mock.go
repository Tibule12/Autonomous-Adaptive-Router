//go:build !linux

package network

import (
	"fmt"
	"math/rand"
	"time"
)

type MockManager struct{
	activeWAN      string
	vpnActive      bool
	currentChannel int
}

func getPlatformManager() Manager {
	fmt.Println("[SIMULATION] Initializing Mock Network Manager (Windows/Mac Mode)")
	return &MockManager{
		activeWAN:      "wan1_primary",
		vpnActive:      true,
		currentChannel: 6,
	}
}

func (m *MockManager) CheckConnectivity() (bool, error) {
	fmt.Println("[SIMULATION] Checking internet connectivity... OK")
	return true, nil
}

func (m *MockManager) GetLatency(target string) (int64, error) {
	// Simulate different performance for Primary vs Backup
	// Primary: Low latency, but occasional spikes (that trigger failover)
	// Backup: Higher consistent latency
	
	var latency int64
	if m.activeWAN == "wan1_primary" {
		// Mostly good (20-40ms), sometimes terrible (150ms+) to trigger Engine
		if rand.Intn(10) > 7 {
			latency = int64(rand.Intn(100) + 120) // Spike!
		} else {
			latency = int64(rand.Intn(20) + 20)
		}
	} else {
		// Backup is slower but stable (80-100ms)
		latency = int64(rand.Intn(20) + 80)
	}
	
	// fmt.Printf("[SIMULATION] Pinging %s via %s... %dms\n", target, m.activeWAN, latency)
	return latency, nil
}

func (m *MockManager) ListInterfaces() ([]string, error) {
	return []string{"sim_eth0", "sim_wlan0", "sim_wg0_vpn"}, nil
}

func (m *MockManager) RestartInterface(name string) error {
	fmt.Printf("[SIMULATION] ⚠️ RESTARTING INTERFACE: %s\n", name)
	time.Sleep(1 * time.Second) // Fake delay
	fmt.Printf("[SIMULATION] Interface %s is back up.\n", name)
	return nil
}

func (m *MockManager) GetActiveWAN() (string, error) {
	return m.activeWAN, nil
}

func (m *MockManager) SwitchWAN(wanInterface string) error {
	fmt.Printf("[SIMULATION] 🔀 Switching WAN to: %s\n", wanInterface)
	time.Sleep(2 * time.Second) // Simulation switching delay
	m.activeWAN = wanInterface
	fmt.Printf("[SIMULATION] New Active WAN: %s\n", m.activeWAN)
	return nil
}

func (m *MockManager) EnableVPN() error {
	fmt.Println("[SIMULATION] 🔒 Enabling VPN (WireGuard)... Success")
	m.vpnActive = true
	return nil
}

func (m *MockManager) DisableVPN() error {
	fmt.Println("[SIMULATION] 🔓 Disabling VPN... Success")
	m.vpnActive = false
	return nil
}

func (m *MockManager) GetVPNStatus() (string, error) {
	if m.vpnActive {
		return "Connected", nil
	}
	return "Disconnected", nil
}

// Wi-Fi Mock Implementation

func (m *MockManager) GetWifiInfo() (int, int, error) {
	// Simulate signal quality fluctuating
	// Allow it to drop low to trigger optimization
	quality := rand.Intn(30) + 70 // 70-100 normally
	
	// Occasionally simulate terrible congestion on Channel 6 (default)
	if m.currentChannel == 6 && rand.Intn(10) > 6 {
		quality = 30 // Bad quality
	}

	return m.currentChannel, quality, nil
}

func (m *MockManager) ScanWifiChannels() ([]WifiChannel, error) {
	fmt.Println("[SIMULATION] 📡 Scanning Wi-Fi Spectrum...")
	time.Sleep(1 * time.Second) // Fake scan delay
	
	// Simulate scan results (Score 0-100, higher is better)
	return []WifiChannel{
		{Channel: 1, Score: 85},
		{Channel: 6, Score: 40}, // Congested!
		{Channel: 11, Score: 92}, // Good candidate
		{Channel: 36, Score: 95},
		{Channel: 161, Score: 98},
	}, nil
}

func (m *MockManager) SetWifiChannel(channel int) error {
	fmt.Printf("[SIMULATION] ⚙️ Tuning Wi-Fi Radio to Channel %d\n", channel)
	time.Sleep(2 * time.Second) // Fake switching delay
	m.currentChannel = channel
	fmt.Printf("[SIMULATION] ✅ Wi-Fi now operating on Channel %d\n", channel)
	return nil
}
