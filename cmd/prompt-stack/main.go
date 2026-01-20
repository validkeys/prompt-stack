package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// osExit is overridden in tests to avoid exiting the test process
var osExit = func(code int) { os.Exit(code) }

var rootCmd = &cobra.Command{
	Use:   "prompt-stack",
	Short: "AI-assisted development workflow tool",
	Long:  `A tool for generating and validating Ralphy YAML files with Plan/Build modes.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("AI-assisted development workflow tool")
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.Version = fmt.Sprintf("%s (commit: %s, built: %s)", Version, Commit, Date)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		osExit(1)
	}
}
