package report

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/zer0ne-hub/z0ne/internal/core"
)

type ReconReport struct {
	Target     string
	Type       string
	ModulesRun []string
	Results    map[string]interface{}
	StartedAt  time.Time
	FinishedAt time.Time
}

// GenerateReport generates a report for a given target
func GenerateReport(targetName string) error {
	//Open and read results.json
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
		Type:       "Probe",
		ModulesRun: []string{},
		Results:    results,
		StartedAt:  time.Now().Add(-time.Hour),
		FinishedAt: time.Now(),
	}

	// make markdown content
	md := fmt.Sprintf("# Reconnaissance Report for %s\n\n", report.Target)
	md += fmt.Sprintf("**Type:** %s\n\n", report.Type)
	md += fmt.Sprintf("**Started:** %s\n\n", report.StartedAt.Format(time.RFC1123))
	md += fmt.Sprintf("**Finished:** %s\n\n", report.FinishedAt.Format(time.RFC1123))

	md += "## Modules Run\n"
	if len(report.ModulesRun) == 0 {
		md += "_No module info available_\n"
	} else {
		for _, m := range report.ModulesRun {
			md += fmt.Sprintf("- %s\n", m)
		}
	}
	md += "\n"

	md += "## Results\n"

	for module, data := range report.Results {
		md += fmt.Sprintf("### %s\n\n", module)
		md += core.FormatMap(data, 0)
		md += "\n"
	}

	// write report.md
	err = os.WriteFile(reportFile, []byte(md), 0644)
	if err != nil {
		return fmt.Errorf("failed to write report.md: %w", err)
	}

	return nil
}
