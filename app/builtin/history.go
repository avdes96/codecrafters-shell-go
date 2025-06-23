package builtin

import (
	"fmt"

	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

type History struct{HistoryList *[]string}

func (h *History) Run(cmd *utils.ShellCommand) {
	for i, prevCommand := range *h.HistoryList {
		fmt.Fprintf(cmd.StdOutFile, "\t%d %s\n", i+1, prevCommand)
	}
}

func NewHistory(h *[]string) *History {
	return &History{HistoryList: h}
}