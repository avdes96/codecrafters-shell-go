package builtin

import (
	"log"
	"os"
	"strconv"

	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

type Exit struct{}

func (e Exit) Run(args []string, config utils.ShellConfig) {
	var exitCode int
	if len(args) == 0 {
		exitCode = 0
	} else {
		exitCodeInt, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalf("%s is an invalid exit code: %v", args[1], err)
		}
		exitCode = exitCodeInt
	}
	os.Exit(exitCode)
}