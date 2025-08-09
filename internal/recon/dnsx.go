package recon

import (
	"fmt"

	"github.com/projectdiscovery/dnsx/libs/dnsx"
)

func RunDnsX(target string) {
	// Create DNS Resolver with default options
	dnsClient, err := dnsx.New(dnsx.DefaultOptions)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	// DNS A question and returns corresponding IPs
	result, err := dnsClient.Lookup(target)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	for idx, msg := range result {
		fmt.Printf("%d: %s\n", idx+1, msg)
	}

	// Query
	rawResp, err := dnsClient.QueryOne(target)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Printf("rawResp: %v\n", rawResp)

	jsonStr, err := rawResp.JSON()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Println(jsonStr)

	return
}
