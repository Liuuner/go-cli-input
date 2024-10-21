package input

import (
	"atomicgo.dev/keyboard/keys"
	"errors"
	"fmt"
	"github.com/liuuner/go-cli-input/cursor"
	"slices"
	"strings"
)

type BooleanState struct {
	text           []rune // The text being written
	defaultBoolean int
	acceptStrings  []string
	declineStrings []string
	position       int // Cursor position within the text
}

func NewBoolean(prompt string, defaultBoolean int) Input[BooleanState] {
	i := newInput[BooleanState]()

	state := BooleanState{
		defaultBoolean: defaultBoolean,
		text:           []rune{},
		acceptStrings:  []string{"y", "ye", "yes"},
		declineStrings: []string{"n", "no"},
		position:       0,
	}

	inputPrompt := "[y/n] "
	if defaultBoolean == 0 {
		inputPrompt = "[y/N] "
	} else if defaultBoolean == 1 {
		inputPrompt = "[Y/n] "
	}

	return Input[BooleanState]{
		render:            renderBoolean,
		handleInput:       handleBoolean,
		close:             closeBoolean,
		userPrompt:        prompt,
		inputPrompt:       inputPrompt,
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
func renderBoolean(s *BooleanState, _ bool) {
	// Move back to the start relative to the initial position
	cursor.MoveHorizontally(-s.position)

	// Clear the line from the current cursor to avoid overwriting
	fmt.Printf("\033[K")

	// Render the current text
	fmt.Print(string(s.text))

	// Move cursor to the current position within the text
	cursor.MoveHorizontally(s.position - len(s.text))
}

func handleBoolean(s *BooleanState, key keys.Key) (stop bool, err error) {
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
		// no whitespaces at start of text
		if len(s.text) == 0 {
			break
		}
		s.text = append(s.text[:s.position], append([]rune{' '}, s.text[s.position:]...)...)
		s.position += len(key.Runes)
		cursor.MoveHorizontally(len(key.Runes))
	case keys.RuneKey:
		// Add the rune to the text at the cursor position
		runes := key.Runes

		// Replace newLines with whitespaces when being pasted
		for i, r := range runes {
			if r == '\r' {
				runes[i] = ' '
			}
		}

		s.text = append(s.text[:s.position], append(key.Runes, s.text[s.position:]...)...)
		s.position += len(key.Runes)
		cursor.MoveHorizontally(len(key.Runes))
	case keys.Enter:
		_, boolErr := s.getBoolean()
		if boolErr != nil {
			stop = false
		} else {
			stop = true
		}
	}

	return
}

func closeBoolean(s *BooleanState, err error) (summary string) {
	// Move back to the start relative to the initial position
	cursor.MoveHorizontally(-s.position)

	// Clear the line from the current cursor to avoid overwriting
	fmt.Printf("\033[K")

	b, _ := s.getBoolean()
	if b {
		summary = "yes"
	} else {
		summary = "no"
	}

	if err != nil {
		summary = err.Error()
	}

	return
}

func (s *BooleanState) Resolve() bool {
	b, _ := s.getBoolean()
	return b
}

func (s *BooleanState) getBoolean() (b bool, err error) {
	text := strings.ToLower(string(s.text))

	if text == "" && s.defaultBoolean == 0 {
		b = false
	} else if text == "" && s.defaultBoolean == 1 {
		b = true
	} else if slices.Contains(s.acceptStrings, text) {
		b = true
	} else if slices.Contains(s.declineStrings, text) {
		b = false
	} else {
		err = errors.New("invalid value")
	}
	return
}
