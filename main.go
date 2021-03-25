package main

import (
	"github.com/jpogran/prototype/cmd"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	rootCmd := cmd.CreateRootCommand()
	rootCmd.AddCommand(cmd.CreateCompletionCommand())
	rootCmd.AddCommand(cmd.CreateNewCommand())

	testCmd := cmd.CreateTestCommand()
	unitCmd := cmd.CreateTestUnitCommand()
	testCmd.AddCommand(unitCmd)

	rootCmd.AddCommand(testCmd)

	version := cmd.Format(version, date, commit)
	versionCmd := cmd.NewCmdVersion(version, date, commit)
	rootCmd.Version = version
	rootCmd.SetVersionTemplate(version)
	rootCmd.AddCommand(versionCmd)

	cobra.OnInitialize(cmd.InitConfig)

	cobra.CheckErr(rootCmd.Execute())
}
