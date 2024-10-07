package go_cli_input

import (
	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"errors"
	"fmt"
	"github.com/liuuner/go-cli-input/colors"
	"github.com/liuuner/go-cli-input/cursor"
)

type Input[T any] struct {
	render            func(s *T, rerender bool)
	handleInput       func(s *T, key keys.Key) (stop bool, err error)
	close             func(s *T, err error) (summary string)
	userPrompt        string
	inputPrompt       string
	promptString      string
	completedString   string
	failedString      string
	hasPrompt         bool
	hasSummary        bool
	isLevelWithPrompt bool // if the input is on the same height as the prompt or it it's on a newline
	state             T
}

var col = colors.CreateColors(true)

func (i *Input[T]) Open() (state T, err error) {

	if i.hasPrompt {
		if i.isLevelWithPrompt {
			fmt.Printf("%s %s: %s",
				col.Cyan(i.promptString),
				i.userPrompt,
				col.Gray(i.inputPrompt),
			)
		} else {
			fmt.Printf("%s %s: %s\n",
				col.Cyan(i.promptString),
				i.userPrompt,
				col.Gray(i.inputPrompt),
			)
		}
	}

	i.render(&i.state, false)

	err = keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		switch key.Code {
		case keys.CtrlC:
			return true, errors.New("terminated with SIGINT (130)")
		case keys.Escape:
			return true, errors.New("canceled")

		case keys.Enter:
			return true, nil
		}

		stop, err = i.handleInput(&i.state, key)
		i.render(&i.state, true)
		return
	})

	summary := i.close(&i.state, err)

	if i.hasPrompt {
		if !i.isLevelWithPrompt {
			cursor.Up()
		}
		cursor.ClearLine()
	}
	cursor.StartOfLine()

	if i.hasSummary {
		if err != nil {
			fmt.Printf("%s %s: %s\n",
				col.Red(i.failedString),
				i.userPrompt,
				col.Gray(summary),
			)
		} else {
			fmt.Printf("%s %s: %s\n",
				col.Green(i.completedString),
				i.userPrompt,
				col.Gray(summary),
			)
		}
	}

	return i.state, err
}

func newInput[T any]() *Input[T] {
	return &Input[T]{
		promptString:      "?",
		completedString:   "✔",
		failedString:      "✖",
		hasPrompt:         true,
		hasSummary:        true,
		isLevelWithPrompt: false,
	}
}
