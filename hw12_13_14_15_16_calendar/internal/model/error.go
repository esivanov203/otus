package model

import (
	"fmt"
)

const (
	emptyID              = "id is required"
	emptyUserID          = "user ID is required"
	emptyEventTitle      = "title is required"
	startDateMoreEndDate = "start date must be before end date"
	emptyDate            = "tart date is required or invalid format"
)

type ValidationError struct {
	Messages []string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation errors: %v", e.Messages)
}

func (e *ValidationError) addMessage(message string) {
	e.Messages = append(e.Messages, message)
}
