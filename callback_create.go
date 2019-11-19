package huh

import (
	"context"
)

var CreateCallback *Callback

func init() {
	CreateCallback = DefaultCallback.Create()
	CreateCallback.Processor.Register(BeforeCreateHandler)
	CreateCallback.Processor.Register(BeforeSaveHandler)
	CreateCallback.Processor.Register(CreateHandler)
	CreateCallback.Processor.Register(AfterSaveHandler)
	CreateCallback.Processor.Register(AfterCreateHandler)
}

func BeforeCreateHandler(ctx context.Context, o *Orm) (*Orm, error) {
	err := o.CallMethod("BeforeCreate")
	return o, err
}

func CreateHandler(ctx context.Context, o *Orm) (*Orm, error) {
	o.model = GetModel(o.result)
	o.parseStatement()

	if !o.do {
		return o, nil
	}

	err := o.Exec(o.String())
	return o, err
}

func AfterCreateHandler(ctx context.Context, o *Orm) (*Orm, error) {
	err := o.CallMethod("AfterCreate")
	return o, err
}

type CreateCallbackProcessor struct {
	CommonCallbackProcessor
}
