package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zer0ne-hub/z0ne/internal/core"
)

var probeCmd = &cobra.Command{
	Use:   "probe [target]",
	Short: "more precise web reconnaissance",
	Long: `This mode will enumerate subdomains, resolve IP addresses, and scan ports... on a given target but
with more precision, using authentication, stealth, and some basic bypass methods. Might be slower.
	Needs a config file to be set up using aim command.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		color.Cyan("🧿 Probing: %s", target)
		core.RunProbe(target)
		color.Green("Scan complete!")
	},
}
