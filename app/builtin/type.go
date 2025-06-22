package builtin

import (
	"fmt"
	"sync"

	"github.com/codecrafters-io/shell-starter-go/app/utils"
)


type Type struct {
	Builtins *map[string]Builtin
}

func (t Type) Run(cmd *utils.ShellCommand, wg *sync.WaitGroup) {
	defer cmd.Close()
	defer wg.Done()
	if len(cmd.Args) == 0 {
		fmt.Println("Usage: type <cmd>")
		return
	}
	command := cmd.Args[0]
	if t.isBuiltin(command) {
		fmt.Fprintf(cmd.StdOutFile, "%s is a shell builtin\n", command)
	} else if ok, path := t.isExecutable(command); ok {
		fmt.Fprintf(cmd.StdOutFile, "%s is %s\n", command, path)
	} else {
		fmt.Fprintf(cmd.StdOutFile, "%s: not found\n", command)
	}
}

func (t Type) isBuiltin(command string) bool {
	builtins := *t.Builtins
	if _, ok := builtins[command]; ok {
		return true
	}
	return false
}

func (t Type) isExecutable (command string) (bool, string) {
	if path := utils.FindExecutablePath(command); path != "" {
		return true, path
	}
	return false, ""
}