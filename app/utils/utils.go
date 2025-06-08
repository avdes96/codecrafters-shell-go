package utils

import "os/exec"

func FindExecutablePath(executable string) string {
	path, err := exec.LookPath(executable)
	if err != nil {
		return ""
	}
	return path
}