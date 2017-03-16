package repl

import (
	"fmt"

	"github.com/piot/hasty-repl/commander"
)

type Repl struct {
	line      *Line
	evaluator Evaluator
}

func NewRepl(commander *commander.Commander) (Repl, error) {
	line := NewLine()
	evaluator := NewEvaluator(commander)
	return Repl{line: line, evaluator: evaluator}, nil
}

func (in *Repl) close() {
	in.line.Close()
}

func (in *Repl) prompt() error {
	input, err := in.line.Prompt()
	if err != nil {
		return err
	}

	if input == "" {
		return nil
	}

	err = in.evaluator.Eval(input)
	if err != nil { /*
			if err == ErrContinue {
				continue
			} else if err == ErrQuit {
				break
			}*/
		fmt.Println(err)
		return err
	}
	in.line.Accepted()
	return nil
}

func (in *Repl) PromptForever() error {
	for {
		err := in.prompt()
		if err != nil {
			in.line.SaveHistory()
			in.close()
			return err
		}
	}
}
