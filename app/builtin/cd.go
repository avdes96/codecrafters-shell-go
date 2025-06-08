package builtin

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"
)

type Cd struct{}

func (c Cd) Run(args []string) {
	if len(args) != 1 {
		fmt.Println("Usage: cd <dir>")
		return
	}
	newDir := args[0]
	if newDir == "~" || strings.HasPrefix(newDir, "~" + string(os.PathSeparator)) {
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Unable to find user home dir: %v", err)
		}
		newDir = strings.Replace(newDir, "~", userHomeDir, 1)
	}
	err := os.Chdir(newDir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			fmt.Println("cd: " + newDir + ": No such file or directory")
			return
		}
		log.Fatalf("Error in changing dir: %v", err)
	}
}