// Package cmd: Handles command line interface using cobra
package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"

	"github.com/spf13/cobra"
)

// rootCmd: represents the base command when the user runs z0ne without any arguments
var rootCmd = &cobra.Command{
	Use:   "z0ne",
	Short: "💀 z0ne - Attack Surface Mapper",
	Long:  `z0ne is a modular attack surface mapper designed for CTFs and pentesting.`,
	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please specify a module command. Run 'z0ne --help' for quick usage of available modules.")
		fmt.Println("For help with a specific module or command, Run 'z0ne <module> --help'.")
		// Print version when the --version flag is set
		showVersion, _ := cmd.Flags().GetBool("version")
		if showVersion {
			color.Cyan("Version: v0.1.0") // TODO: Get version from release source
			return
		}
	},
}

// init: Initialize the root command with flags and subcommands
func init() {
	rootCmd.Flags().Bool("version", false, "Show version info")
	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(probeCmd)
	rootCmd.AddCommand(reportCmd)
}

// Execute: Executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		color.Red("Error executing command: %v", err)
		os.Exit(1)
	}
}
