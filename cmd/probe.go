// Package cmd: Handles command line interface using cobra
package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zer0ne-hub/z0ne/internal/core"
)

// probeCmd: represents the probe command
var probeCmd = &cobra.Command{
	Use:   "probe [target]",
	Short: "more precise web reconnaissance",
	Long: `This mode will enumerate subdomains, resolve IP addresses, and scan ports... on a given target but
with more precision, using authentication, stealth, and some basic bypass methods. Might be slower,
	require higher privileges and support flags for authentication to some tools.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		color.Cyan("🧿 Probing: %s", target)

		currentProbekeys := core.ProbeKeys{}
		currentProbekeys.ShodanKey, _ = cmd.Flags().GetString("shodan")
		err := core.RunProbe(target, currentProbekeys)
		if err != nil {
			color.Red("Probe Error: %s", err)
		}
		color.Green("\n\nProbe complete! Results saved to /z0ne-out/results.json")
		color.Green("You can now run 'z0ne report [targetName]' to generate a MD report.")
	},
}

// init: Initialize the probe command with flags
func init() {
	probeCmd.Flags().String("shodan", "", "Shodan API key")
}
