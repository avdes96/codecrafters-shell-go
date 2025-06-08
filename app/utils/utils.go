package utils

import (
	"os/exec"
	"unicode"
)

func FindExecutablePath(executable string) string {
	path, err := exec.LookPath(executable)
	if err != nil {
		return ""
	}
	return path
}


func ParseArgs(s string) []string {
	output := []string{}
	inSingleQuote := false
	buffer := ""
	for _, c := range s {
		if c == '\'' {
			inSingleQuote = !inSingleQuote
		} else if unicode.IsSpace(c) && !inSingleQuote {
			if buffer != "" {
				output = append(output, buffer)
				buffer = ""
			}
		} else {
			buffer += string(c)
		}
	}
	
	if buffer != "" {
		output = append(output, buffer)
	}
		
	return output
}