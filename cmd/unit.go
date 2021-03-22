package cmd

import (
	"fmt"

	"github.com/jpogran/prototype/internal/pkg/pdkshell"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var testCmd = createTestCommand()

func init() {
	rootCmd.AddCommand(testCmd)
}

func createTestCommand() *cobra.Command {
	tmpTest := &cobra.Command{
		Use:       "test",
		Short:     "run tests",
		Long:      `run tests`,
		Args:      cobra.ExactValidArgs(1),
		ValidArgs: []string{"unit"},
		Run: func(cmd *cobra.Command, args []string) {
			argsV := []string{"test", args[0]}
			cmd.Flags().VisitAll(func(f *pflag.Flag) {
				if f.Changed {
					switch f.Value.Type() {
					case "bool":
						argsV = append(argsV, fmt.Sprintf("--%v", f.Name))
					case "string":
						argsV = append(argsV, fmt.Sprintf("--%v %v", f.Name, f.Value))
					}
				}
			})
			fmt.Printf("Config.UnitTest.OutputFormat: %v\n", Config.UnitTest.OutputFormat)
			fmt.Printf("args: %v", argsV)
			shell := pdkshell.New()
			shell.Execute(argsV)
		},
	}
	tmpTest.Flags().BoolVarP(&Config.UnitTest.TestCmdDebug, "debug", "d", false, "enable debug output")

	tmpTest.Flags().BoolVarP(&Config.UnitTest.UnitTestCleanFixtures, "clean-fixtures", "c", false, "clean up downloaded fixtures after the test run")
	tmpTest.Flags().BoolVar(&Config.UnitTest.ListUnitTestFiles, "list", false, "list all available unit test files")
	tmpTest.Flags().BoolVar(&Config.UnitTest.ParallelUnitTests, "parallel", false, "run unit tests in parallel")

	tmpTest.Flags().StringVar(&Config.UnitTest.PuppetDevSourceVersion, "puppet-dev", "", "When specified, PDK will validate or test against the current Puppet source from github.com. To use this option, you must have network access to https://github.com.")
	tmpTest.Flags().StringVar(&Config.UnitTest.PuppetVersion, "puppet-version", "", "Puppet version to run tests or validations against")
	tmpTest.Flags().StringVar(&Config.UnitTest.UnitTestsToRun, "tests", "", "Specify a comma-separated list of unit test files to run. (default: )")

	tmpTest.Flags().StringVarP(&Config.UnitTest.OutputFormat, "format", "f", "", "formating (default is junit)")
	tmpTest.Flags().BoolVar(&Config.UnitTest.VerboseUnitTestOutput, "verbose", false, "more verbose --list output. displays a list of examples in each unit test file")

	return tmpTest
}
