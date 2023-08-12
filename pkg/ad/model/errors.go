package model

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrParseResponse       = errors.New("parse response")
	ErrNotFound            = errors.New("not found")
)
