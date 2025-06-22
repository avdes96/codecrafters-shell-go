package builtin

import "github.com/codecrafters-io/shell-starter-go/app/utils"

type Builtin interface {
	Run(cmd *utils.ShellCommand)
}