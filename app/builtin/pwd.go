package builtin

import (
	"fmt"
	"log"
	"os"

	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

type Pwd struct{}

func NewPwd() *Pwd {
	return &Pwd{}
}

func (p *Pwd) Run(cmd *utils.ShellCommand) int {
	wd, err := os.Getwd(); 
	if err != nil {
		log.Printf("Unable to get current working directory: %v", err)
		return 1
	}
	fmt.Println(wd)
	return -1
}