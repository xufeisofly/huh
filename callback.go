package huh

import "context"

type Callback struct {
	create *CallbackProcessor
}

var DefaultCallback = &Callback{}

type CallbackProcessor struct {
	ttype    string
	handlers []CallbackHandler
}

type CallbackHandler func(context.Context, *Orm) error

func (c *Callback) Create() *CallbackProcessor {
	if c.create == nil {
		c.create = &CallbackProcessor{
			ttype: "create",
		}
	}
	return c.create
}

func (cp *CallbackProcessor) Register(handler CallbackHandler) {
	cp.handlers = append(cp.handlers, handler)
}
