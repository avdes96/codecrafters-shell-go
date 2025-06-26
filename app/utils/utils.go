package utils

import (
	"log"
	"os"
	"os/exec"
	"unicode"
)

type ShellCommand struct  {
	Command string
	Args []string
	StdInFile *os.File
	StdOutFile *os.File
	StdErrFile *os.File
}

func NewShellCommand() ShellCommand {
	return ShellCommand{
		Command: "",
		Args: []string{},
		StdInFile: os.Stdin,
		StdOutFile: os.Stdout,
		StdErrFile: os.Stderr,
	}
}

func (s *ShellCommand) addArg(arg string) {
	s.Args = append(s.Args, arg)
}

func (s *ShellCommand) Close() {
	if s.StdOutFile != os.Stdout {
		s.StdOutFile.Close()
	}
	if s.StdErrFile != os.Stderr {
		s.StdErrFile.Close()
	}
	if s.StdInFile != os.Stdin {
		s.StdInFile.Close()
	}
}

func FindExecutablePath(cmd string) string {
	path, err := exec.LookPath(cmd)
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


func ParseString(s string) []string {
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

func ParseInput(input []string) ([]ShellCommand) {
	cmds := []ShellCommand{}
	currentCmd := NewShellCommand()
	i := 0
	for i < len(input) {
		if currentCmd.Command == "" {
			currentCmd.Command = input[i]
			i += 1
		} else {
			arg := input[i]
			switch arg {
			case ">", "1>":
				if i < len(input) - 1 {
					currentCmd.StdOutFile = GetOSFile(input[i+1], true)
				}
				i += 2
			case ">>", "1>>":
				if i < len(input) - 1 {
					currentCmd.StdOutFile = GetOSFile(input[i+1], false)
				}
				i += 2
			case "2>":
				if i < len(input) - 1 {
					currentCmd.StdErrFile = GetOSFile(input[i+1], true)
				}
				i += 2
			case "2>>":
				if i < len(input) - 1 {
					currentCmd.StdErrFile = GetOSFile(input[i+1], false)
				}
				i += 2
			case "|":
				r, w, err := os.Pipe()
				if err != nil {
					log.Fatalf("Error creating pipe: %s", err)
				}
				currentCmd.StdOutFile = w
				cmds = append(cmds, currentCmd)
				currentCmd = NewShellCommand()
				currentCmd.StdInFile = r
				i += 1
			default:
				currentCmd.addArg(arg)
				i += 1
			}
		}
	}
	cmds = append(cmds, currentCmd)
	return cmds
}

func GetOSFile(filename string, overwrite bool) *os.File {
	var flags int = os.O_CREATE | os.O_WRONLY
	if !overwrite {
		flags |= os.O_APPEND
	}
	file, err := os.OpenFile(filename, flags, 0644)
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	return file
}