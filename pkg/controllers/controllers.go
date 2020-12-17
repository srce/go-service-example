package controllers

import "errors"

var (
	ErrRequredValue   = errors.New("value is required")
	ErrNotImplemented = errors.New("not implemented")
)
