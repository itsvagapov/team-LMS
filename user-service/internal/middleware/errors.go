package middleware

import "errors"

var (
	ErrUserIDNotFound   = errors.New("user id not found in context")
	ErrUserRoleNotFound = errors.New("user role not found in context")
)