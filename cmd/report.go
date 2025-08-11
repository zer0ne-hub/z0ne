package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zer0ne-hub/z0ne/internal/report"
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate report",
	Long: `This mode will generate a report for a given target name,
	from the /z0ne-out/results.json file. It checks for the file so other modules
	may be run before generating the report. The format is markdown.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetName := args[0]
		fmt.Println("📄Generating report for:", targetName)
		err := report.GenerateReport(targetName)
		if err != nil {
			fmt.Println("Error generating report:", err)
		}
		fmt.Println("Report generated successfully!")
	},
}
