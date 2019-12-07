package huh

import (
	"context"
)

var DestroyCallback *Callback

func init() {
	DestroyCallback = DefaultCallback.Destroy()
	DestroyCallback.Processor.Register(BeforeDestroyHandler)
	DestroyCallback.Processor.Register(DestroyHandler)
	DestroyCallback.Processor.Register(AfterDestroyHandler)
}

func BeforeDestroyHandler(ctx context.Context, o *Orm) (*Orm, error) {
	err := o.CallMethod("BeforeDestroy")
	return o, err
}

func DestroyHandler(ctx context.Context, o *Orm) (*Orm, error) {
	o, err := o.ParseStatement()
	if err != nil {
		return o, err
	}

	if !o.do {
		return o, nil
	}

	err = o.Exec(o.String())
	if err != nil {
		return o, ErrInvalidSQL
	}
	return o, nil
}

func AfterDestroyHandler(ctx context.Context, o *Orm) (*Orm, error) {
	err := o.CallMethod("AfterDestroy")
	return o, err
}

type DestroyCallbackProcessor struct {
	CommonCallbackProcessor
}
