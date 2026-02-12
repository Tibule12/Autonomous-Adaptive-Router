package engine

import (
	"time"

	"github.com/TMtshwelo/aar/network"
	"github.com/TMtshwelo/aar/pkg/logger"
)

// Config holds the thresholds for auto-recovery triggering
type Config struct {
	MaxLatencyMs      int64
	MaxPacketLoss     float64
	CheckInterval     time.Duration
	StabilityDuration int // How many checks must fail before action
}

// AutoPilot is the brain of the router.
// It monitors metrics and executes recovery actions via the Network Manager.
type AutoPilot struct {
	manager        network.Manager
	config         Config
	failureStreak  int
	channelScores  map[int]int // AI Memory: Channel -> Score (0-100)
	currentTraffic string      // AI Memory: Current Traffic Type
	learningRate   float64     // How fast the AI adapts (0.0 - 1.0)
}

func NewAutoPilot(mgr network.Manager) *AutoPilot {
	ap := &AutoPilot{
		manager: mgr,
		config: Config{
			MaxLatencyMs:      150,             // Stricter Threshold: 150ms for faster reaction
			MaxPacketLoss:     5.0,             // Stricter Threshold: 5% packet loss
			CheckInterval:     2 * time.Second, // Check more frequently
			StabilityDuration: 3,
		},
		channelScores: map[int]int{
			6:   50, // Start neutral
			161: 50, // Start neutral
			1:   50,
			11:  50,
		},
		learningRate: 0.1,
	}

	// Load long-term memory immediately
	ap.LoadMemory()
	return ap
}

// Start begins the monitoring loop
func (ap *AutoPilot) Start() {
	logger.Info("Auto-Pilot Engaged. Monitoring network health...")
	go ap.controlLoop()
}

func (ap *AutoPilot) controlLoop() {
	ticker := time.NewTicker(ap.config.CheckInterval)
	wifiTicker := time.NewTicker(5 * time.Second) // Check Wi-Fi more frequently

	for {
		select {
		case <-ticker.C:
			// Continuous Learning Loop
			ap.learnAndEvaluate()
		case <-wifiTicker.C:
			// Proactive Optimization
			ap.optimizeWifi()
		}
	}
}

func (ap *AutoPilot) learnAndEvaluate() {
	// 1. Gather Sensory Data (Metrics)
	metrics, err := ap.manager.GetNetworkMetrics("8.8.8.8")
	activeWAN, _ := ap.manager.GetActiveWAN()
	currentChannel, _, _ := ap.manager.GetWifiInfo()

	// Handling Total Network Failure (Ping Failed)
	if err != nil {
		metrics.LatencyMs = 9999
		metrics.PacketLoss = 100.0
	}

	// 2. AI Scoring (Reinforcement Learning - Simple Q-Learning implementation)
	// Reward: Low Latency, Low Loss
	// Punishment: High Latency, High Loss

	// Calculate performance score (0-100, where 100 is perfect)
	// Base score starts at 100, subtract for latency and loss
	performance := 100.0
	performance -= float64(metrics.LatencyMs) / 10.0 // -10 points per 100ms
	performance -= metrics.PacketLoss * 2.0          // -2 points per 1% loss
	if performance < 0 {
		performance = 0
	}

	// Update Knowledge Base (Channel Score) with Learning Rate
	oldScore := float64(ap.channelScores[currentChannel])
	newScore := oldScore + ap.learningRate*(performance-oldScore) // Bellman-like update
	ap.channelScores[currentChannel] = int(newScore)

	// Save Updated Brain to Disk
	ap.Persist()

	logger.Debug("AI Analysis | Channel: %d | Performance: %.1f | Updated Score: %d",
		currentChannel, performance, ap.channelScores[currentChannel])

	// 3. Decide (Based on Thresholds AND AI Score)
	isLatencyBad := metrics.LatencyMs > ap.config.MaxLatencyMs
	isLossBad := metrics.PacketLoss > ap.config.MaxPacketLoss
	isScoreBad := ap.channelScores[currentChannel] < 40 // If AI thinks this channel sucks

	if isLatencyBad || isLossBad || isScoreBad {
		ap.failureStreak++

		reason := "Unknown"
		if isLossBad {
			reason = "High Packet Loss"
		}
		if isScoreBad {
			reason = "Low AI Score (Bad Reputation)"
		}

		logger.Warn("AI ALERT: Network Degraded. Cause: %s. Failure Streak: %d", reason, ap.failureStreak)

		// Trigger Action if persistent
		if ap.failureStreak >= ap.config.StabilityDuration {
			ap.takeCorrectiveAction(metrics)
			ap.failureStreak = 0 // Reset after taking action
		}
	} else {
		if ap.failureStreak > 0 {
			logger.Success("CONNECTION STABLE. AI Confidence: %d%%", ap.channelScores[currentChannel])
		}
		ap.failureStreak = 0
	}

	// 4. Traffic Awareness
	trafficType, _ := ap.manager.GetTrafficAnalysis()
	if trafficType != ap.currentTraffic {
		// New Context Detected!
		logger.Info("AI Context Switch: Traffic changed from '%s' to '%s'", ap.currentTraffic, trafficType)
		ap.currentTraffic = trafficType
		// Immediate check if we need to optimize for this updated context
		ap.optimizeForTraffic(trafficType, activeWAN)
	}
}

