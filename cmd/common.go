package cmd

import (
	"fmt"
	"os"
)

// osExit is a copy of `os.Exit` to ease the "exit status" test.
// See: https://stackoverflow.com/a/40801733/8367711
var osExit = os.Exit

// EchoStdErrIfError is an STDERR wrappter and returns 0(zero) or 1.
// It does nothing if the error is nil and returns 0.
func EchoStdErrIfError(err error) int {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)

		return 1
	}

	return 0
}

// Execute is the main function of `cmd` package.
func Execute() error {
	return rootCmd.Execute()
}
