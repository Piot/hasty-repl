package repl

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
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
	l := &Line{State: rl}
	l.ReadHistory()
	return l
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

func historyPath() (string, error) {
	home, homeErr := homedir.Dir()
	if homeErr != nil {
		return "", homeErr
	}
	saveDir := filepath.Join(home, ".hasty-repl/")
	os.MkdirAll(saveDir, 0777)
	savePath := filepath.Join(saveDir, "history")
	return savePath, nil
}

func (in *Line) SaveHistory() error {
	savePath, saveErr := historyPath()
	if saveErr != nil {
		return saveErr
	}
	if f, err := os.Create(savePath); err != nil {
		log.Print("Error writing history file: ", err)
		return err
	} else {
		in.State.WriteHistory(f)
		f.Close()
	}
	return nil
}

func (in *Line) ReadHistory() error {
	savePath, saveErr := historyPath()
	if saveErr != nil {
		return saveErr
	}

	if f, err := os.Open(savePath); err == nil {
		in.State.ReadHistory(f)
		f.Close()
	}

	return nil
}

func (in *Line) Clear() {
	in.buffer = ""
	in.depth = 0
}
