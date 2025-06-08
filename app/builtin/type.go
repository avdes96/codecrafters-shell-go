package builtin

import (
	"fmt"
	"os"
	"strings"
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
	path := os.Getenv("PATH")
	dirs := strings.Split(path, string(os.PathListSeparator))
	for _, dir := range dirs {
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, file := range files {
			if file.Name() == command {
				return true, dir + string(os.PathSeparator) + file.Name()
			}
		}
	}

	return false, ""
}