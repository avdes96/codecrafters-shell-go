package executable

import (
	"log"
	"os"
	"os/exec"
)

func RunExecutable(path string, args []string) {
	cmd := exec.Command(path, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Error running command %s: %v", path, err)
	}
}