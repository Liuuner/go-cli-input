package main

import (
	input "github.com/liuuner/go-cli-input"
)

func main() {
	var err error

	items := []string{
		"Item 1",
		"Item 2",
		"Item 3",
		"Item 4",
	}

	s := input.NewSelect(items, func(s string) string {
		return s
	})

	_, err = s.Open()
	if err != nil {
		return
	}

	s2 := input.NewText("myDefault", false)

	_, err = s2.Open()
	if err != nil {
		return
	}

	s3 := input.NewText("myDefault", true)

	_, err = s3.Open()
	if err != nil {
		return
	}

	s4 := input.NewCheckbox(items, func(s string) string {
		return s
	})

	_, err = s4.Open()
	if err != nil {
		return
	}
}
