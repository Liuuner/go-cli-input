package cursor

import (
	"fmt"
)

func Show() {
	fmt.Printf("\u001B[?25h") // Show Cursor
}

func Hide() {
	fmt.Printf("\u001B[?25l") // Hide Cursor
}

func UpN(n int) {
	fmt.Printf("\u001B[%dA", n)
}

func Up() {
	UpN(1)
}

func ClearLine() {
	fmt.Print("\u001B[2K") // ANSI escape code to clear the line
}

func StartOfLine() {
	fmt.Print("\r")
}

// Move the cursor left or right by a certain number of columns
func MoveHorizontally(offset int) {
	if offset > 0 {
		fmt.Printf("\033[%dC", offset) // Move right
	} else if offset < 0 {
		fmt.Printf("\033[%dD", -offset) // Move left
	}
}
