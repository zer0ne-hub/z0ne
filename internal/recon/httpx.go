// Package recon: Handles all Reconnaissance modules independently
package recon

import (
	"fmt"
	"log"

	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
	"github.com/projectdiscovery/httpx/runner"
)

// RunHttpX scans the target and returns results in a format ready for JSON saving
func RunHttpX(target string) (interface{}, error) {
	gologger.DefaultLogger.SetMaxLevel(levels.LevelVerbose)

	var results []map[string]interface{}

	options := runner.Options{
		Methods:         "GET",
		InputTargetHost: goflags.StringSlice{target},
		OnResult: func(r runner.Result) {
			if r.Err != nil {
				results = append(results, map[string]interface{}{
					"input": target,
					"error": r.Err.Error(),
				})
				return
			}
			results = append(results, map[string]interface{}{
				"input":               r.Input,
				"host":                r.Host,
				"url":                 r.URL,
				"status_code":         r.StatusCode,
				"title":               r.Title,
				"A":                   r.A,
				"AAAA":                r.AAAA,
				"CNAMES":              r.CNAMEs,
				"ASN":                 r.ASN,
				"Domains":             r.Domains,
				"isCDN":               r.CDN,
				"CDN":                 r.CDN,
				"CDNType":             r.CDNType,
				"location":            r.Location,
				"screeshotbytes":      r.ScreenshotBytes,
				"server":              r.WebServer,
				"vhost":               r.VHost,
				"technologies":        r.Technologies,
				"technologiesDetails": r.TechnologyDetails,
			})
		},
	}

	if err := options.ValidateOptions(); err != nil {
		log.Fatal(err)
	}

	httpxRunner, err := runner.New(&options)
	if err != nil {
		log.Fatal(err)
	}
	defer httpxRunner.Close()
	httpxRunner.RunEnumeration()

	fmt.Println("HTTPX Results on ", results[0]["input"])
	fmt.Println("Host is: ", results[0]["host"])
	if results[0]["isCDN"] == true {
		fmt.Println("CDN found")
		fmt.Println("CDN is: ", results[0]["CDN"])
		fmt.Println("CDN Type: ", results[0]["CDNType"])
	}
	fmt.Println("Technologies identified are: ", results[0]["technologies"])

	return results, nil
}
