package builtin

import (
	"fmt"
	"log"
	"os"
)

type Pwd struct{}

func (p Pwd) Run(args []string) {
	wd, err := os.Getwd(); 
	if err != nil {
		log.Fatalf("Unable to get current working directory: %v", err)
	}
	fmt.Println(wd)
}