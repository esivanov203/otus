package model

import (
	"time"
)

type Event struct {
	ID          string    `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	DateStart   time.Time `db:"date_start" json:"date_start"`
	DateEnd     time.Time `db:"date_end" json:"date_end"`
	UserID      string    `db:"user_id" json:"user_id"`
}

func (e *Event) ValidateCreate() error {
	var ve ValidationError

	if e.UserID == "" {
		ve.addMessage(emptyUserID)
	}
	if e.Title == "" {
		ve.addMessage(emptyEventTitle)
	}
	if e.DateStart.After(e.DateEnd) {
		ve.addMessage(startDateMoreEndDate)
	}

	if len(ve.Messages) > 0 {
		return ve
	}
	return nil
}

func (e *Event) ValidateUpdate() error {
	var ve ValidationError

	if e.ID == "" {
		ve.addMessage(emptyID)
	}
	if e.Title == "" {
		ve.addMessage(emptyEventTitle)
	}
	if e.DateStart.After(e.DateEnd) {
		ve.addMessage(startDateMoreEndDate)
	}

	if len(ve.Messages) > 0 {
		return ve
	}
	return nil
}

func (e *Event) ValidateOne() error {
	var ve ValidationError

	if e.ID == "" {
		ve.addMessage(emptyID)
	}
	return nil
}

func (e *Event) ValidateList() error {
	var ve ValidationError

	if e.UserID == "" {
		ve.addMessage(emptyID)
	}
	if e.DateStart.IsZero() {
		ve.addMessage(emptyDate)
	}
	return nil
}
