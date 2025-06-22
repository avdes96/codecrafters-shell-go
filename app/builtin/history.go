package builtin

import (
	"sync"

	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

type History struct{}

func (h History) Run(cmd *utils.ShellCommand, wg *sync.WaitGroup) {
	
}