package storage

import "errors"

var (
	ErrAlreadyExist = errors.New("user already exist")
	ErrNotFound     = errors.New("user not found")
)
