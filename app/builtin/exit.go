package builtin

import (
	"log"
	"os"
	"strconv"

	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

type Exit struct{}

func NewExit() *Exit {
	return &Exit{}
}

func (e *Exit) Run(cmd *utils.ShellCommand) {
	var exitCode int
	if len(cmd.Args) == 0 {
		exitCode = 0
	} else {
		exitCodeInt, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			log.Fatalf("%s is an invalid exit code: %v", cmd.Args[1], err)
		}
		exitCode = exitCodeInt
	}
	os.Exit(exitCode)
}