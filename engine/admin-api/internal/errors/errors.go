package errors

import (
	"fmt"
)

// Wrapper creates a function that returns errors starts with a given message.
//
//nolint:goerr113 // this is a wrapper for errors
func Wrapper(message string) func(params ...interface{}) error {
	return func(params ...interface{}) error {
		return fmt.Errorf(message, params...)
	}
}
