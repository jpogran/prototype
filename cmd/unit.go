package cmd

import (
	"fmt"

	"github.com/jpogran/prototype/internal/pkg/pdkshell"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	cleanFixtures          bool
	listUnitTestFiles      bool
	parallelUnitTests      bool
	puppetDevSourceVersion string
	puppetVersion          string
	unitTestsToRun         string
	verboseUnitTestOutput  bool
)

func CreateTestUnitCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "unit",
		Short: "Run unit tests",
		Long:  `Run unit tests`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			logrus.Trace("test unit PreRunE")
			bindAllFlagsForCommand(cmd)
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			logrus.Trace("test unit Run")
			argsV := []string{"test", "unit"}

			flagsToIgnore := []string{"log-level"}

			argsV = getListOfFlags(cmd, argsV, flagsToIgnore)

			logrus.Tracef("OutputFormat: %v\n", format)
			logrus.Tracef("args: %v", argsV)

			shell := pdkshell.New()
			shell.Execute(argsV)
		},
	}
	tmp.Flags().BoolVarP(&cleanFixtures, "clean-fixtures", "c", false, "clean up downloaded fixtures after the test run")
	tmp.Flags().BoolVar(&listUnitTestFiles, "list", false, "list all available unit test files")
	tmp.Flags().BoolVar(&parallelUnitTests, "parallel", false, "run unit tests in parallel")

	tmp.Flags().StringVar(&puppetDevSourceVersion, "puppet-dev", "", "When specified, PDK will validate or test against the current Puppet source from github.com. To use this option, you must have network access to https://github.com.")
	tmp.Flags().StringVar(&puppetVersion, "puppet-version", "", "Puppet version to run tests or validations against")
	tmp.Flags().StringVar(&unitTestsToRun, "tests", "", "Specify a comma-separated list of unit test files to run. (default: )")

	tmp.Flags().BoolVar(&verboseUnitTestOutput, "verbose", false, "more verbose --list output. displays a list of examples in each unit test file")
	return tmp
}

func init() {

}

func getListOfFlags(cmd *cobra.Command, argsV []string, flagsToIgnore []string) []string {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if !contains(flagsToIgnore, f.Name) {
			if f.Changed {
				switch f.Value.Type() {
				case "bool":
					argsV = append(argsV, fmt.Sprintf("--%v", f.Name))
				case "string":
					argsV = append(argsV, fmt.Sprintf("--%v %v", f.Name, f.Value))
				}
			}
		}
	})
	return argsV
}
