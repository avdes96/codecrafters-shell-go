package builtin

import (
	"fmt"
	"os/exec"
)


type Type struct {
	Builtins *map[string]Builtin
}

func (t Type) Run(args []string) {
	command := args[0]
	if t.isBuiltIn(command) {
		fmt.Println(command + " is a shell builtin")
	} else if ok, exec := t.isExecutable(command); ok {
		fmt.Println(command + " is " + exec)
	} else {
		fmt.Println(command + ": not found")
	}
}

func (t Type) isBuiltIn(command string) bool {
	builtins := *t.Builtins
	if _, ok := builtins[command]; ok {
		return true
	}
	return false
}

func (t Type) isExecutable (command string) (bool, string) {
	path, err := exec.LookPath(command)
	if err != nil {
		return false, ""
	}
	return true, path
}