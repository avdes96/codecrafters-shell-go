package builtin

import "fmt"

type Type struct {
	Builtins *map[string]Builtin
}

func (t Type) Run(args []string) {
	builtins := *t.Builtins
	if _, ok := builtins[args[0]]; ok {
		fmt.Println(args[0] + " is a shell builtin")
		return
	}
	fmt.Println(args[0] + ": not found")
}