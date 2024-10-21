package main

import (
	"fmt"
	"github.com/liuuner/go-cli-input"
)

func main() {
	var err error

	// List of items
	items := []string{
		"Item 1",
		"Item 2",
		"Item 3",
		"Item 4",
	}

	// Select input
	s := input.NewSelect("Select an Option:", items, func(s string) string {
		return s
	})

	state, err := s.Open()
	if err != nil {
		return
	}
	selectedOption := state.Resolve()

	// Text input for email
	s2 := input.NewText("Enter your Email:", input.WithDefaultText("example@mail.com"))

	state2, err := s2.Open()
	if err != nil {
		return
	}
	email := state2.Resolve()

	// Text input for password (sensitive input)
	s3 := input.NewText("Enter your Password:", input.WithIsSensitive(true))

	state3, err := s3.Open()
	if err != nil {
		return
	}
	password := state3.Resolve()

	// Checkbox input
	s4 := input.NewCheckbox("Select one or more options:", items, func(s string) string {
		return s
	})

	state4, err := s4.Open()
	if err != nil {
		return
	}
	selectedCheckboxes := state4.Resolve()

	// Boolean input
	s5 := input.NewBoolean("Are you sure?", 0)

	state5, err := s5.Open()
	if err != nil {
		return
	}
	confirmation := state5.Resolve()

	// Print all results
	fmt.Println("Selected Option:", selectedOption)
	fmt.Println("Entered Email:", email)
	fmt.Println("Entered Password:", password)
	fmt.Println("Selected Checkboxes:", selectedCheckboxes)
	fmt.Println("Confirmation:", confirmation)

}
