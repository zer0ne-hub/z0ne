package recon

import (
	"fmt"

	"github.com/projectdiscovery/dnsx/libs/dnsx"
)

// RunDnsX performs DNS lookups and returns structured results for saving
func RunDnsX(target string) (interface{}, error) {
	// Create DNS Resolver with default options
	dnsClient, err := dnsx.New(dnsx.DefaultOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create dnsx client: %w", err)
	}

	// Lookup IPs
	ips, err := dnsClient.Lookup(target)
	if err != nil {
		return nil, fmt.Errorf("lookup error: %w", err)
	}

	// Query for raw DNS response
	rawResp, err := dnsClient.QueryOne(target)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	// Get raw JSON string from dnsx
	jsonStr, err := rawResp.JSON()
	if err != nil {
		return nil, fmt.Errorf("json encode error: %w", err)
	}

	// Return all in a structured map
	results := map[string]interface{}{
		"ips":         ips,
		"rawResponse": jsonStr,
	}

	return results, nil
}
