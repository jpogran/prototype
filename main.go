package main

import (
	"os"

	"github.com/jpogran/prototype/cmd"
)

// osExit is a copy of `os.Exit` to ease the "exit status" test.
// See: https://stackoverflow.com/a/40801733/8367711
var osExit = os.Exit

func main() {
	if err := cmd.Execute(); err != nil {
		osExit(1)
	}
}
