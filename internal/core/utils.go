// Package core: Handles core engine logic
package core

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

// TargetType: Enumeration for different types of targets
type TargetType int

const (
	UNKNOWN TargetType = iota
	IP
	DOMAIN
	URL
	FILE
)

// detectTargetType: Detects the type of the target passed as an argument
func detectTargetType(input string) TargetType {
	input = strings.TrimSpace(input)
	if _, err := os.Stat(input); err == nil || filepath.IsAbs(input) {
		return FILE
	}
	if net.ParseIP(input) != nil {
		return IP
	}
	domainRegex := regexp.MustCompile(`^(?:[a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`)
	if domainRegex.MatchString(input) {
		return DOMAIN
	}
	if u, err := url.Parse(input); err == nil && u.Scheme != "" && u.Host != "" {
		return URL
	}
	return UNKNOWN
}

// Global mutex to lock file writes across goroutines
var resultsFileLock sync.Mutex

// SaveResultToJSON: safely writes or updates a tool's result in /z0ne-out/results.json
// Thread-safe because only one goroutine can modify the file at a time.
// Update existing results if file exists or create a new file
func SaveResultToJSON(toolName string, resultData interface{}) error {
	resultsFileLock.Lock()
	defer resultsFileLock.Unlock()

	outputDir := "z0ne-out"
	outputFile := filepath.Join(outputDir, "results.json")

	if err := os.MkdirAll(outputDir, fs.ModePerm); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	results := make(map[string]interface{})
	if fileData, err := os.ReadFile(outputFile); err == nil {
		if len(fileData) > 0 {
			if err := json.Unmarshal(fileData, &results); err != nil {
				return fmt.Errorf("failed to parse existing JSON: %w", err)
			}
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to read existing results file: %w", err)
	}
	results[toolName] = resultData
	fileData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode results to JSON: %w", err)
	}

	if err := os.WriteFile(outputFile, fileData, 0644); err != nil {
		return fmt.Errorf("failed to write results file: %w", err)
	}

	return nil
}

// ProbeKeys struct to hold API keys
type ProbeKeys struct {
	ShodanKey string
}
