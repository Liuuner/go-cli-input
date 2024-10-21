package input

import (
	"atomicgo.dev/keyboard/keys"
	"fmt"
	"github.com/liuuner/go-cli-input/colors"
	"github.com/liuuner/go-cli-input/cursor"
	"strings"
)

type CheckboxItem[T any] struct {
	value   T
	checked bool
}

type CheckboxState[T any] struct {
	items     []CheckboxItem[T]
	GetName   func(T) string
	GetColor  func(T) colors.Formatter
	cursorPos int
}

func NewCheckbox[T any](prompt string, items []T, getName func(T) string) Input[CheckboxState[T]] {
	i := newInput[CheckboxState[T]]()

	checkboxItems := make([]CheckboxItem[T], len(items))
	for i, item := range items {
		checkboxItems[i] = CheckboxItem[T]{value: item, checked: false}
	}

	state := CheckboxState[T]{
		items:     checkboxItems,
		cursorPos: 0,
		GetName:   getName,
	}

	return Input[CheckboxState[T]]{
		render:          renderCheckbox[T],
		handleInput:     handleCheckbox[T],
		close:           closeCheckbox[T],
		userPrompt:      prompt,
		inputPrompt:     "› - Use arrow-keys. Return to submit.",
		hasPrompt:       i.hasPrompt,
		hasSummary:      i.hasSummary,
		failedString:    i.failedString,
		completedString: i.completedString,
		promptString:    i.promptString,
		state:           state,
	}
}

func renderCheckbox[T any](s *CheckboxState[T], rerender bool) {
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

		menuItemText := s.GetName(item.value)
		if s.GetColor != nil {
			menuItemText = s.GetColor(item.value)(menuItemText)
		}
		checkboxString := "[ ]"
		if index == s.cursorPos { // for color or other effects
			checkboxString = fmt.Sprintf("[%s]", col.Gray("X"))
			menuItemText = col.Underline(menuItemText)
		}
		if item.checked {
			checkboxString = "[X]"
		}

		fmt.Printf("\r  %s %s%s", checkboxString, menuItemText, newline)
	}
}

func handleCheckbox[T any](s *CheckboxState[T], key keys.Key) (stop bool, err error) {
	switch key.Code {
	case keys.Up:
		s.cursorPos--
		s.keepPosInBoundaries()
	case keys.Down:
		s.cursorPos++
		s.keepPosInBoundaries()
	case keys.Left:
		//select none
		s.setAllCheckedState(false)
	case keys.Right:
		//select all
		s.setAllCheckedState(true)
	case keys.Space:
		//check/uncheck current
		s.items[s.cursorPos].checked = !s.items[s.cursorPos].checked
	case keys.Enter:
		stop, err = true, nil
	}

	return
}

func closeCheckbox[T any](s *CheckboxState[T], err error) (summary string) {
	cursor.ClearLine()
	for i := 0; i < len(s.items)-1; i++ {
		cursor.Up()
		cursor.ClearLine()
	}

	checkedItems := s.getCheckedItems()

	if len(checkedItems) == 0 {
		summary = "none"
	} else {
		names := make([]string, len(checkedItems))
		for i, item := range checkedItems {
			names[i] = s.GetName(item.value)
		}
		summary = strings.Join(names, ", ")
	}

	if err != nil {
		summary = err.Error()
	}
	cursor.Show()
	return
}

func (s *CheckboxState[T]) Resolve() []T {
	return s.toItems(s.getCheckedItems())
}

func (s *CheckboxState[T]) setAllCheckedState(checked bool) {
	for i := range s.items {
		s.items[i].checked = checked
	}
}

func (s *CheckboxState[T]) getCheckedItems() (checkedItems []CheckboxItem[T]) {
	for _, item := range s.items {
		if item.checked {
			checkedItems = append(checkedItems, item)
		}
	}
	return checkedItems
}

func (s *CheckboxState[T]) toItems(items []CheckboxItem[T]) []T {
	nonCheckboxItems := make([]T, len(items))
	for i, checkboxItem := range items {
		nonCheckboxItems[i] = checkboxItem.value
	}
	return nonCheckboxItems
}

func (s *CheckboxState[T]) keepPosInBoundaries() {
	s.cursorPos = (s.cursorPos + len(s.items)) % len(s.items)
}
