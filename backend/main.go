package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/TMtshwelo/aar/engine"
	"github.com/TMtshwelo/aar/network"
	"github.com/TMtshwelo/aar/pkg/logger"
	"github.com/TMtshwelo/aar/pkg/storage"
)

type StatusResponse struct {
	Status       string  `json:"status"`
	Connectivity bool    `json:"connectivity"`
	VPNStatus    string  `json:"vpn_status"`
	Latency      int64   `json:"latency_ms"`
	PacketLoss   float64 `json:"packet_loss"` // New Metric
	ActiveWAN    string  `json:"active_wan"`
	TrafficType  string  `json:"traffic_type"`
	WifiChannel  int     `json:"wifi_channel"`
	WifiQuality  int     `json:"wifi_quality"`
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()
	logger.Init()
	logger.Info("Starting AAR Control Daemon...")

	// 1. Initialize Network Manager (Hardware Layer)
	netMgr := network.NewManager()

	// 2. Initialize and Start Auto-Pilot (Decision Layer)
	brain := engine.NewAutoPilot(netMgr)
	brain.Start()

	// Global state for simple demo
	var currentStats network.NetworkMetrics
	var wifiCh, wifiQual int
	var requestCounter int64 // Counter for stress testing

	// Separate loop just for the dashboard API data updates
	go func() {
		for {
			// Update Simulation with Real Traffic Load
			// We measure requests in the last 2 seconds (interval)
			load := atomic.SwapInt64(&requestCounter, 0)
			// Normalize to req/sec (approximate since sleep is 2s)
			netMgr.SetSimulatedLoad(int(load / 2))

			var err error
			currentStats, err = netMgr.GetNetworkMetrics("8.8.8.8")
			if err != nil {
				// Ignore metric errors for the dashboard ticker
			}

			// Auto-Save History
			storage.SaveMetric(currentStats.LatencyMs, currentStats.PacketLoss)

			wifiCh, wifiQual, _ = netMgr.GetWifiInfo()

			time.Sleep(2 * time.Second)
		}
	}()

	// Stress Test Endpoint (Hit this from Laptop B)
	http.HandleFunc("/api/stress", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		atomic.AddInt64(&requestCounter, 1)
		w.WriteHeader(http.StatusOK)
	})

	// Health Check / Status API
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		isConnected, _ := netMgr.CheckConnectivity()
		vpnStatus, _ := netMgr.GetVPNStatus()
		activeWAN, _ := netMgr.GetActiveWAN()
		trafficType, _ := netMgr.GetTrafficAnalysis()

		resp := StatusResponse{
			Status:       "Running (Auto-Pilot Active)",
			Connectivity: isConnected,
			VPNStatus:    vpnStatus,
			Latency:      currentStats.LatencyMs,
			PacketLoss:   currentStats.PacketLoss,
			ActiveWAN:    activeWAN,
			TrafficType:  trafficType,
			WifiChannel:  wifiCh,
			WifiQuality:  wifiQual,
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		json.NewEncoder(w).Encode(resp)
	})

	// VPN Control API
	http.HandleFunc("/vpn/toggle", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

		if r.Method == "OPTIONS" {
			return
		}

		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		currentStatus, _ := netMgr.GetVPNStatus()
		var err error
		var newState string

		if currentStatus == "Connected" {
			err = netMgr.DisableVPN()
			newState = "Disconnected"
		} else {
			err = netMgr.EnableVPN()
			newState = "Connected"
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"status": newState})
	})

	// Chaos Control API
	http.HandleFunc("/chaos/lag", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

		if r.Method == "OPTIONS" {
			return
		}

		type ChaosRequest struct {
			Enable bool   `json:"enable"`
			Type   string `json:"type"` // "lag" or "loss"
		}
		var req ChaosRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if req.Type == "loss" {
			netMgr.SetSimulatedPacketLoss(req.Enable)
			fmt.Printf("[API] Chaos Packet Loss set to: %v\n", req.Enable)
		} else if req.Type == "interference" {
			netMgr.SetSimulatedInterference(req.Enable)
			fmt.Printf("[API] Chaos Wi-Fi Interference set to: %v\n", req.Enable)
		} else {
			netMgr.SetSimulatedLag(req.Enable)
			fmt.Printf("[API] Chaos Lag set to: %v\n", req.Enable)
		}

		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// Traffic Simulation API
	http.HandleFunc("/simulation/traffic", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

		if r.Method == "OPTIONS" {
			return
		}

		type TrafficRequest struct {
			Type string `json:"type"` // "Gaming", "Streaming", "Default"
		}
		var req TrafficRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		netMgr.SetSimulatedTraffic(req.Type)
		logger.Info("Simulating User Activity: %s", req.Type)

		json.NewEncoder(w).Encode(map[string]string{"status": "updated", "type": req.Type})
	})

	// System Logs API
	http.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(logger.GetLogs())
	})

	// Metrics History API
	http.HandleFunc("/history", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(storage.GetHistory())
	})

	// Connected Devices API (GET / POST)
	http.HandleFunc("/api/devices", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE")
		w.Header().Set("Content-Type", "application/json")

		if r.Method == "OPTIONS" {
			return
		}

		if r.Method == "GET" {
			devices, err := netMgr.GetConnectedDevices()
			if err != nil {
				json.NewEncoder(w).Encode([]network.Device{})
				return
			}
			json.NewEncoder(w).Encode(devices)
			return
		}

		// Handle BLOCK/UNBLOCK via POST/DELETE
		type DeviceAction struct {
			MAC    string `json:"mac"`
			Action string `json:"action"` // "block" or "unblock"
		}
		var req DeviceAction
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if req.Action == "block" {
			netMgr.BlockDevice(req.MAC)
			logger.Info("Blocked device: %s", req.MAC)
		} else {
			netMgr.UnblockDevice(req.MAC)
			logger.Info("Unblocked device: %s", req.MAC)
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	})

	// Serve Frontend (Robust check for Portable vs Dev mode)
	var frontendDir string
	if _, err := os.Stat("./frontend_dist"); err == nil {
		frontendDir = "./frontend_dist"
	} else if _, err := os.Stat("../frontend/dist"); err == nil {
		frontendDir = "../frontend/dist"
	} else {
		fmt.Println("Warning: Frontend dist folder not found. Dashboard will be blank.")
	}

	if frontendDir != "" {
		fmt.Printf("Serving frontend from: %s\n", frontendDir)
		fs := http.FileServer(http.Dir(frontendDir))
		http.Handle("/", fs)
	}

	fmt.Println("Server listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Error("Server failed: %s", err)
	}
}
