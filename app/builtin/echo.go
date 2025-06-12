package builtin

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

type Echo struct{}

func (e Echo,) Run(args []string, config utils.ShellConfig) {
	fmt.Fprintln(config.StdOutFile, strings.Join(args, " "))
}