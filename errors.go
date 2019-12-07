package huh

import "errors"

var (
	ErrDialectNotSupported = errors.New("dialect not supported")
	ErrInvalidOperator     = errors.New("invalid operation")
	ErrMethodNotFound      = errors.New("method not found")
	ErrRecordNotFound      = errors.New("record not found")
	ErrInvalidSQL          = errors.New("invalid SQL")
	ErrInvalidTransaction  = errors.New("no valid transaction")
	ErrUnchangable         = errors.New("unchangable value")
	ErrUnknownFieldType    = errors.New("unknown field type")
	ErrResultUnassignable  = errors.New("record is unassignable")
	ErrNeedPtrParam        = errors.New("need a pointer param")
)
