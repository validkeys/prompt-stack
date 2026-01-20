package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Generate implementation plans",
	Long:  `Generate implementation plans from requirements or templates.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("plan command: Generate implementation plans")
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(planCmd)
}
