package huh

import "errors"

var (
	ErrDialectNotSupported = errors.New("dialect not supported")
	ErrInvalidOperator     = errors.New("sql operator is invalid")
	ErrMethodNotFound      = errors.New("method not found")
)
