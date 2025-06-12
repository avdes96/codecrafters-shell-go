package builtin

import (
	"fmt"
	"log"
	"os"

	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

type Pwd struct{}

func (p Pwd) Run(args []string, config utils.ShellConfig) {
	wd, err := os.Getwd(); 
	if err != nil {
		log.Fatalf("Unable to get current working directory: %v", err)
	}
	fmt.Println(wd)
}