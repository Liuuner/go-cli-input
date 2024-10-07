package go_cli_input

import (
	"atomicgo.dev/keyboard/keys"
	"fmt"
	"github.com/liuuner/go-cli-input/cursor"
)

type TextState struct {
	text        []rune // The text being written
	defaultText []rune
	position    int // Cursor position within the text
}

func NewText(defaultValue string) Input[TextState] {
	i := newInput[TextState]()

	state := TextState{
		defaultText: []rune(defaultValue),
		text:        []rune(defaultValue),
		position:    len(defaultValue),
		//position:    0,
	}

	return Input[TextState]{
		render:            renderText,
		handleInput:       handleText,
		close:             closeText,
		userPrompt:        "What is your name",
		hasPrompt:         i.hasPrompt,
		hasSummary:        i.hasSummary,
		failedString:      i.failedString,
		completedString:   i.completedString,
		promptString:      i.promptString,
		isLevelWithPrompt: true,
		state:             state,
	}
}

// Render the text relative to the initial position
func renderText(s *TextState, rerender bool) {
	if rerender {
		// Move back to the start relative to the initial position
		cursor.MoveHorizontally(-s.position)

		// Clear the line from the current cursor to avoid overwriting
		fmt.Printf("\033[K")

		// Render the current text
		fmt.Print(string(s.text))

		// Move cursor to the current position within the text
		cursor.MoveHorizontally(s.position - len(s.text))
	} else {
		fmt.Print(string(s.text))
		//cursor.MoveHorizontally(s.position - len(s.text))
	}
}

func handleText(s *TextState, key keys.Key) (stop bool, err error) {
	switch key.Code {
	case keys.Left:
		if s.position > 0 {
			s.position--
			cursor.MoveHorizontally(-1)
		}
	case keys.Right:
		if s.position < len(s.text) {
			s.position++
			cursor.MoveHorizontally(1)
		}
	case keys.Backspace:
		if s.position > 0 {
			s.text = append(s.text[:s.position-1], s.text[s.position:]...)
			s.position--
			cursor.MoveHorizontally(-1)
		}
	case keys.Delete:
		if s.position < len(s.text) {
			// Remove the character at the current position
			s.text = append(s.text[:s.position], s.text[s.position+1:]...)
		}
	case keys.Space:
		s.text = append(s.text[:s.position], append([]rune{' '}, s.text[s.position:]...)...)
		s.position += len(key.Runes)
		cursor.MoveHorizontally(len(key.Runes))
	case keys.RuneKey:
		// Add the rune to the text at the cursor position
		s.text = append(s.text[:s.position], append(key.Runes, s.text[s.position:]...)...)
		s.position += len(key.Runes)
		cursor.MoveHorizontally(len(key.Runes))
	}

	return false, nil
}

func closeText(s *TextState, err error) (summary string) {
	// Move back to the start relative to the initial position
	cursor.MoveHorizontally(-s.position)

	// Clear the line from the current cursor to avoid overwriting
	fmt.Printf("\033[K")

	summary = string(s.text)

	if err != nil {
		summary = err.Error()
	}
	return
}

func (s *TextState) Resolve() string {
	return string(s.text)
}
