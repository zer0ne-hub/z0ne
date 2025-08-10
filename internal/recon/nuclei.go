package recon

import (
	"context"
	"fmt"

	nuclei "github.com/projectdiscovery/nuclei/v3/lib"
)

func RunNuclei(target string) error {
	ne, err := nuclei.NewNucleiEngineCtx(context.Background(),
		nuclei.WithTemplateFilters(nuclei.TemplateFilters{Tags: []string{"oast"}}),
		nuclei.EnableStatsWithOpts(nuclei.StatsOptions{MetricServerPort: 6064}), // optionally enable metrics server for better observability
	)
	if err != nil {
		fmt.Println(err)
	}
	// load targets and optionally probe non http/https targets
	ne.LoadTargets([]string{target}, false)
	err = ne.ExecuteWithCallback(nil)
	if err != nil {
		fmt.Println(err)
	}
	defer ne.Close()

	return nil
}
