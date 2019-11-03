package huh

import "context"

var destroyCallback *Callback

func init() {
	destroyCallback = DefaultCallback.Destroy()
	destroyCallback.processor.Register(BeforeDestroyHandler)
	destroyCallback.processor.Register(DestroyHandler)
	destroyCallback.processor.Register(AfterDestroyHandler)
}

func BeforeDestroyHandler(ctx context.Context, o *Orm) (*Orm, error) {
	err := o.CallMethod("BeforeDestroy")
	return o, err
}

func DestroyHandler(ctx context.Context, o *Orm) (*Orm, error) {
	o.model = GetModel(o.result)
	o.parseStatement()

	if !o.do {
		return o, nil
	}

	err := o.Exec(o.String())
	return o, err
}

func AfterDestroyHandler(ctx context.Context, o *Orm) (*Orm, error) {
	err := o.CallMethod("AfterDestroy")
	return o, err
}

type DestroyCallbackProcessor struct {
	CommonCallbackProcessor
}
