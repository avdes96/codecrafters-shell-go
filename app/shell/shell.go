package shell

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/builtin"
	"github.com/codecrafters-io/shell-starter-go/app/executable"
	"github.com/codecrafters-io/shell-starter-go/app/linereader"
	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

type shell struct {
	builtins map[string]builtin.Builtin
	executables []string
	completionsCache []string
}

func New() shell {
	return shell{
		builtins: getBuiltins(),
		executables: getExecutables(),
		completionsCache: []string{},
	}
}

func getBuiltins() map[string]builtin.Builtin {
	builtins := map[string]builtin.Builtin{}
	builtins["exit"] = builtin.Exit{}
	builtins["echo"] = builtin.Echo{}
	builtins["type"] = builtin.Type{Builtins: &builtins}
	builtins["pwd"] = builtin.Pwd{}
	builtins["cd"] = builtin.Cd{}
	return builtins
}

func getExecutables() []string {
	path := os.Getenv("PATH")
	dirs := strings.Split(path, string(os.PathListSeparator))
	execs := []string{}
	for _, dir := range dirs {
		contents, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, obj := range contents {
			if isExec(obj) {
				execs = append(execs, obj.Name())
			}
		}
	}
	return execs
}

func isExec(obj os.DirEntry) bool {
	info, err := obj.Info()
	if err != nil {
		return false
	}
	mode := info.Mode()
	return mode.IsRegular() && mode.Perm()&0111 != 0
}


func (s shell) Run() {
	var invoker string
	var userInput string = ""
	isNewLine := true
	reader := linereader.New(bufio.NewReader(os.Stdin))
	for {
		if isNewLine {
			fmt.Fprint(os.Stdout, "$ ")
			isNewLine = false
		}
		var err error
		invoker, userInput, err = reader.ReadLine(userInput)
		if err != nil {
			log.Fatalf("Error whilst reading user input: %v", err)
		}
		switch invoker {
		case "\n":
			s.executeCommand(userInput)
			isNewLine, userInput = true, ""
			s.completionsCache = []string{}
		case "\t\t":
			if len(s.completionsCache) > 0 {
				fmt.Fprintf(os.Stdout, "\r\n%s\r\n$ %s", strings.Join(s.completionsCache, "  "), userInput)
				s.completionsCache = []string{}
			}
		case "\t":
			completions := s.getPossibleAutocompletions(userInput)
			switch len(completions) {
			case 0:
				fmt.Fprint(os.Stdout, "\x07")
			case 1:
				end := strings.TrimPrefix(completions[0], userInput)
				fmt.Fprint(os.Stdout, end + " ")
				userInput += end + " "
			default:
				commonPrefix := getCommonPrefix(completions)
				if commonPrefix == "" {
					fmt.Fprint(os.Stdout, "\x07")
				} else {
					end := strings.TrimPrefix(commonPrefix, userInput)
					if end == "" {
						fmt.Fprint(os.Stdout, "\x07")
					} else {
						fmt.Fprint(os.Stdout, end)
						userInput += end
					}
					s.completionsCache = completions
				}
			}
		}
	}
}


func (s shell) executeCommand(userInput string) {
	var config utils.ShellConfig
	fmt.Fprint(os.Stdout, "\n")
	userInputSplit := utils.ParseString(userInput)
	command := userInputSplit[0]
	args := []string{}
	if len(userInputSplit) > 1 {
		args = userInputSplit[1:]
	}
	args, config = utils.ParseArgs(args)
	if b, ok := s.builtins[command]; ok {
		b.Run(args, config)
	} else if path := utils.FindExecutablePath(command); path != "" {
		executable.RunExecutable(command, config, args)
	} else {
		fmt.Fprint(config.StdErrFile, command + ": command not found\n")
	}
}


func (s shell) getPossibleAutocompletions(userInput string) []string {
	completions := []string{}
	for builtin := range s.builtins {
		if strings.HasPrefix(builtin, userInput) {
			completions = append(completions, builtin)
		}
	}
	for _, exec := range s.executables {
		if strings.HasPrefix(exec, userInput) {
			completions = append(completions, exec)
		}
	}
	slices.Sort(completions)
	return completions
}

func getCommonPrefix(strings []string) string {
	prefix := ""
	n, i := len(strings[0]), 0
	for i < n {
		p := strings[0][i]
		for _, s := range strings[1:] {
			if s[i] != p {
				return prefix
			}
		}
		prefix += string(p)
		i += 1
	}
	return prefix
}