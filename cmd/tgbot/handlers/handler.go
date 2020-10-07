package handlers

import (
	"errors"
)

var ErrCommandUnavailable = errors.New("This command is unavailable for current user")

type Handler struct {
	Name        string
	Description string
}
