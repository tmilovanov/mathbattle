package models

import "errors"

var (
	ErrNotFound       = errors.New("not found")
	ErrWrongUserInput = errors.New("wrong user input")
)
