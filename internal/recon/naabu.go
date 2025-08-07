package recon

import (
	"context"
	"fmt"
	"log"

	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/naabu/v2/pkg/result"
	"github.com/projectdiscovery/naabu/v2/pkg/runner"
)

func RunNaabu(target string, ports string, scanType string) {
	options := runner.Options{
		Host:     goflags.StringSlice{target},
		ScanType: scanType,
		OnResult: func(hr *result.HostResult) {
			fmt.Println(hr.Host, hr.Ports)
			//TODO: Add to report
		},
		Ports: ports,
	}

	naabuRunner, err := runner.NewRunner(&options)
	if err != nil {
		log.Fatal(err)
	}
	defer naabuRunner.Close()
	naabuRunner.RunEnumeration(context.Background())
}
