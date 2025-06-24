package linereader

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/term"
)

// Credit to https://github.com/danilovict2/go-shell;
// Used this to understand how to get byte-by-byte input

type linereader struct {
	reader *bufio.Reader
}

func New(r *bufio.Reader) linereader {
	return linereader{reader: r}
}

func (l linereader) ReadLine(initialState string) (string, string, error) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return "", "", err
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	userInput := initialState
	onlyTab := true
	for {
		b, err := l.reader.ReadByte()
		if err != nil {
			return "", "", err
		}
		switch b {
		case '\n':
			return string(b), userInput, nil
		case '\t':
			if onlyTab {
				return "\t\t", userInput, nil
			}
			return string(b), userInput, nil
		case '\x1b':
			b, err = l.reader.ReadByte()
			if err != nil {
				return "", "", err
			}
			switch b {
				case '[':
					b, err = l.reader.ReadByte()
					if err != nil {
						return "", "", err
					}
					switch b {
					case 'A':
						return "arrowUp", "", nil
					}
			}
		default:
			onlyTab = false
			userInput += string(b)
			fmt.Fprint(os.Stdout, string(b))
		}

	}
}