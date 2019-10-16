package huh

type Operator int

const (
	OperatorCreate Operator = iota
	OperatorUpdate
	OperatorDelete
	OperatorSelect
)
