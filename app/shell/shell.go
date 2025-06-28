package shell

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
	"sync"

	"github.com/codecrafters-io/shell-starter-go/app/builtin"
	"github.com/codecrafters-io/shell-starter-go/app/executable"
	"github.com/codecrafters-io/shell-starter-go/app/linereader"
	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

const histfile string = "HISTFILE"

type shell struct {
	builtins map[string]builtin.Builtin
	executables []string
	completionsCache []string
	history []string
	historyPtr int
}

func New() *shell {
	s:= shell{
		builtins: getBuiltins(),
		executables: getExecutables(),
		completionsCache: []string{},
		history: []string{},
		historyPtr: -1,
	}
	s.builtins["history"] = builtin.NewHistory(&s.history)
	if h, ok := s.builtins["history"].(*builtin.History); ok {
		h.ReadFromFile(os.Getenv(histfile), true)
	}
	return &s
}

func getBuiltins() map[string]builtin.Builtin {
	builtins := map[string]builtin.Builtin{}
	builtins["exit"] = builtin.NewExit()
	builtins["echo"] = builtin.NewEcho()
	builtins["type"] = builtin.NewType(&builtins)
	builtins["pwd"] = builtin.NewPwd()
	builtins["cd"] = builtin.NewCd()
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


func (s *shell) Run() {
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
			s.history = append(s.history, userInput)
			s.executeCommand(userInput)
			isNewLine, userInput = true, ""
			s.reset()
		case "\t\t":
			if len(s.completionsCache) > 0 {
				fmt.Fprintf(os.Stdout, "\r\n%s\r\n$ %s", strings.Join(s.completionsCache, "  "), userInput)
				s.reset()
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
		case "arrowUp":
			if len(s.history) == 0 {
				continue
			}
			s.clearLine()
			if 0 < s.historyPtr && s.historyPtr <= len(s.history) {
				s.historyPtr -= 1
				userInput = s.history[s.historyPtr]
				fmt.Fprint(os.Stdout, userInput)
			}
		case "arrowDown":
			if len(s.history) == 0 {
				continue
			}
			s.clearLine()
			if 0 <=s.historyPtr && s.historyPtr < len(s.history) {
				s.historyPtr += 1
				userInput = s.history[s.historyPtr]
				fmt.Fprint(os.Stdout, userInput)

			}
		}
	}
}


func (s *shell) reset() {
	s.completionsCache = []string{}
	s.historyPtr = len(s.history)
}

func (s* shell) clearLine() {
	fmt.Fprint(os.Stdout, "\n\033[1A\033[K") //source: https://groups.google.com/g/golang-nuts/c/k6l_LhI8CO0/m/Nu-bWHb8BwAJ
	fmt.Fprint(os.Stdout, "$ ")
}

func (s *shell) executeCommand(userInput string) {
	fmt.Fprint(os.Stdout, "\n")
	userInputSplit := utils.ParseString(userInput)
	cmds := utils.ParseInput(userInputSplit)
	var wg sync.WaitGroup
	for _, cmd := range cmds {
		wg.Add(1)
		go func(cmd utils.ShellCommand) {
			defer cmd.Close()
			defer wg.Done()
			if b, ok := s.builtins[cmd.Command]; ok {
				exitCode := b.Run(&cmd)
				s.dealWithExitCode(exitCode)

			} else if path := utils.FindExecutablePath(cmd.Command); path != "" {
				exitCode := executable.Run(&cmd)
				s.dealWithExitCode(exitCode)
			} else {
				fmt.Fprintf(cmd.StdErrFile, "%s: command not found\n", cmd.Command)
			}
		}(cmd)
	}
	wg.Wait()
}

func (s *shell) dealWithExitCode(e int) {
	if e < 0 {
		return
	}
	if h, ok := s.builtins["history"].(*builtin.History); ok {
		if f := os.Getenv(histfile); f != "" {
			h.AppendToFile(f)
		}
	}
	os.Exit(e)
}

func (s *shell) getPossibleAutocompletions(userInput string) []string {
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