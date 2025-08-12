// Package cmd: Handles command line interface using cobra
package cmd

import (
	"github.com/fatih/color"

	"github.com/spf13/cobra"
	"github.com/zer0ne-hub/z0ne/internal/report"
)

// reportCmd: represents the report command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate report",
	Long: `This mode will generate a report for a given target name,
	from the /z0ne-out/results.json file. It checks for the file so other modules
	may be run before generating the report. The format is markdown.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetName := args[0]
		color.Cyan("📄Generating report for: %s", targetName)
		err := report.GenerateReport(targetName)
		if err != nil {
			color.Red("Error generating report: %v", err)
		}
		color.Green("Report generated successfully! File saved to /z0ne-out/report.md")
		color.Green("You can convert it to PDF or HTML with Any tool like Pandoc.")
		color.Green("You can now run 'z0ne scan' or 'z0ne probe' to scan again. Happy Hunting!")
	},
}