func (ap *AutoPilot) takeCorrectiveAction(metrics network.NetworkMetrics) {
	logger.Info("🤖 AI TAKING CONTROL: Executing Mitigation Strategy...")

	// Option 1: Channel Hop (If current channel score is low)
	currentChannel, _, _ := ap.manager.GetWifiInfo()
	if ap.channelScores[currentChannel] < 40 {
		// Find best channel
		bestChannel := currentChannel
		bestScore := -1
		for ch, score := range ap.channelScores {
			if score > bestScore && ch != currentChannel {
				bestScore = score
				bestChannel = ch
			}
		}
		if bestChannel != currentChannel {
			logger.Info("Decision: Switching to Channel %d (Score: %d) to escape interference.", bestChannel, bestScore)
			ap.manager.SetWifiChannel(bestChannel)
			return // Action taken
		}
	}

	// Option 2: WAN Failover (If Wi-Fi wasn't the problem, or couldn't be fixed)
	activeWAN, _ := ap.manager.GetActiveWAN()
	ap.triggerRecovery(metrics, activeWAN)
}

func (ap *AutoPilot) optimizeForTraffic(trafficType, activeWAN string) {
	targetWAN := ""

	switch trafficType {
	case "Gaming":
		// Gaming requires Low Latency (Fiber/5G)
		// If we are on a high-latency backup or busy line, switch!
		if activeWAN != "wan2_low_latency" {
			targetWAN = "wan2_low_latency" // This matches our simulated "Low Latency" WAN
		}
		// Also, Gaming might want to bypass VPN
		if status, _ := ap.manager.GetVPNStatus(); status == "Connected" {
			logger.Info("🎮 GAMING MODE: Disabling VPN to reduce latency...")
			ap.manager.DisableVPN()
		}

	case "Streaming":
		// Streaming requires High Bandwidth (Starlink/Cable)
		if activeWAN != "wan1_primary" {
			targetWAN = "wan1_primary"
		}
	}

	if targetWAN != "" && targetWAN != activeWAN {
		logger.Info("AI DETECTED %s TRAFFIC. Optimizing route...", trafficType)
		logger.Warn("Traffic Optimization: Switching to %s for %s", targetWAN, trafficType)
		ap.manager.SwitchWAN(targetWAN)
	}
}

func (ap *AutoPilot) triggerRecovery(metrics network.NetworkMetrics, activeWAN string) {

	triggerReason := "Latency"
	if metrics.PacketLoss > ap.config.MaxPacketLoss {
		triggerReason = "Packet Loss"
	}

	logger.Error("ACTION REQUIRED | Detected High %s (Lat: %dms, Loss: %.1f%%)",
		triggerReason, metrics.LatencyMs, metrics.PacketLoss)

	if activeWAN == "wan1_primary" {
		logger.Warn("FAILOVER STRATEGY: Switching to Backup WAN...")
		err := ap.manager.SwitchWAN("wan2_backup")
		if err != nil {
			logger.Error("Failover Failed: %v", err)
		} else {
			logger.Success("Switched to Secondary WAN (Backup)")
		}
	} else {
		logger.Warn("RECOVERY STRATEGY: Already on Backup. Restarting Interface...")
		err := ap.manager.RestartInterface(activeWAN)
		if err != nil {
			logger.Error("Interface Restart Failed: %v", err)
		} else {
			logger.Success("Interface Restarted")
		}
	}
}

func (ap *AutoPilot) optimizeWifi() {
	// 1. Check current health
	currentChannel, quality, err := ap.manager.GetWifiInfo()
	if err != nil {
		return
	}

	// 2. Decide: If quality is bad (<40), trigger scan
	if quality < 40 {
		logger.Warn("Wi-Fi Quality degraded on Channel %d (Score: %d/100)", currentChannel, quality)
		logger.Info("Initiating Intelligent Channel Optimization...")

		bestChannel := ap.findBestChannel()
		if bestChannel != -1 && bestChannel != currentChannel {
			err := ap.manager.SetWifiChannel(bestChannel)
			if err != nil {
				logger.Error("Wi-Fi Optimization Failed: %v", err)
			} else {
				logger.Success("Wi-Fi Optimized! Switched Intelligently to Channel %d", bestChannel)
			}
		} else {
			logger.Info("Current channel is still the best available option.")
		}
	}
}

func (ap *AutoPilot) findBestChannel() int {
	candidates, err := ap.manager.ScanWifiChannels()
	if err != nil {
		logger.Error("Scan failed: %v", err)
		return -1
	}

	bestScore := -1
	bestCh := -1

	for _, c := range candidates {
		// AI FUSION: Combine simulated scan data with learned experience
		// If we have history for this channel, weigh it in.
		learnedScore, exists := ap.channelScores[c.Channel]
		if !exists {
			learnedScore = 50 // Neutral bias for unknown channels
		}

		// Combined Score: 70% Scan (Current Reality), 30% AI Memory (History)
		// This prevents jumping to a channel that historically sucks.
		combinedScore := int((float64(c.Score) * 0.7) + (float64(learnedScore) * 0.3))

		logger.Debug("Channel Candidate %d | Scan: %d | AI Memory: %d | Final: %d", c.Channel, c.Score, learnedScore, combinedScore)

		if combinedScore > bestScore {
			bestScore = combinedScore
			bestCh = c.Channel
		}
	}
	return bestCh
}
