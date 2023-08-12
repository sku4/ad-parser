package model

import "errors"

var (
	ErrLastPage            = errors.New("this is last page")
	ErrProfileNotMightAuth = errors.New("profile not might auth")
	ErrTooManyRequests     = errors.New("too many requests")
)
