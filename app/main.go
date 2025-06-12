package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/codecrafters-io/shell-starter-go/app/builtin"
	"github.com/codecrafters-io/shell-starter-go/app/executable"
	"github.com/codecrafters-io/shell-starter-go/app/logger"
	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

func main() {
	logger.InitLogger()

	builtins := map[string]builtin.Builtin{}
	builtins["exit"] = builtin.Exit{}
	builtins["echo"] = builtin.Echo{}
	builtins["type"] = builtin.Type{Builtins: &builtins}
	builtins["pwd"] = builtin.Pwd{}
	builtins["cd"] = builtin.Cd{}
	
	var config utils.ShellConfig

	for {
		fmt.Fprint(os.Stdout, "$ ")
		userInput, err := bufio.NewReader(os.Stdin).ReadString('\n')
		
		if err != nil {
			log.Fatalf("Error whilst reading user input: %v", err)
		}
		userInput = userInput[:len(userInput)-1]
		userInputSplit := utils.ParseString(userInput)
		command := userInputSplit[0]
		args := []string{}
		if len(userInputSplit) > 1 {
			args = userInputSplit[1:]
		}
		args, config = utils.ParseArgs(args)
		if b, ok := builtins[command]; ok {
			b.Run(args, config)
		} else if path := utils.FindExecutablePath(command); path != "" {
			executable.RunExecutable(command, config, args)
		} else {
			fmt.Fprint(config.StdErrFile, command + ": command not found\n")
		}
	}
}
