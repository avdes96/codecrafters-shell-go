package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)


func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")
		userInput, err := bufio.NewReader(os.Stdin).ReadString('\n')
		
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error whilst reading user input: %s", err)
			os.Exit(1)
		}
		userInput = userInput[:len(userInput)-1]
		userInputSplit := strings.Split(userInput, " ")
		command := userInputSplit[0]
		args := []string{}
		if len(userInputSplit) > 1 {
			args = userInputSplit[1:]
		}
		
		switch command {
		case "exit":
			exitCode, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s is an invalid exit code: %s", args[1], err)
				os.Exit(1)
			}
			os.Exit(exitCode)
		case "echo":
			fmt.Println(strings.Join(args, " "))
		default: 
			fmt.Println(command + ": command not found")
		}
	}
}
