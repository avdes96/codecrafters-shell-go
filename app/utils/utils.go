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
	print := false
	for i, c := range s {
		if print {
			buffer += string(c)
			print = false
		} else if c == '\'' {
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
		} else if status == IN_SINGLE_QUOTE {
			buffer += string(c)
		} else if c == '\\' {
			if status == NOT_IN_QUOTE {
				print = true
			} else {
				if next, ok := peek(s, i); !ok || !isPrintableChar(next) {
					buffer += string(c)
				} else {
					print = true
				}
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

func peek(s string, i int) (rune, bool) {
	if i >= len(s) - 1 {
		return 0, false
	}
	return rune(s[i+1]), true
}

func isPrintableChar(c rune) bool {
	return c == '\\' || c == '\n' || c == '$' || c == '"'
}