package huh

import "context"

type Callback struct {
	processor CallbackProcessor
}

var DefaultCallback = &Callback{}

type CallbackProcessor interface {
	Register(CallbackHandler)
	Process(context.Context, *Orm) error
}

type CallbackHandler func(context.Context, *Orm) error

func (c *Callback) Create() *Callback {
	c.processor = &CreateCallbackProcessor{}
	return c
}
