package builtin

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

type Echo struct{}

func NewEcho() *Echo {
	return &Echo{}
}

func (e *Echo) Run(cmd *utils.ShellCommand) int {
	fmt.Fprintln(cmd.StdOutFile, strings.Join(cmd.Args, " "))
	return -1
}