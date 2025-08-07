package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate report",
	Long: `This mode will generate a report for a given target name,
	from other recon or attack surface mapping modules results.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetName := args[0]
		fmt.Println("📄Generating report for:", targetName)
		//TODO: Report generation

	},
}
