package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var reconCmd = &cobra.Command{
	Use:   "recon [target]",
	Short: "Complete target reconnaissance",
	Long:  `This mode will enumerate subdomains, resolve IP addresses, and scan ports on a given target.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		fmt.Println("👀 Running RECON mode on:", target)
	},
}

func init() {
	rootCmd.AddCommand(reconCmd)
}
