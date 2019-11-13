package huh

type Operator int

const (
	OperatorSelect Operator = iota
	OperatorCreate
	OperatorUpdate
	OperatorDelete
)
