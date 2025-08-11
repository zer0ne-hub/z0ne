// performs utilities tasks like parsing target input
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
	"sort"
	"strings"
	"sync"
)

type TargetType int

const (
	UNKNOWN TargetType = iota
	IP
	DOMAIN
	URL
	FILE
)

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

// SaveResultToJSON safely writes or updates a tool's result in /z0ne-out/results.json
// Thread-safe: only one goroutine can modify the file at a time.
func SaveResultToJSON(toolName string, resultData interface{}) error {
	resultsFileLock.Lock()
	defer resultsFileLock.Unlock()

	outputDir := "z0ne-out"
	outputFile := filepath.Join(outputDir, "results.json")

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, fs.ModePerm); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Load existing results if file exists
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

	// Update or add the tool's result
	results[toolName] = resultData

	// Encode and write updated results
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
	// Add more flags API keys as needed
}

// formatMap formats a map[string]interface{} recursively into a Markdown string
// used in the report generation
func FormatMap(data interface{}, indent int) string {
	prefix := strings.Repeat("  ", indent)
	switch v := data.(type) {
	case map[string]interface{}:
		md := ""
		for key, val := range v {
			if s, ok := val.(string); ok && len(s) > 60 && (containsNewline(s) || looksLikeJSON(s)) {
				md += fmt.Sprintf("%s- **%s**:\n\n```\n%s\n```\n\n", prefix, key, s)
			} else {
				md += fmt.Sprintf("%s- **%s**: %s\n", prefix, key, FormatMap(val, indent+1))
			}
		}
		return md
	case []interface{}:
		if len(v) > 0 {
			if _, ok := v[0].(map[string]interface{}); ok {
				return formatTable(v)
			}
		}
		md := ""
		for _, val := range v {
			md += fmt.Sprintf("%s- %s\n", prefix, FormatMap(val, indent+1))
		}
		return md
	case string:
		return v
	case float64:
		return fmt.Sprintf("%v", v)
	case bool:
		return fmt.Sprintf("%v", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func containsNewline(s string) bool {
	return strings.Contains(s, "\n")
}

func looksLikeJSON(s string) bool {
	return strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}")
}

func formatTable(items []interface{}) string {
	colSet := make(map[string]struct{})
	for _, item := range items {
		if m, ok := item.(map[string]interface{}); ok {
			for k := range m {
				colSet[k] = struct{}{}
			}
		}
	}
	cols := make([]string, 0, len(colSet))
	for k := range colSet {
		cols = append(cols, k)
	}
	sort.Strings(cols)

	md := "| " + strings.Join(cols, " | ") + " |\n"
	md += "| " + strings.Repeat("--- | ", len(cols)) + "\n"

	for _, item := range items {
		if m, ok := item.(map[string]interface{}); ok {
			row := make([]string, len(cols))
			for i, col := range cols {
				if val, exists := m[col]; exists {
					row[i] = fmt.Sprintf("%v", val)
				} else {
					row[i] = ""
				}
			}
			md += "| " + strings.Join(row, " | ") + " |\n"
		}
	}
	return md + "\n"
}
