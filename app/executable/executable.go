package executable

import (
	"log"
	"os/exec"

	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

func Run(s *utils.ShellCommand) int {
	cmd := exec.Command(s.Command, s.Args...)
	cmd.Stdin = s.StdInFile
	cmd.Stdout = s.StdOutFile
	cmd.Stderr = s.StdErrFile
	if err := cmd.Start(); err != nil {
		log.Printf("Error running command %s with args %s: %v", s.Command, s.Args, err)
	}
	cmd.Wait()
	return -1
}

