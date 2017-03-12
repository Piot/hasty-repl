package repl

import (
	"fmt"
	"io"
	"strings"

	"github.com/peterh/liner"
)

const (
	promptDefault  = "hasty> "
	promptContinue = "..... "
	indent         = "    "
)

type Line struct {
	*liner.State
	buffer string
	depth  int
}

func (in *Line) promptString() string {
	if in.buffer != "" {
		return promptContinue + strings.Repeat(indent, in.depth)
	}

	return promptDefault
}

func NewLine() *Line {
	rl := liner.NewLiner()
	rl.SetCtrlCAborts(true)
	return &Line{State: rl}
}

func (in *Line) Prompt() (string, error) {
	line, err := in.State.Prompt(in.promptString())
	if err == io.EOF {
		if in.buffer != "" {
			// cancel line continuation
			in.Accepted()
			fmt.Println()
			err = nil
		} else {
			fmt.Println("You pressed EOF")
		}
	} else if err == liner.ErrPromptAborted {
		err = nil
		if in.buffer != "" {
			in.Accepted()
		} else {
			return "", fmt.Errorf("You pressed ctrl-c")
		}
	} else if err == nil {
		if in.buffer != "" {
			in.buffer = in.buffer + "\n" + line
		} else {
			in.buffer = line
		}
	}

	return in.buffer, err
}

func (in *Line) Accepted() {
	in.State.AppendHistory(in.buffer)
	in.buffer = ""
}

func (in *Line) Clear() {
	in.buffer = ""
	in.depth = 0
}
