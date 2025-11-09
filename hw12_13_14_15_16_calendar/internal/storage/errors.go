package storage

import "errors"

var (
	ErrDateBusy = errors.New("event time is busy by another event")
	ErrNotFound = errors.New("event not found")
)
