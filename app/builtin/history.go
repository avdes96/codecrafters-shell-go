package builtin

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

type History struct{
	HistoryList *[]string
	lastAppended int
}

const usageStr = "Usage: history [limit] | [-r <path>]"

func (h *History) Run(cmd *utils.ShellCommand) int {
	switch len(cmd.Args) {
	case 0:
		n := len(*h.HistoryList)
		h.printHistory(cmd.StdOutFile, n)
		return -1
	case 1:
		nStr := cmd.Args[0]
		n, err := strconv.Atoi(nStr)
		if err != nil {
			fmt.Fprintf(cmd.StdOutFile, "history: %s: numeric argument required\n", nStr)
			return -1
		}
		if n < 0 {
			fmt.Fprintf(cmd.StdOutFile, "history: %d: numeric argument must be non-negative\n", n)
			return -1
		}
		h.printHistory(cmd.StdOutFile, n)
		return -1
	case 2:
		switch cmd.Args[0] {
		case "-r":
			err := h.ReadFromFile(cmd.Args[1], false)
			if err != nil {
				log.Printf("Error reading from file %s: %s", cmd.Args[1], err)
				return -1
			}
		case "-w":
			err := h.writeToFile(cmd.Args[1])
			if err != nil {
				log.Printf("Error writing to file %s: %s", cmd.Args[1], err)
				return -1
			}
		case "-a":
			err := h.AppendToFile(cmd.Args[1])
			if err != nil {
				log.Printf("Error appending to file %s: %s", cmd.Args[1], err)
				return -1
			}
		default:
			fmt.Fprint(cmd.StdOutFile, usageStr)
			return -1
		}
	default:
		fmt.Fprint(cmd.StdOutFile, usageStr)
		return -1
	}
	return -1
}

func NewHistory(h *[]string) *History {
	return &History{
		HistoryList: h,
		lastAppended: 0,
	}
}

func (h *History) printHistory(out *os.File, n int) {
	start := max(0, len(*h.HistoryList) - n)
	for i, prevCommand := range (*h.HistoryList)[start:] {
		fmt.Fprintf(out, "\t%d %s\n", i+1, prevCommand)
	}
}

func (h *History) ReadFromFile(filename string, onStart bool) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("Unable to open file: %s", err)
	}
	buffer := make([]byte, 1024)
	s := ""
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("Unable to read file: %s", err)
		}
		s += strings.TrimSpace(string(buffer[:n]))
	}
	cmds := strings.Split(s, "\n")
	cmdsFiltered := []string{}
	for _, cmd := range cmds {
		if c := strings.TrimSpace(cmd); c != "" {
			cmdsFiltered = append(cmdsFiltered, c)
		}
	}
	*h.HistoryList = append(*h.HistoryList, cmdsFiltered...)
	if onStart {
		h.lastAppended = len(*h.HistoryList)
	}
	return nil
}

func (h *History) writeToFile(filename string) error {
	file := utils.GetOSFile(filename, true)
	for _, cmd := range *h.HistoryList {
		_, err := file.WriteString(cmd + "\n")
		if err != nil {
			return fmt.Errorf("Unable to write to file")
		}
	}
	return nil
}

func (h *History) AppendToFile(filename string) error {
	file := utils.GetOSFile(filename, false)
	for i, cmd := range (*h.HistoryList)[h.lastAppended:] {
		_, err := file.WriteString(cmd + "\n")
		if err != nil {
			h.lastAppended = h.lastAppended + i
			return fmt.Errorf("Unable to append to file")
		}
	}
	h.lastAppended = len(*h.HistoryList)
	return nil
}