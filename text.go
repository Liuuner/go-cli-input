package input

import (
	"atomicgo.dev/keyboard/keys"
	"fmt"
	"github.com/liuuner/go-cli-input/cursor"
)

type TextState struct {
	text        []rune // The text being written
	defaultText []rune
	position    int // Cursor position within the text
	isSensitive bool
}

type TextOption func(*Input[TextState])

func NewText(prompt string, opts ...TextOption) Input[TextState] {
	i := newInput[TextState]()

	ts := TextState{
		defaultText: []rune{},
		text:        []rune{},
		isSensitive: false,
		position:    0,
	}

	s := Input[TextState]{
		render:            renderText,
		handleInput:       handleText,
		close:             closeText,
		userPrompt:        prompt,
		hasPrompt:         i.hasPrompt,
		hasSummary:        i.hasSummary,
		failedString:      i.failedString,
		completedString:   i.completedString,
		promptString:      i.promptString,
		isLevelWithPrompt: true,
		state:             ts,
	}

	for _, opt := range opts {
		opt(&s)
	}

	return s
}

// Render the text relative to the initial position
func renderText(s *TextState, _ bool) {
	if len(s.text) == 0 {
		fmt.Print(col.Gray(string(s.makeSensitiveIfNecessary(s.defaultText))))
		// Move cursor back to start
		cursor.MoveHorizontally(-len(s.defaultText))
	} else {
		// Move back to the start relative to the initial position
		cursor.MoveHorizontally(-s.position)

		// Clear the line from the current cursor to avoid overwriting
		fmt.Printf("\033[K")

		// Render the current text
		fmt.Print(string(s.makeSensitiveIfNecessary(s.text)))

		// Move cursor to the current position within the text
		cursor.MoveHorizontally(s.position - len(s.text))
	}
}

func (s *TextState) makeSensitiveIfNecessary(text []rune) []rune {
	if !s.isSensitive {
		return text
	}

	sensitiveText := make([]rune, len(text))
	for i := range sensitiveText {
		sensitiveText[i] = '*'
	}
	return sensitiveText
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
		stop, err = true, nil
	}

	return
}

func closeText(s *TextState, err error) (summary string) {
	// Move back to the start relative to the initial position
	cursor.MoveHorizontally(-s.position)

	// Clear the line from the current cursor to avoid overwriting
	fmt.Printf("\033[K")

	if err != nil {
		summary = err.Error()
	} else {
		summary = string(s.makeSensitiveIfNecessary(s.text))
		if summary == "" {
			summary = string(s.makeSensitiveIfNecessary(s.defaultText))
		}
	}
	return
}

func (s *TextState) Resolve() string {
	if len(s.text) == 0 {
		return string(s.defaultText)
	}
	return string(s.text)
}

func WithDefaultText(v string) TextOption {
	return func(s *Input[TextState]) {
		s.state.defaultText = []rune(v)
	}
}

func WithText(v string) TextOption {
	return func(s *Input[TextState]) {
		s.state.text = []rune(v)
		s.state.position = len(v)
	}
}

func WithIsSensitive(b bool) TextOption {
	return func(s *Input[TextState]) {
		s.state.isSensitive = b
	}
}
