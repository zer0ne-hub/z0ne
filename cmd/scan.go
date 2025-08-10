package cmd

import (
	"z0ne/internal/core"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "recon [target]",
	Short: "Basic web reconnaissance",
	Long: `This mode will enumerate subdomains, resolve IP addresses, and scan ports... on a given target.
	fast, minimal but large scale web reconnaissance.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		color.Cyan("🧿 Scanning: %s", target)
		core.RunRecon(target)
		color.Green("Scan complete!")
	},
}
