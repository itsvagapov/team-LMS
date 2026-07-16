package service

import "errors"

var (
	ErrUserNotFound           = errors.New("user not found")
	ErrEmailAlreadyExists     = errors.New("email already exists")
	ErrInvalidPassword        = errors.New("invalid password")
	ErrNameRequired           = errors.New("name is required")
	ErrInvalidCharacters      = errors.New("invalid special characters")
	ErrInvalidNameLength      = errors.New("name must be longer than 2 characters and shorter than 30")
	ErrNameContainDigits      = errors.New("name cannot contain digits")
	ErrInvalidEmail           = errors.New("invalid email")
	ErrPasswordContainsSpaces = errors.New("password cannot contain spaces")
	ErrNameContainsSpaces     = errors.New("name cannot contain spaces")
	ErrEmailContainsSpaces    = errors.New("email cannot contain spaces")
	ErrPasswordRequired       = errors.New("password is required")
	ErrEmailRequired          = errors.New("email is required")
	ErrInvalidLoginOrPassword = errors.New("invalid login or password")
	ErrRoleAlreadyAssigned    = errors.New("role already assigned")
)
