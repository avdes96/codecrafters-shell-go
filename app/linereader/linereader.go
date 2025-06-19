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
	
	for {
		b, err := l.reader.ReadByte()
		if err != nil {
			return "", "", err
		}
		switch s := string(b); s {
		case "\n", "\t":
			return s, userInput, nil
		default:
			userInput += s
			fmt.Fprint(os.Stdout, s)
		}

	}
}