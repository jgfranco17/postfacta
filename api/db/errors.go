package db

import (
	"errors"
)

var (
	ErrNotFound = errors.New("incident not found")
	ErrConflict = errors.New("incident conflict")
)
