package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build project from implementation plan",
	Long:  `Build project components based on implementation plan tasks.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("build command: Build project from implementation plan")
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
