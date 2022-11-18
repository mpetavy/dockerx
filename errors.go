package main

import "fmt"

type ErrUndefinedAction struct {
}

func (e *ErrUndefinedAction) Error() string {
	return fmt.Sprintf("undefined action")
}
