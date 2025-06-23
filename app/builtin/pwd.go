package builtin

import (
	"fmt"
	"log"
	"os"

	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

type Pwd struct{}

func (p Pwd) Run(cmd *utils.ShellCommand) {
	wd, err := os.Getwd(); 
	if err != nil {
		log.Fatalf("Unable to get current working directory: %v", err)
	}
	fmt.Println(wd)
}