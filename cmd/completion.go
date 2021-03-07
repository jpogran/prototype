package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:       "completion",
	Short:     "Generate shell completions for the chosen shell",
	Long:      `Generate shell completions for the chosen shell`,
	ValidArgs: []string{"bash", "fish", "pwsh", "zsh"},
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		switch args[0] {
		case "bash":
			err = cmd.Root().GenBashCompletion(os.Stdout)
		case "fish":
			err = cmd.Root().GenFishCompletion(os.Stdout, true)
		case "pwsh":
			err = cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		case "zsh":
			err = cmd.Root().GenZshCompletion(os.Stdout)
		default:
			log.Printf("unsupported shell type %q", args[0])
		}

		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
