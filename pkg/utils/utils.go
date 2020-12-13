package utils

import (
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	// Stdout points to the output buffer to send screen output
	Stdout io.Writer = os.Stdout
	// Stderr points to the output buffer to send errors to the screen
	Stderr io.Writer = os.Stderr
)

// Printf is just like fmt.Printf except that it send the output to Stdout. It
// is equal to fmt.Fprintf(util.Stdout, format, args)
func Printf(format string, args ...interface{}) {
	fmt.Fprintf(Stdout, format, args...)
}

// Eprintf prints the errors to the output buffer Stderr. It is equal to
// fmt.Fprintf(util.Stderr, format, args)
func Eprintf(format string, args ...interface{}) {
	fmt.Fprintf(Stderr, format, args...)
}

// Info prints to screen
func Info(format string, args ...interface{}) {
	Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// CheckIfError handles errors and exits if necessary
func CheckIfError(err error) {
	if err == nil {
		return
	}

	Eprintf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

// ListContainsSubString returns true when string s is found in the list
func ListContainsSubString(list []string, s string) bool {
	for _, value := range list {
		if strings.Contains(value, s) {
			return true
		}
	}
	return false
}

// ListContains returns true when string s is found in the list
func ListContains(list []string, s string) bool {
	for _, value := range list {
		if value == s {
			return true
		}
	}
	return false
}
