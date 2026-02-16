package storage

import (
	"encoding/json"
	"os"
	"sync"
)

type RouterConfig struct {
	BlockedMACs []string `json:"blocked_macs"`
	WifiSSID    string   `json:"wifi_ssid"`
	WifiPass    string   `json:"wifi_pass"`
	GamerMode   bool     `json:"gamer_mode"`
	Theme       string   `json:"theme"` // For UI preferences
}

var (
	configFile   = "router_config.json"
	configMu     sync.Mutex
	cachedConfig *RouterConfig
)

// LoadConfig reads the configuration from disk
func LoadConfig() (*RouterConfig, error) {
	configMu.Lock()
	defer configMu.Unlock()

	if cachedConfig != nil {
		return cachedConfig, nil
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			// specific default
			return &RouterConfig{
				BlockedMACs: []string{},
				Theme:       "cyberpunk",
			}, nil
		}
		return nil, err
	}

	var config RouterConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	cachedConfig = &config
	return &config, nil
}

// SaveConfig writes the configuration to disk
func SaveConfig(config *RouterConfig) error {
	configMu.Lock()
	defer configMu.Unlock()

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return err
	}

	cachedConfig = config
	return nil
}

// Helper: AddBlockedMAC
func AddBlockedMAC(mac string) error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}

	// Check for duplicate
	for _, m := range cfg.BlockedMACs {
		if m == mac {
			return nil // Already blocked
		}
	}

	cfg.BlockedMACs = append(cfg.BlockedMACs, mac)
	return SaveConfig(cfg)
}

// Helper: RemoveBlockedMAC
func RemoveBlockedMAC(mac string) error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}

	newMACs := []string{}
	for _, m := range cfg.BlockedMACs {
		if m != mac {
			newMACs = append(newMACs, m)
		}
	}

	cfg.BlockedMACs = newMACs
	return SaveConfig(cfg)
}

// Helper: IsBlocked
func IsBlocked(mac string) bool {
	cfg, err := LoadConfig()
	if err != nil {
		return false
	}
	for _, m := range cfg.BlockedMACs {
		if m == mac {
			return true
		}
	}
	return false
}
