package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run tests",
	Long:  `Run tests`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Trace("test called")
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "enable debug output")
	testCmd.PersistentFlags().StringVarP(&format, "format", "f", "junit", "formating (default is junit)")
}
