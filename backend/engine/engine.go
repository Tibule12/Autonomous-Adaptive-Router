package engine

import (
	"fmt"
	"time"

	"github.com/TMtshwelo/aar/network"
)

// Config holds the thresholds for auto-recovery triggering
type Config struct {
	MaxLatencyMs      int64
	CheckInterval     time.Duration
	StabilityDuration int // How many checks must fail before action
}

// AutoPilot is the brain of the router.
// It monitors metrics and executes recovery actions via the Network Manager.
type AutoPilot struct {
	manager       network.Manager
	config        Config
	failureStreak int
}

func NewAutoPilot(mgr network.Manager) *AutoPilot {
	return &AutoPilot{
		manager: mgr,
		config: Config{
			MaxLatencyMs:      200, // Threshold: 200ms (tuned for WSL/WiFi jitter)
			CheckInterval:     3 * time.Second,
			StabilityDuration: 3, // React after 3 bad checks
		},
	}
}

// Start begins the monitoring loop
func (ap *AutoPilot) Start() {
	fmt.Println("[ENGINE] Auto-Pilot Engaged. Monitoring network health...")
	go ap.controlLoop()
}

func (ap *AutoPilot) controlLoop() {
	ticker := time.NewTicker(ap.config.CheckInterval)
	wifiTicker := time.NewTicker(10 * time.Second) // Check Wi-Fi less frequently

	for {
		select {
		case <-ticker.C:
			ap.evaluateNetwork()
		case <-wifiTicker.C:
			ap.optimizeWifi()
		}
	}
}

func (ap *AutoPilot) evaluateNetwork() {
	// 1. Monitor Latency
	latency, err := ap.manager.GetLatency("8.8.8.8")

	// Handling Total Network Failure (Ping Failed)
	if err != nil {
		fmt.Printf("[ENGINE] ❌ CRITICAL: Internet Unreachable (Ping Failed)\n")
		// Treat this as max latency to force a failure streak increment
		latency = 9999
	} else {
		// Valid measurement
		// fmt.Printf("[ENGINE] Latency: %dms\n", latency)
	}

	activeWAN, _ := ap.manager.GetActiveWAN()

	// 2. Decide
	if latency > ap.config.MaxLatencyMs {
		ap.failureStreak++
		fmt.Printf("[ENGINE] ⚠️ Abnormal Network State on %s: %dms (Streak: %d/%d)\n",
			activeWAN, latency, ap.failureStreak, ap.config.StabilityDuration)
	} else {
		// Recovery Logic: If we were failing, and now we are good, announce it loudly.
		if ap.failureStreak > 0 {
			fmt.Printf("[ENGINE] ✅ CONNECTION RESTORED. Latency: %dms. Reseting failure counter.\n", latency)
		}
		ap.failureStreak = 0
	}

	// 3. Act
	if ap.failureStreak >= ap.config.StabilityDuration {
		ap.triggerRecovery(latency, activeWAN)
		ap.failureStreak = 0 // Reset after action
	}
}

func (ap *AutoPilot) triggerRecovery(currentLatency int64, activeWAN string) {
	fmt.Println("------------------------------------------------")
	fmt.Printf("[ENGINE] 🚨 ACTION REQUIRED | Latency %dms exceeded threshold\n", currentLatency)

	if activeWAN == "wan1_primary" {
		fmt.Println("[ENGINE] 🔄 FAILOVER STRATEGY: Switching to Backup WAN...")
		err := ap.manager.SwitchWAN("wan2_backup")
		if err != nil {
			fmt.Printf("[ENGINE] Failover Failed: %v\n", err)
		} else {
			fmt.Println("[ENGINE] ✨ Switched to Secondary WAN (Backup)")
		}
	} else {
		fmt.Println("[ENGINE] 🔄 RECOVERY STRATEGY: Already on Backup. Restarting Interface...")
		err := ap.manager.RestartInterface(activeWAN)
		if err != nil {
			fmt.Printf("[ENGINE] Interface Restart Failed: %v\n", err)
		} else {
			fmt.Println("[ENGINE] ✨ Interface Restarted")
		}
	}
	fmt.Println("------------------------------------------------")
}

func (ap *AutoPilot) optimizeWifi() {
	// 1. Check current health
	currentChannel, quality, err := ap.manager.GetWifiInfo()
	if err != nil {
		return
	}

	// 2. Decide: If quality is bad (<40), trigger scan
	if quality < 40 {
		fmt.Printf("[ENGINE] ⚠️ Wi-Fi Quality degraded on Channel %d (Score: %d/100)\n", currentChannel, quality)
		fmt.Println("[ENGINE] 📡 Initiating Intelligent Channel Optimization...")

		bestChannel := ap.findBestChannel()
		if bestChannel != -1 && bestChannel != currentChannel {
			err := ap.manager.SetWifiChannel(bestChannel)
			if err != nil {
				fmt.Printf("[ENGINE] Wi-Fi Optimization Failed: %v\n", err)
			} else {
				fmt.Printf("[ENGINE] ✨ Wi-Fi Optimized! Switched Intelligently to Channel %d\n", bestChannel)
			}
		} else {
			fmt.Println("[ENGINE] ℹ️ Current channel is still the best available option.")
		}
	}
}

func (ap *AutoPilot) findBestChannel() int {
	candidates, err := ap.manager.ScanWifiChannels()
	if err != nil {
		fmt.Printf("[ENGINE] Scan failed: %v\n", err)
		return -1
	}

	bestScore := -1
	bestCh := -1

	for _, c := range candidates {
		if c.Score > bestScore {
			bestScore = c.Score
			bestCh = c.Channel
		}
	}
	return bestCh
}
