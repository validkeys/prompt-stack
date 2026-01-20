package main

import (
	"fmt"
	"os"

	"github.com/kyledavis/prompt-stack/pkg/executor"
	"github.com/spf13/cobra"
)

var ralphyDryRun bool

var ralphyCmd = &cobra.Command{
	Use:   "ralphy",
	Short: "Run Ralphy shell executor",
	Long:  `Execute Ralphy shell script for AI-assisted task execution. Use --dry-run to generate reports without executing.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if ralphyDryRun {
			return runRalphyDryRun()
		}
		return runRalphyLive()
	},
}

func runRalphyLive() error {
	return fmt.Errorf("live execution not implemented yet, use --dry-run flag")
}

func init() {
	rootCmd.AddCommand(ralphyCmd)
	ralphyCmd.Flags().BoolVar(&ralphyDryRun, "dry-run", false, "Generate reports without executing")
}

func runRalphyDryRun() error {
	fmt.Println("=== Ralphy Dry-Run Mode ===")
	fmt.Println()

	execr := executor.NewExecutor(".", true)

	config := executor.ExecutionConfig{
		Task:       "dry-run",
		AIEngine:   "opencode",
		DryRun:     true,
		WorkingDir: ".",
	}

	result, err := execr.Execute(config)
	if err != nil {
		return fmt.Errorf("dry-run execution failed: %w", err)
	}

	if !result.Success {
		fmt.Println("Dry-run validation failed:")
		fmt.Println(result.Stderr)
		return fmt.Errorf("dry-run validation failed")
	}

	fmt.Println("✓ Dry-run completed successfully")
	fmt.Println()
	fmt.Println("Generated files:")
	fmt.Printf("  - .prompt-stack/report.txt\n")
	fmt.Printf("  - .prompt-stack/audit.log\n")
	fmt.Println()

	if _, err := os.Stat(".prompt-stack/report.txt"); err == nil {
		fmt.Println("Report preview (first 200 chars):")
		content, _ := os.ReadFile(".prompt-stack/report.txt")
		if len(content) > 200 {
			fmt.Println(string(content[:200]) + "...")
		} else {
			fmt.Println(string(content))
		}
		fmt.Println()
	}

	return nil
}
