package builtin

import (
	"fmt"
	"os"
	"strconv"
)

type Exit struct{}

func (e Exit) Run(args []string) {
	var exitCode int
	if len(args) == 0 {
		exitCode = 0
	} else {
		exitCodeInt, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s is an invalid exit code: %s", args[1], err)
			os.Exit(1)
		}
		exitCode = exitCodeInt
	}
	os.Exit(exitCode)
}