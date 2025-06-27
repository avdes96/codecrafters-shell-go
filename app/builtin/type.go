package builtin

import (
	"fmt"

	"github.com/codecrafters-io/shell-starter-go/app/utils"
)


type Type struct {
	Builtins *map[string]Builtin
}

func NewType(b *map[string]Builtin) *Type {
	return &Type{Builtins: b}
}

func (t *Type) Run(cmd *utils.ShellCommand) int {
	if len(cmd.Args) == 0 {
		fmt.Println("Usage: type <cmd>")
		return -1
	}
	command := cmd.Args[0]
	if t.isBuiltin(command) {
		fmt.Fprintf(cmd.StdOutFile, "%s is a shell builtin\n", command)
	} else if ok, path := t.isExecutable(command); ok {
		fmt.Fprintf(cmd.StdOutFile, "%s is %s\n", command, path)
	} else {
		fmt.Fprintf(cmd.StdOutFile, "%s: not found\n", command)
	}
	return -1
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