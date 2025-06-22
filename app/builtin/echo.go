package builtin

import (
	"fmt"
	"strings"
	"sync"

	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

type Echo struct{}

func (e Echo) Run(cmd *utils.ShellCommand, wg *sync.WaitGroup) {
	defer cmd.Close()
	defer wg.Done()
	fmt.Fprintln(cmd.StdOutFile, strings.Join(cmd.Args, " "))
}