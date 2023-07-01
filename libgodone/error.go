package libgodone

import (
	"fmt"
	"os"
)

type ErrorCategory int

const (
	USAGE   ErrorCategory = iota
	RUNTIME ErrorCategory = iota
)

// Prints the error message approriately to its category and quits
// with exit status 1.
func Die(category ErrorCategory, message string) {
	var formatString string

	switch category {
	case USAGE:
		formatString = "Usage Error: %s\n"

	case RUNTIME:
		formatString = "Runtime Error: %s\n"
	}

	fmt.Fprintf(os.Stderr, formatString, message)
	os.Exit(1)
}
