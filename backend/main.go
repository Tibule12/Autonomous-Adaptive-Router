package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/TMtshwelo/aar/engine"
	"github.com/TMtshwelo/aar/network"
)

type StatusResponse struct {
	Status       string `json:"status"`
	Connectivity bool   `json:"connectivity"`
	VPNStatus    string `json:"vpn_status"`
	Latency      int64  `json:"latency_ms"`
	ActiveWAN    string `json:"active_wan"`
	WifiChannel  int    `json:"wifi_channel"`
	WifiQuality  int    `json:"wifi_quality"`
}

func main() {
	fmt.Println("Starting AAR Control Daemon...")

	// 1. Initialize Network Manager (Hardware Layer)
	netMgr := network.NewManager()

	// 2. Initialize and Start Auto-Pilot (Decision Layer)
	brain := engine.NewAutoPilot(netMgr)
	brain.Start()

	// Global state for simple demo (used only for API response)
	// The real monitoring is now done inside the engine.
	var currentLatency int64
	var wifiCh, wifiQual int
	
	// Separate loop just for the dashboard API data updates
	go func() {
		for {
			var err error
			currentLatency, err = netMgr.GetLatency("8.8.8.8")
			if err != nil {
				// Ignore metric errors for the dashboard ticker
			}
			
			wifiCh, wifiQual, _ = netMgr.GetWifiInfo()

			time.Sleep(2 * time.Second)
		}
	}()

	// Health Check / Status API
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		isConnected, _ := netMgr.CheckConnectivity()
		vpnStatus, _ := netMgr.GetVPNStatus()
		activeWAN, _ := netMgr.GetActiveWAN()

		resp := StatusResponse{
			Status:       "Running (Auto-Pilot Active)",
			Connectivity: isConnected,
			VPNStatus:    vpnStatus,
			Latency:      currentLatency,
			ActiveWAN:    activeWAN,
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

	fmt.Println("Server listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed: %s\n", err)
	}
}
