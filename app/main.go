package main

import (
	"bufio"
	"fmt"
	"os"
)


func main() {
	running := true
	for running {
		fmt.Fprint(os.Stdout, "$ ")
		// Wait for user input
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error whilst reading user input: %s", err)
			os.Exit(1)
		}
		fmt.Println(command[:len(command)-1] + ": command not found")
	}
}
