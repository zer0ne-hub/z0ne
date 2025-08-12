// Package report: Handles report generation
package report

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ReconReport: Struct to hold report data
type ReconReport struct {
	Target     string
	Results    map[string]interface{}
	RedactedAt time.Time
}

// GenerateReport: generates a markdown report for a given target from the json results
func GenerateReport(targetName string) error {
	var outputDir = "z0ne-out"
	var resultFile = filepath.Join(outputDir, "results.json")
	var reportFile = filepath.Join(outputDir, "report.md")

	jsonFile, err := os.Open(resultFile)
	if err != nil {
		return fmt.Errorf("failed to open results.json: %w", err)
	}
	defer jsonFile.Close()

	jsonBytes, err := os.ReadFile(resultFile)
	if err != nil {
		return fmt.Errorf("failed to read results.json: %w", err)
	}

	var results map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &results); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// build the report
	report := ReconReport{
		Target:     targetName,
		Results:    results,
		RedactedAt: time.Now(),
	}

	// make markdown content
	md := fmt.Sprintf("# %s\n\n", report.Target)
	md += fmt.Sprintf("**Report geberated at: ** %s\n\n", report.RedactedAt.Format(time.RFC1123))
	md += "## Assets Discovered\n\n"
	md += "### Domains and IPs\n\n"
	md += ""

	// write report.md
	err = os.WriteFile(reportFile, []byte(md), 0644)
	if err != nil {
		return fmt.Errorf("failed to write report.md: %w", err)
	}

	return nil
}
