package main

import (
	"github.com/spf13/cobra"
)

var helpCmd = &cobra.Command{
	Use:   "help [command]",
	Short: "Help about any command",
	Long: `Help provides help for any command in the CLI.
Run "prompt-stack help [command]" for more information about a command.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Parent().Help()
			return
		}

		c, _, err := cmd.Parent().Find(args)
		if err != nil {
			cmd.Printf("Unknown help topic: %s\n", args[0])
			cmd.Parent().Help()
			return
		}

		c.Help()
	},
}

func init() {
	rootCmd.SetHelpCommand(helpCmd)
}
