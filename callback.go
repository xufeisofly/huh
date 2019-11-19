package huh

import "context"

type Callback struct {
	Processor CallbackProcessor
}

var DefaultCallback = &Callback{}

var mainHandlers = map[string]bool{
	"CreateHandler":  true,
	"UpdateHandler":  true,
	"SelectHandler":  true,
	"DestroyHandler": true,
}

type CallbackProcessor interface {
	Register(CallbackHandler)
	Process(context.Context, *Orm) (*Orm, error)
	GetHandlers() []CallbackHandler
}

type CallbackHandler func(context.Context, *Orm) (*Orm, error)

func (c *Callback) Create() *Callback {
	cc := c.clone()
	cc.Processor = &CreateCallbackProcessor{}
	return cc
}

func (c *Callback) Update() *Callback {
	cc := c.clone()
	cc.Processor = &UpdateCallbackProcessor{}
	return cc
}

func (c *Callback) Destroy() *Callback {
	cc := c.clone()
	cc.Processor = &DestroyCallbackProcessor{}
	return cc
}

func (c *Callback) Select() *Callback {
	cc := c.clone()
	cc.Processor = &SelectCallbackProcessor{}
	return cc
}

func (c *Callback) clone() *Callback {
	return &Callback{
		Processor: c.Processor,
	}
}
