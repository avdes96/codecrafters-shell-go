package builtin

import (
	"fmt"

	"github.com/codecrafters-io/shell-starter-go/app/utils"
)


type Type struct {
	Builtins *map[string]Builtin
}

func (t Type) Run(args []string, config utils.ShellConfig) {
	command := args[0]
	if t.isBuiltin(command) {
		fmt.Println(command + " is a shell builtin")
	} else if ok, path := t.isExecutable(command); ok {
		fmt.Println(command + " is " + path)
	} else {
		fmt.Println(command + ": not found")
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