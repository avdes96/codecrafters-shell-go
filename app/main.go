package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

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
		
	for {
		fmt.Fprint(os.Stdout, "$ ")
		userInput, err := bufio.NewReader(os.Stdin).ReadString('\n')
		
		if err != nil {
			log.Fatalf("Error whilst reading user input: %v", err)
		}
		userInput = userInput[:len(userInput)-1]
		userInputSplit := strings.SplitN(userInput, " ", 2)
		command := userInputSplit[0]
		argsStr := ""
		if len(userInputSplit) > 1 {
			argsStr = userInputSplit[1]
		}
		args := utils.ParseArgs(argsStr)
		if b, ok := builtins[command]; ok {
			b.Run(args)
		} else if path := utils.FindExecutablePath(command); path != "" {
			executable.RunExecutable(command, args)
		} else {
			fmt.Println(command + ": command not found")
		}
	}
}
