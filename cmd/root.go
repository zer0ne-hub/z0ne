// root command: when z0ne is called without any arguments
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
	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Print banner and quick usage
		fmt.Println("Please specify a module command. Run 'z0ne --help' for usage.")

		// Print version when the --version flag is set
		showVersion, _ := cmd.Flags().GetBool("version")
		if showVersion {
			fmt.Println("Version: v0.1.0") // TODO: Get version from release source
			return
		}
	},
}

func init() {
	rootCmd.Flags().Bool("version", false, "Show version info")
	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(probeCmd)
	rootCmd.AddCommand(reportCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error executing command: ", err)
		os.Exit(1)
	}
}
