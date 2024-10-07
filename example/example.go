package main

import (
	input "github.com/liuuner/go-cli-input"
)

func main() {
	items := []string{
		"Item 1",
		"Item 2",
		"Item 3",
		"Item 4",
	}
	s := input.NewSelect[string](items, func(s string) string {
		return s
	})

	_, err := s.Open()
	if err != nil {
		return
	}

	s2 := input.NewText("myDefault")

	_, err = s2.Open()
	if err != nil {
		return
	}
}
