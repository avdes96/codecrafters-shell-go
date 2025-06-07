package builtin

import (
	"fmt"
	"strings"
)

type Echo struct {}

func (e Echo) Run(args []string) {
	fmt.Println(strings.Join(args, " "))
}