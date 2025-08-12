// Package cmd: Handles command line interface using cobra
package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zer0ne-hub/z0ne/internal/core"
)

// scanCmd: represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan [target]",
	Short: "Basic web reconnaissance",
	Long: `This mode will enumerate subdomains, resolve IP addresses, and scan ports... on a given target.
	fast, minimal but large scale web reconnaissance.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		color.Cyan("👀 Scanning: %s", target)
		err := core.RunRecon(target)
		if err != nil {
			color.Red("Scan Error: %v", err)
		}
		color.Green("\n\nScan complete! Results saved to /z0ne-out/results.json")
		color.Green("You can now run 'z0ne report [targetName]' to generate a MD report.")
	},
}
