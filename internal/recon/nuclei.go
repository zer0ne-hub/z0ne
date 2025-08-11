package recon

import (
	"context"

	nuclei "github.com/projectdiscovery/nuclei/v3/lib"
	"github.com/projectdiscovery/nuclei/v3/pkg/installer"
	"github.com/projectdiscovery/nuclei/v3/pkg/output"
	syncutil "github.com/projectdiscovery/utils/sync"
)

func RunNuclei(target string) (interface{}, error) {

	var results []map[string]interface{}

	ctx := context.Background()
	tm := installer.TemplateManager{}
	if err := tm.FreshInstallIfNotExists(); err != nil {
		panic(err)
	}

	// create nuclei engine with options
	ne, err := nuclei.NewThreadSafeNucleiEngineCtx(ctx)
	if err != nil {
		panic(err)
	}
	// setup sizedWaitgroup to handle concurrency
	sg, err := syncutil.New(syncutil.WithSize(10))
	if err != nil {
		panic(err)
	}

	//callback function to save results
	onResult := func(result *output.ResultEvent) {
		results = append(results, map[string]interface{}{
			"ip":   result.IP,
			"port": result.Port,
			"host": result.Host,
			"info": result.Info,
		})
	}
	ne.GlobalResultCallback(onResult)

	// scan 1 = run dns templates on target
	sg.Add()
	go func() {
		defer sg.Done()
		err = ne.ExecuteNucleiWithOpts([]string{target},
			nuclei.WithTemplateFilters(nuclei.TemplateFilters{ProtocolTypes: "dns"}),
			nuclei.WithHeaders([]string{"X-Bug-Bounty: pdteam"}),
			nuclei.EnablePassiveMode(),
		)
		if err != nil {
			panic(err)
		}
	}()

	// scan 2 = run templates with oast tags on target
	sg.Add()
	go func() {
		defer sg.Done()
		err = ne.ExecuteNucleiWithOpts([]string{"http://" + target},
			nuclei.WithTemplateFilters(nuclei.TemplateFilters{Tags: []string{"oast"}}))
		if err != nil {
			panic(err)
		}
	}()

	// wait for all scans to finish
	sg.Wait()
	defer ne.Close()

	// Output:
	return results, nil
}
