package apiserver

import "errors"

var (
	errInncorrectEmailOrPassword = errors.New("Incorrect email or password")
	errNotAuthenticated          = errors.New("Not authenticated")
)
