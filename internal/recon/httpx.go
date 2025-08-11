package recon

import (
	"log"

	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
	"github.com/projectdiscovery/httpx/runner"
)

func RunHttpX(target string) (interface{}, error) {
	gologger.DefaultLogger.SetMaxLevel(levels.LevelVerbose) // increase the verbosity (optional)

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
				"input":       r.Input,
				"host":        r.Host,
				"url":         r.URL,
				"status_code": r.StatusCode,
				"title":       r.Title,
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

	return results, nil
}
