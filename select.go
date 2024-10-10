package input

import (
	"atomicgo.dev/keyboard/keys"
	"fmt"
	"github.com/liuuner/go-cli-input/colors"
	"github.com/liuuner/go-cli-input/cursor"
)

type SelectState[T any] struct {
	items      []T
	GetName    func(T) string
	GetColor   func(T) colors.Formatter
	cursorRune rune
	cursorPos  int
}

func NewSelect[T any](items []T, getName func(T) string) Input[SelectState[T]] {
	i := newInput[SelectState[T]]()

	state := SelectState[T]{
		items:      items,
		cursorRune: '❯',
		cursorPos:  0,
		GetName:    getName,
	}

	return Input[SelectState[T]]{
		render:          renderSelect[T],
		handleInput:     handleSelect[T],
		close:           closeSelect[T],
		userPrompt:      "Select an Option",
		inputPrompt:     "› - Use arrow-keys. Return to submit.",
		hasPrompt:       i.hasPrompt,
		hasSummary:      i.hasSummary,
		failedString:    i.failedString,
		completedString: i.completedString,
		promptString:    i.promptString,
		state:           state,
	}
}

func renderSelect[T any](s *SelectState[T], rerender bool) {
	if rerender {
		// Move cursor to top
		cursor.UpN(len(s.items) - 1)
	} else {
		cursor.Hide()
	}

	for index, item := range s.items {
		var newline = "\n"
		if index == len(s.items)-1 {
			// Adding a new line on the last option will move the cursor position out of range
			// For out redrawing
			newline = ""
		}

		menuItemText := s.GetName(item)
		if s.GetColor != nil {
			menuItemText = s.GetColor(item)(menuItemText)
		}
		cursorString := "   "
		if index == s.cursorPos { // for color or other effects
			cursorString = col.Cyan(string(s.cursorRune), "  ")
			menuItemText = col.Underline(menuItemText)
		}

		fmt.Printf("\r%s %s%s", cursorString, menuItemText, newline)
	}
}

func handleSelect[T any](s *SelectState[T], key keys.Key) (stop bool, err error) {
	switch key.Code {
	case keys.Left:
		s.cursorPos = 0
	case keys.Right:
		s.cursorPos = len(s.items) - 1
	case keys.Up:
		s.cursorPos--
		s.keepPosInBoundaries()
	case keys.Down:
		s.cursorPos++
		s.keepPosInBoundaries()
	case keys.Enter:
		stop, err = true, nil
	}

	return
}

func closeSelect[T any](s *SelectState[T], err error) (summary string) {
	cursor.ClearLine()
	for i := 0; i < len(s.items)-1; i++ {
		cursor.Up()
		cursor.ClearLine()
	}

	summary = s.GetName(s.items[s.cursorPos])

	if err != nil {
		summary = err.Error()
	}
	cursor.Show()
	return
}

func (s *SelectState[T]) Resolve() T {
	return s.items[s.cursorPos]
}

func (s *SelectState[T]) keepPosInBoundaries() {
	s.cursorPos = (s.cursorPos + len(s.items)) % len(s.items)
}
