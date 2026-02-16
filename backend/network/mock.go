//go:build !linux

package network

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/TMtshwelo/aar/pkg/storage"
)

type MockManager struct {
	activeWAN                 string
	vpnActive                 bool
	currentChannel            int
	simulatedLag              bool
	simulatedLoss             bool
	simulatedWiFiInterference bool
	trafficType               string
	currentLoad               int      // Requests per second
	connectedDevices          []Device // Mock devices for UI testing
}

func getPlatformManager() Manager {
	fmt.Println("[SIMULATION] Initializing Mock Network Manager (Windows/Mac Mode)")
	devList := []Device{
		{Name: "Dad's iPhone", IP: "192.168.1.10", MAC: "00:1A:2B:3C:4D:5E"},
		{Name: "Living Room TV", IP: "192.168.1.15", MAC: "11:22:33:44:55:66"},
		{Name: "PlayStation 5", IP: "192.168.1.20", MAC: "AA:BB:CC:DD:EE:FF"}, // Check persistence for block
		{Name: "Unknown Device", IP: "192.168.1.50", MAC: "CA:FE:BA:BE:00:11"},
	}

	// Dynamic Persistence Check
	for i := range devList {
		devList[i].IsBlocked = storage.IsBlocked(devList[i].MAC)
	}

	return &MockManager{
		activeWAN:                 "wan1_primary",
		vpnActive:                 true,
		currentChannel:            6,
		simulatedLag:              false,
		simulatedLoss:             false,
		simulatedWiFiInterference: false,
		connectedDevices:          devList, // Injected persistent state
		trafficType:               "Default",
		currentLoad:               0,
	}
}

func (m *MockManager) CheckConnectivity() (bool, error) {
	fmt.Println("[SIMULATION] Checking internet connectivity... OK")
	return true, nil
}

func (m *MockManager) GetNetworkMetrics(target string) (NetworkMetrics, error) {
	// Logic to return bad metrics if chaos mode is on
	// Base metrics
	metrics := NetworkMetrics{
		LatencyMs:  45,
		PacketLoss: 0.0,
		JitterMs:   5,
	}

	// 1. Apply High Traffic (DDoS/Stress) Effects
	// 50 reps/sec is the threshold in this simulation
	if m.currentLoad > 0 {
		// Log scale impact: more requests = exponential latency
		metrics.LatencyMs += int64(m.currentLoad * 2)
		metrics.JitterMs += int64(m.currentLoad / 2)

		if m.currentLoad > 50 {
			metrics.PacketLoss += 5.0
			m.trafficType = "High Load"
		}
		if m.currentLoad > 100 {
			metrics.PacketLoss += 15.0
			m.trafficType = "Congestion"
		}
	}

	if m.simulatedLoss {
		metrics.PacketLoss = 25.0
	}

	if m.simulatedLag {
		metrics.LatencyMs = 500
	}

	return metrics, nil
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

// --- Device Management (SIMULATED) ---
func (m *MockManager) GetConnectedDevices() ([]Device, error) {
	return m.connectedDevices, nil
}

func (m *MockManager) BlockDevice(mac string) error {
	fmt.Printf("[MOCK] 🚫 Device Blocked: %s\n", mac)
	// Persist for next restart
	storage.AddBlockedMAC(mac)

	// Update Runtime Memory
	for i, dev := range m.connectedDevices {
		if dev.MAC == mac {
			m.connectedDevices[i].IsBlocked = true
		}
	}
	return nil
}

func (m *MockManager) UnblockDevice(mac string) error {
	fmt.Printf("[MOCK] ✅ Device Unblocked: %s\n", mac)
	// Persist for next restart
	storage.RemoveBlockedMAC(mac)

	// Update Runtime Memory
	for i, dev := range m.connectedDevices {
		if dev.MAC == mac {
			m.connectedDevices[i].IsBlocked = false
		}
	}
	return nil
}

// Wi-Fi Mock Implementation

// Wi-Fi Mock Implementation

func (m *MockManager) GetWifiInfo() (int, int, error) {
	// Simulate signal quality fluctuating
	// Allow it to drop low to trigger optimization
	quality := rand.Intn(30) + 70 // 70-100 normally

	// Chaos: Massive interference forces quality down
	if m.simulatedWiFiInterference {
		quality = 25 // Critical interference
	}

	// Occasionally simulate terrible congestion on Channel 6 (default)
	if m.currentChannel == 6 && rand.Intn(10) > 6 {
		quality = 30 // Bad quality
	}

	return m.currentChannel, quality, nil
}

func (m *MockManager) ScanWifiChannels() ([]WifiChannel, error) {
	fmt.Println("[SIMULATION] 📡 Scanning Wi-Fi Spectrum...")
	time.Sleep(1 * time.Second) // Fake scan delay

	// Base Candidates
	candidates := []WifiChannel{
		{Channel: 1, Score: 85},
		{Channel: 6, Score: 40},  // Congested!
		{Channel: 11, Score: 92}, // Good candidate
		{Channel: 36, Score: 95},
		{Channel: 161, Score: 98},
	}

	// Dynamic Adjustment based on simulation state
	for i := range candidates {
		// If Jamming is active, the CURRENT channel is garbage
		if m.simulatedWiFiInterference && candidates[i].Channel == m.currentChannel {
			candidates[i].Score = 20 // Crushed by interference
		}
	}

	return candidates, nil
}

func (m *MockManager) SetWifiChannel(channel int) error {
	fmt.Printf("[SIMULATION] ⚙️ Tuning Wi-Fi Radio to Channel %d\n", channel)
	time.Sleep(2 * time.Second) // Fake switching delay
	m.currentChannel = channel
	fmt.Printf("[SIMULATION] ✅ Wi-Fi now operating on Channel %d\n", channel)
	return nil
}

func (m *MockManager) SetSimulatedLoad(requestsPerSecond int) {
	m.currentLoad = requestsPerSecond
	if requestsPerSecond > 10 {
		fmt.Printf("[SIMULATION] ⚠️ High Traffic Detected: %d req/sec\n", requestsPerSecond)
	}

	// Auto-Detect Gaming Traffic based on "signature" if simulated
	// We'll simulate this by saying if traffic type is explicitly set to Gaming, we enforce it.
	// In a real router, this would be packet inspection (DPI).
}

func (m *MockManager) SetSimulatedTraffic(trafficType string) {
	m.trafficType = trafficType
	fmt.Printf("[SIMULATION] Traffic Pattern Changed: %s\n", trafficType)
}

func (m *MockManager) SetSimulatedInterference(enabled bool) {
	m.simulatedWiFiInterference = enabled
	fmt.Printf("[SIMULATION] Toggling Fake Wi-Fi Interference: %v\n", enabled)
}

func (m *MockManager) SetSimulatedLag(enabled bool) {
	m.simulatedLag = enabled
	fmt.Printf("[SIMULATION] Toggling Fake Lag: %v\n", enabled)
}

func (m *MockManager) SetSimulatedPacketLoss(enabled bool) {
	m.simulatedLoss = enabled
	fmt.Printf("[SIMULATION] Toggling Fake Packet Loss: %v\n", enabled)
}

func (m *MockManager) GetTrafficAnalysis() (string, error) {
	return m.trafficType, nil
}
