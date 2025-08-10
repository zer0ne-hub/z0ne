package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type MetaInfo struct {
	StartedAt   time.Time `json:"started_at"`
	LastUpdated time.Time `json:"last_updated,omitempty"`
}

type Z0neResults struct {
	Target string   `json:"target"`
	Meta   MetaInfo `json:"meta"`
	//More tool output sections here
}

// Config structure (global config file)
type Z0neConfig struct {
	APIKeys map[string]string `json:"api_keys"`
}

func SetupTarget(target string) error {
	workspace := "out"
	configDir := "config"
	targetDir := filepath.Join(workspace, "target") // fixed folder
	resultsDir := filepath.Join(targetDir, "results")
	reportsDir := filepath.Join(targetDir, "reports")
	resultsFile := filepath.Join(targetDir, "target.json")

	// Create necessary directories
	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		return fmt.Errorf("failed to create results directory: %w", err)
	}
	if err := os.MkdirAll(reportsDir, 0755); err != nil {
		return fmt.Errorf("failed to create reports directory: %w", err)
	}

	// Always delete the previous results JSON if it exists
	if _, err := os.Stat(resultsFile); err == nil {
		if err := os.Remove(resultsFile); err != nil {
			return fmt.Errorf("failed to remove existing results file: %w", err)
		}
	}

	// Create new results JSON
	initial := Z0neResults{
		Target: target,
		Meta: MetaInfo{
			StartedAt: time.Now(),
		},
	}

	f, err := os.Create(resultsFile)
	if err != nil {
		return fmt.Errorf("failed to create results file: %w", err)
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(initial); err != nil {
		return fmt.Errorf("failed to write results JSON: %w", err)
	}

	// Ensure global config exists (only create if missing)
	globalConfigPath := filepath.Join(configDir, "config.json")
	if _, err := os.Stat(globalConfigPath); os.IsNotExist(err) {
		cfg := Z0neConfig{
			APIKeys: map[string]string{
				"uncover": "",
				"shodan":  "",
				"openai":  "",
				// Add more API keys here
			},
		}
		cf, err := os.Create(globalConfigPath)
		if err != nil {
			return fmt.Errorf("failed to create config file: %w", err)
		}
		defer cf.Close()

		enc := json.NewEncoder(cf)
		enc.SetIndent("", "  ")
		if err := enc.Encode(cfg); err != nil {
			return fmt.Errorf("failed to write config file: %w", err)
		}
		fmt.Println("[+] Global config created at", globalConfigPath)
	}

	fmt.Println("[+] Workspace and results initialized for", target)
	return nil
}
