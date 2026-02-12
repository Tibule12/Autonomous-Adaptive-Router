package engine

import (
	"encoding/json"
	"os"

	"github.com/TMtshwelo/aar/pkg/logger"
)

const MemoryFile = "ai_memory.json"

// LoadMemory attempts to read the AI's past experiences from disk
func (ap *AutoPilot) LoadMemory() {
	file, err := os.ReadFile(MemoryFile)
	if err != nil {
		logger.Info("Ez: No previous AI memory found. Starting fresh.")
		// Default initialized in NewAutoPilot
		return
	}

	var memory map[int]int
	if err := json.Unmarshal(file, &memory); err != nil {
		logger.Error("Failed to parse AI memory: %v", err)
		return
	}

	ap.channelScores = memory
	logger.Success("🧠 AI Memory Loaded! Knowledge restored for %d channels.", len(memory))
}

// SaveMemory writes the current channel scores to disk
func (ap *AutoPilot) SaveMemory() {
	// Simple throttling: Only save if we haven't saved in the last 10 seconds?
	// For now, simpler is better. We'll save on every update.
	// In a real production system, we'd debounce this.

	// Ensure we don't save empty states
	if len(ap.channelScores) == 0 {
		return
	}

	data, err := json.MarshalIndent(ap.channelScores, "", "  ")
	if err != nil {
		logger.Error("Failed to serialize AI memory: %v", err)
		return
	}

	// Write to file (Truncate/Create)
	if err := os.WriteFile(MemoryFile, data, 0644); err != nil {
		logger.Error("Failed to save AI memory: %v", err)
	}
}

// Persist calls SaveMemory with a timestamp log (optional, effectively a wrapper)
func (ap *AutoPilot) Persist() {
	go func() {
		ap.SaveMemory()
		// logger.Debug("AI Memory Persisted to disk.")
	}()
}
