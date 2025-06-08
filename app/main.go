package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/builtin"
	"github.com/codecrafters-io/shell-starter-go/app/logger"
)



func main() {
	logger.InitLogger()

	builtins := map[string]builtin.Builtin{}
	builtins["exit"] = builtin.Exit{}
	builtins["echo"] = builtin.Echo{}
	builtins["type"] = builtin.Type{Builtins: &builtins}
		
	for {
		fmt.Fprint(os.Stdout, "$ ")
		userInput, err := bufio.NewReader(os.Stdin).ReadString('\n')
		
		if err != nil {
			log.Fatalf("Error whilst reading user input: %s", err)
		}
		userInput = userInput[:len(userInput)-1]
		userInputSplit := strings.Split(userInput, " ")
		command := userInputSplit[0]
		args := []string{}
		if len(userInputSplit) > 1 {
			args = userInputSplit[1:]
		}
		if b, ok := builtins[command]; ok {
			b.Run(args)
		} else {
			fmt.Println(command + ": command not found")
		}
	}
}
