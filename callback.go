package huh

import "context"

type Callback struct {
	processor CallbackProcessor
}

var DefaultCallback = &Callback{}

type CallbackProcessor interface {
	Register(CallbackHandler)
	Process(context.Context, *Orm) (*Orm, error)
}

type CallbackHandler func(context.Context, *Orm) (*Orm, error)

func (c *Callback) Create() *Callback {
	cc := c.clone()
	cc.processor = &CreateCallbackProcessor{}
	return cc
}

func (c *Callback) Update() *Callback {
	cc := c.clone()
	cc.processor = &UpdateCallbackProcessor{}
	return cc
}

func (c *Callback) Select() *Callback {
	cc := c.clone()
	cc.processor = &SelectCallbackProcessor{}
	return cc
}

func (c *Callback) clone() *Callback {
	return &Callback{
		processor: c.processor,
	}
}
