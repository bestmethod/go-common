package gocommon

import (
	"errors"
	"fmt"
)

// because I shouldn't need to do that
func makeError(format string, args ...interface{}) error {
	return errors.New(fmt.Sprintf(format, args))
}

// a lightly more advanced version, puts latest error in line on top and previous error as appended underneath
func appendErrorln(err error, format string, args ...interface{}) error {
	if err != nil {
		var newError string
		if len(args) != 0 {
			newError = fmt.Sprintf(format, args)
		} else {
			newError = format
		}
		return errors.New(fmt.Sprintf("%s\n%s", newError, err.Error()))
	} else {
		return nil
	}
}

// like above, but on a single line
func appendError(err error, format string, args ...interface{}) error {
	if err != nil {
		var newError string
		if len(args) != 0 {
			newError = fmt.Sprintf(format, args)
		} else {
			newError = format
		}
		return errors.New(fmt.Sprintf("%s: %s", newError, err.Error()))
	} else {
		return nil
	}
}
