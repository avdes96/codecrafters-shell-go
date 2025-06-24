package builtin

import (
	"fmt"
	"os"
	"strconv"

	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

type History struct{HistoryList *[]string}

func (h *History) Run(cmd *utils.ShellCommand) {
	if len(cmd.Args) > 1 {
		fmt.Fprint(cmd.StdOutFile, "Usage: history <n>\n")
		return
	}
	if len(cmd.Args) == 0 {
		n := len(*h.HistoryList)
		h.printHistory(cmd.StdOutFile, n)
		return
	}

	nStr := cmd.Args[0]
	n, err := strconv.Atoi(nStr)
	if err != nil {
		fmt.Fprintf(cmd.StdOutFile, "history: %s: numeric argument required\n", nStr)
		return
	}
	if n < 0 {
		fmt.Fprintf(cmd.StdOutFile, "history: %d: numeric argument must be non-negative\n", n)
		return
	}
	h.printHistory(cmd.StdOutFile, n)
	return
}

func NewHistory(h *[]string) *History {
	return &History{HistoryList: h}
}

func (h *History) printHistory(out *os.File, n int) {
	start := max(0, len(*h.HistoryList) - n)
	for i, prevCommand := range (*h.HistoryList)[start:] {
		fmt.Fprintf(out, "\t%d %s\n", i+1, prevCommand)
	}
}