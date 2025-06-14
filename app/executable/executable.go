package executable

import (
	"log"
	"os/exec"

	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

func RunExecutable(path string, config utils.ShellConfig, args []string) {
	cmd := exec.Command(path, args...)
	cmd.Stdout = config.StdOutFile
	cmd.Stderr = config.StdErrFile
	if err := cmd.Run(); err != nil {
		log.Printf("Error running command %s: %v", path, err)
	}
}