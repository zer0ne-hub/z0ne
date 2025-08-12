// Package recon: Handles all Reconnaissance modules independently
package recon

import (
	"context"
	"fmt"

	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/naabu/v2/pkg/result"
	"github.com/projectdiscovery/naabu/v2/pkg/runner"
)

// RunNaabu scans the target and returns results in a format ready for JSON saving
func RunNaabu(target, ports, scanType string) (interface{}, error) {
	var results []map[string]interface{}

	options := runner.Options{
		Host:             goflags.StringSlice{target},
		ScanType:         scanType,
		Ports:            ports,
		Verify:           true,
		Nmap:             true,
		Ping:             true,
		ServiceDiscovery: true,
		ServiceVersion:   true,
		OutputCDN:        true,
		OnResult: func(hr *result.HostResult) {
			entry := map[string]interface{}{
				"host":            hr.Host,
				"ip":              hr.IP,
				"ports":           hr.Ports,
				"confidenceLevel": hr.Confidence,
			}
			results = append(results, entry)
			fmt.Println("Naabu found:")
			fmt.Println(hr.Host, hr.Ports)
		},
	}

	naabuRunner, err := runner.NewRunner(&options)
	if err != nil {
		return nil, fmt.Errorf("could not create naabu runner: %w", err)
	}
	defer naabuRunner.Close()

	if err := naabuRunner.RunEnumeration(context.Background()); err != nil {
		return nil, fmt.Errorf("naabu enumeration failed: %w", err)
	}

	return results, nil
}
