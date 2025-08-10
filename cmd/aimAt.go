package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zer0ne-hub/z0ne/internal/core"
)

var aimAtCmd = &cobra.Command{
	Use:   "aimAt [target]",
	Short: "Set a target to attack",
	Long: `This mode will setup the workspace for a given target. ALWAYS deleting previous results,
	and creating a new results file.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		color.Cyan("🎯 Aimming at %s", target)
		core.SetupTarget(target)
		color.Green("Target locked!")
	},
}
