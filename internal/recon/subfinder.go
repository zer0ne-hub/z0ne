// Package recon: Handles all Reconnaissance modules independently
package recon

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"

	"github.com/projectdiscovery/subfinder/v2/pkg/runner"
)

// RunSubfinder enumerates subdomains for a target and returns them in a JSON-friendly format
func RunSubfinder(target string) (interface{}, error) {
	subfinderOpts := &runner.Options{
		Threads:            10,
		Timeout:            30,
		MaxEnumerationTime: 10,
	}

	log.SetFlags(0)

	subfinderRunner, err := runner.NewRunner(subfinderOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create subfinder runner: %w", err)
	}

	output := &bytes.Buffer{}
	sourceMap, err := subfinderRunner.EnumerateSingleDomainWithCtx(
		context.Background(),
		target,
		[]io.Writer{output},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to enumerate single domain: %w", err)
	}

	// Build results slice
	var results []map[string]interface{}
	for subdomain, sources := range sourceMap {
		srcList := make([]string, 0, len(sources))
		for s := range sources {
			srcList = append(srcList, s)
		}
		results = append(results, map[string]interface{}{
			"subdomain": subdomain,
			"sources":   srcList,
			"count":     len(srcList),
		})
	}

	fmt.Println("Subfinder found:", results[0]["subdomain"])

	return results, nil
}
