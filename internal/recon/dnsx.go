package recon

import (
	"fmt"

	"github.com/projectdiscovery/dnsx/libs/dnsx"
)

func RunDnsX(target string) error {
	// Create DNS Resolver with default options
	dnsClient, err := dnsx.New(dnsx.DefaultOptions)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}

	// DNS A question and returns corresponding IPs
	result, err := dnsClient.Lookup(target)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	for idx, msg := range result {
		fmt.Printf("%d: %s\n", idx+1, msg)
	}

	// Query
	rawResp, err := dnsClient.QueryOne(target)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	fmt.Printf("rawResp: %v\n", rawResp)

	jsonStr, err := rawResp.JSON()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	fmt.Println(jsonStr)

	return nil
}
