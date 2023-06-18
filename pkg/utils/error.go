package utils

import (
	"errors"
	"fmt"
)

// To append the message to the error
func AppendMessageToError(err error, message string) error {
	if err == nil {
		return errors.New(message)
	}
	return fmt.Errorf("%w \n%s", err, message)
}

// To prepend the message to the error
func PrependMessageToError(err error, message string) error {
	if err == nil {
		return errors.New(message)
	}
	return fmt.Errorf("%s \n%w", message, err)
}
