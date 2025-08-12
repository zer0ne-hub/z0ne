// Package recon: Handles all Reconnaissance modules independently
package recon

import (
	"fmt"

	"github.com/projectdiscovery/cdncheck"
	"github.com/projectdiscovery/dnsx/libs/dnsx"
)

// RunDnsX performs DNS lookups and returns structured results for saving
// Creates a DNS Resolver (dnsClient) with default options and performs various lookups
func RunDnsX(target string) (interface{}, error) {

	dnsClient, err := dnsx.New(dnsx.DefaultOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create dnsx client: %w", err)
	}

	cdnClient, err := cdncheck.NewWithOpts(3, cdncheck.DefaultResolvers)
	if err != nil {
		return nil, fmt.Errorf("failed to init cdncheck client: %w", err)
	}

	ips, err := dnsClient.Lookup(target)
	if err != nil {
		return nil, fmt.Errorf("lookup error: %w", err)
	}
	if len(ips) == 0 {
		return nil, fmt.Errorf("no results for %s", target)
	}

	trace, err := dnsClient.Trace(target)
	if err != nil {
		return nil, fmt.Errorf("trace error: %w", err)
	}

	isCDN, cdn, cdntype, err := cdnClient.CheckDomainWithFallback(target)
	if err != nil {
		return nil, fmt.Errorf("cdn check error: %w", err)
	}

	results := map[string]interface{}{
		"ips":      ips,
		"DNStrace": trace,
		"isCDN":    isCDN,
		"CDNType":  cdntype,
		"CDN":      cdn,
	}

	fmt.Println("DNSx found ips: ")
	for _, ip := range ips {
		fmt.Println(ip)
	}

	return results, nil
}
