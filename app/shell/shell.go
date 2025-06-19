package shell

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/codecrafters-io/shell-starter-go/app/builtin"
	"github.com/codecrafters-io/shell-starter-go/app/executable"
	"github.com/codecrafters-io/shell-starter-go/app/linereader"
	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

type shell struct {
	builtins map[string]builtin.Builtin
}

func New() shell {
	return shell{
		builtins: getBuiltins(),
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
	if userInput == "ech" {
		return "o "
	}
	return "t "
}