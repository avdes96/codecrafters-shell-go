package builtin

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

type Cd struct{}

func NewCd() *Cd {
	return &Cd{}
}

func (c *Cd) Run(cmd *utils.ShellCommand) int {
	if len(cmd.Args) != 1 {
		fmt.Println("Usage: cd <dir>")
		return -1
	}
	newDir := cmd.Args[0]
	if newDir == "~" || strings.HasPrefix(newDir, "~" + string(os.PathSeparator)) {
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			log.Printf("Unable to find user home dir: %v", err)
			return 1
		}
		newDir = strings.Replace(newDir, "~", userHomeDir, 1)
	}
	err := os.Chdir(newDir)
	if err == nil {
		return -1
	}
	if errors.Is(err, fs.ErrNotExist) {
		fmt.Println("cd: " + newDir + ": No such file or directory")
		return -1
	}
	log.Printf("Error in changing dir: %v", err)
	return 1
}