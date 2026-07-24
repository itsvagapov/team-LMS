package jwt

import "errors"

var (
	ErrSecretNotSet = errors.New("jwt secret is not set")
	
)