package builtin

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

type Echo struct{}

func (e Echo) Run(cmd *utils.ShellCommand) {
	fmt.Fprintln(cmd.StdOutFile, strings.Join(cmd.Args, " "))
}