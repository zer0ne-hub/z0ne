// cmd/root.go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "z0ne",
	Short: "💀 z0ne - Attack Surface Mapper",
	Long:  `z0ne is a modular attack surface mapper designed for CTFs and pentesting.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please specify a command. Run 'z0ne --help' for usage.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error executing command: ", err)
		os.Exit(1)
	}
}
