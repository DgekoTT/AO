package storage

import (
	"errors"
)

var (
	ErrUserExists   = errors.New("modelUser already exist")
	ErrUserNotFound = errors.New("modelUser not found")
	ErrAppNotFound  = errors.New("app not found")
)
