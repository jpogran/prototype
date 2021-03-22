package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = createRootCommand()

func createRootCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:     "prototype",
		Short:   "PDK prototype content template system",
		Long:    `PDK prototype content template system`,
		Version: "3.0.0",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Working with OutOrStdout/OutOrStderr allows us to unit test our command easier
			// out := cmd.OutOrStdout()
		},
	}
	return tmp
}

func init() {}
