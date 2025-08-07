package report

import (
	"errors"
	"time"
)

type ReportFormat string

const (
	JSON     ReportFormat = "json"
	Markdown ReportFormat = "md"
	PDF      ReportFormat = "pdf"
)

type ReconReport struct {
	Target     string
	Type       string
	ModulesRun []string
	Results    map[string]interface{}
	StartedAt  time.Time
	FinishedAt time.Time
}

func Generate(report ReconReport, format ReportFormat, outPath string) error {
	//TODO: Generate report
	return errors.New("not implemented")
}
