package util

import (
	"errors"
	"fmt"
	"strings"
)

func ConfirmAction(prompt string) error {
	fmt.Printf("%s (y/N): ", prompt)
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil { // e.g. new line, no input, no tty - treat as a no
		return errors.New("operation cancelled by user")
	}

	response = strings.TrimSpace(strings.ToLower(response))
	if response != "y" && response != "yes" {
		return errors.New("operation cancelled by user")
	}

	return nil
}

func FormatHeader(title string) string {
	width := len(title) + 2
	width = max(width, 30)
	padding := (width - len(title)) / 2

	line := strings.Repeat("=", padding)
	return fmt.Sprintf("%s[ %s ]%s\n", line, title, line)
}
