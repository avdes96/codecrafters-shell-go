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

const (
	NOT_IN_QUOTE = iota
	IN_SINGLE_QUOTE
	IN_DOUBLE_QUOTE
)


func ParseArgs(s string) []string {
	output := []string{}
	status := NOT_IN_QUOTE
	buffer := ""
	for _, c := range s {
		if c == '\'' {
			switch status {
			case NOT_IN_QUOTE:
				status = IN_SINGLE_QUOTE
			case IN_SINGLE_QUOTE:
				status = NOT_IN_QUOTE
			case IN_DOUBLE_QUOTE:
				buffer += string(c)
			}
		} else if c == '"' {
			switch status {
			case NOT_IN_QUOTE:
				status = IN_DOUBLE_QUOTE
			case IN_DOUBLE_QUOTE:
				status = NOT_IN_QUOTE
			case IN_SINGLE_QUOTE:
				buffer += string(c)
			}
		} else if unicode.IsSpace(c) && status == NOT_IN_QUOTE {
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