package network

type WifiChannel struct {
	Channel int
	Score   int // 0-100 (100 is best quality/least execution)
}

type NetworkMetrics struct {
	LatencyMs  int64
	PacketLoss float64 // Percentage 0.0 - 100.0
	JitterMs   int64
}

// Manager defines the interface for router network operations.
// It abstracts away the OS-specific commands.
type Manager interface {
	// Health Checks
	CheckConnectivity() (bool, error)
	GetNetworkMetrics(target string) (NetworkMetrics, error)

	// Interface Management
	ListInterfaces() ([]string, error)
	RestartInterface(name string) error

	// Multi-WAN Management
	GetActiveWAN() (string, error)
	SwitchWAN(wanInterface string) error

	// VPN Management
	EnableVPN() error
	DisableVPN() error
	GetVPNStatus() (string, error)

	// Wi-Fi Management
	SetSimulatedPacketLoss(enabled bool)
	SetSimulatedLoad(requestsPerSecond int) // New: For multi-device stress testing
	GetWifiInfo() (int, int, error)         // Returns (currentChannel, signalQuality 0-100)
	ScanWifiChannels() ([]WifiChannel, error)
	SetWifiChannel(channel int) error

	// Chaos Engineering (Simulation)
	SetSimulatedLag(enabled bool)
	SetSimulatedInterference(enabled bool)

	// Traffic Analysis (AI Simulation)
	SetSimulatedTraffic(trafficType string)
	GetTrafficAnalysis() (string, error)
}

// NewManager creates a platform-specific network manager.
// The actual implementation returned depends on the build tags (OS).
func NewManager() Manager {
	return getPlatformManager()
}
