package recon

import (
	"math"
	"sync"

	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/katana/pkg/engine/standard"
	"github.com/projectdiscovery/katana/pkg/output"
	"github.com/projectdiscovery/katana/pkg/types"
)

func RunKatana(target string) (interface{}, error) {
	var (
		resultsMu sync.Mutex
		results   []map[string]interface{}
	)

	options := &types.Options{
		MaxDepth:     3,
		FieldScope:   "rdn",
		BodyReadSize: math.MaxInt,
		Timeout:      10,
		Concurrency:  10,
		Parallelism:  10,
		Delay:        0,
		RateLimit:    150,
		Strategy:     "depth-first",
		OnResult: func(result output.Result) {
			resultsMu.Lock()
			results = append(results, map[string]interface{}{
				"url":    result.Request.URL,
				"method": result.Request.Method,
				"body":   result.Request.Body,
			})
			resultsMu.Unlock()
			gologger.Info().Msg(result.Request.URL)
		},
	}

	crawlerOptions, err := types.NewCrawlerOptions(options)
	if err != nil {
		return nil, err
	}
	defer crawlerOptions.Close()

	crawler, err := standard.New(crawlerOptions)
	if err != nil {
		return nil, err
	}
	defer crawler.Close()

	err = crawler.Crawl("https://" + target)
	if err != nil {
		gologger.Warning().Msgf("Could not crawl %s: %s", target, err.Error())
		// return partial results anyway
	}

	return results, nil
}
