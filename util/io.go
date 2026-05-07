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
	if err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	response = strings.TrimSpace(strings.ToLower(response))
	if response != "y" && response != "yes" {
		return errors.New("operation cancelled by user")
	}

	return nil
}
