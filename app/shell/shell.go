package shell

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/builtin"
	"github.com/codecrafters-io/shell-starter-go/app/executable"
	"github.com/codecrafters-io/shell-starter-go/app/linereader"
	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

type shell struct {
	builtins map[string]builtin.Builtin
	executables []string
}

func New() shell {
	return shell{
		builtins: getBuiltins(),
		executables: getExecutables(),
	}
}

func getBuiltins() map[string]builtin.Builtin {
	builtins := map[string]builtin.Builtin{}
	builtins["exit"] = builtin.Exit{}
	builtins["echo"] = builtin.Echo{}
	builtins["type"] = builtin.Type{Builtins: &builtins}
	builtins["pwd"] = builtin.Pwd{}
	builtins["cd"] = builtin.Cd{}
	return builtins
}

func getExecutables() []string {
	path := os.Getenv("PATH")
	dirs := strings.Split(path, string(os.PathListSeparator))
	execs := []string{}
	for _, dir := range dirs {
		contents, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, obj := range contents {
			if isExec(obj) {
				execs = append(execs, obj.Name())
			}
		}
	}
	return execs
}

func isExec(obj os.DirEntry) bool {
	info, err := obj.Info()
	if err != nil {
		return false
	}
	mode := info.Mode()
	return mode.IsRegular() && mode.Perm()&0111 != 0
}


func (s shell) Run() {
	var invoker string
	var userInput string = ""
	isNewLine := true
	reader := linereader.New(bufio.NewReader(os.Stdin))
	for {
		if isNewLine {
			fmt.Fprint(os.Stdout, "$ ")
			isNewLine = false
		}
		var err error
		invoker, userInput, err = reader.ReadLine(userInput)
		if err != nil {
			log.Fatalf("Error whilst reading user input: %v", err)
		}
		switch invoker {
		case "\n":
			s.executeCommand(userInput)
			isNewLine, userInput = true, ""
		case "\t":
			endOfWord := s.autocomplete(userInput)
			fmt.Fprint(os.Stdout, endOfWord)
			userInput += endOfWord
		}		
	}
}

func (s shell) executeCommand(userInput string) {
	var config utils.ShellConfig
	fmt.Fprint(os.Stdout, "\n")
	userInputSplit := utils.ParseString(userInput)
	command := userInputSplit[0]
	args := []string{}
	if len(userInputSplit) > 1 {
		args = userInputSplit[1:]
	}
	args, config = utils.ParseArgs(args)
	if b, ok := s.builtins[command]; ok {
		b.Run(args, config)
	} else if path := utils.FindExecutablePath(command); path != "" {
		executable.RunExecutable(command, config, args)
	} else {
		fmt.Fprint(config.StdErrFile, command + ": command not found\n")
	}
}

func (s shell) autocomplete(userInput string) string {
	for builtin := range s.builtins {
		if ok, completionStr := canAutocomplete(builtin, userInput); ok {
			return completionStr
		}
	}
	for _, exec := range s.executables {
		if ok, completionStr := canAutocomplete(exec, userInput); ok {
			return completionStr
		}
	}
	fmt.Fprint(os.Stdout, "\x07")
	return ""
}

func canAutocomplete(s string, prefix string) (bool, string) {
	if strings.HasPrefix(s, prefix) {
		return true, strings.TrimPrefix(s, prefix) + " "
	}
	return false, ""
}