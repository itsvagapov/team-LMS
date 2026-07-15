package repository

import "errors"

var (
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrUserIsNil = errors.New("user is nill")
)