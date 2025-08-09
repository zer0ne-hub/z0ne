package recon

import (
	"bytes"
	"context"
	"io"
	"log"

	"github.com/projectdiscovery/subfinder/v2/pkg/runner"
)

func RunSubfinder(target string) {
	subfinderOpts := &runner.Options{
		Threads:            10,
		Timeout:            30,
		MaxEnumerationTime: 10,
		// ResultCallback: func(s *resolve.HostEntry) {
		// callback function executed after each unique subdomain is found
		// },
		// ProviderConfig: "your_provider_config.yaml",
		// and other config related options
	}

	// disable timestamps in logs / configure logger
	log.SetFlags(0)

	subfinder, err := runner.NewRunner(subfinderOpts)
	if err != nil {
		log.Fatalf("failed to create subfinder runner: %v", err)
	}
	output := &bytes.Buffer{}
	var sourceMap map[string]map[string]struct{}

	// To run subdomain enumeration on a single domain
	if sourceMap, err = subfinder.EnumerateSingleDomainWithCtx(context.Background(), target, []io.Writer{output}); err != nil {
		log.Fatalf("failed to enumerate single domain: %v", err)
	}

	// To run subdomain enumeration on a list of domains from file/reader
	// file, err := os.Open("domains.txt")
	// if err != nil {
	// 	log.Fatalf("failed to open domains file: %v", err)
	// }
	// defer file.Close()
	// if err = subfinder.EnumerateMultipleDomainsWithCtx(context.Background(), file, []io.Writer{output}); err != nil {
	// 	log.Fatalf("failed to enumerate subdomains from file: %v", err)
	// }

	// print the output
	log.Println(output.String())

	// Or use sourceMap to access the results in your application
	for subdomain, sources := range sourceMap {
		sourcesList := make([]string, 0, len(sources))
		for source := range sources {
			sourcesList = append(sourcesList, source)
		}
		log.Printf("%s %s (%d)\n", subdomain, sourcesList, len(sources))
	}
}
