package builtin

import "github.com/codecrafters-io/shell-starter-go/app/utils"

type Builtin interface {
	Run(args []string, config utils.ShellConfig)
}