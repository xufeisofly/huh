package huh

import "context"

var updateCallback *Callback

func init() {
	updateCallback = DefaultCallback.Update()
	updateCallback.processor.Register(BeforeUpdateHandler)
	updateCallback.processor.Register(BeforeSaveHandler)
	updateCallback.processor.Register(UpdateHandler)
	updateCallback.processor.Register(AfterSaveHandler)
	updateCallback.processor.Register(AfterUpdateHandler)
}

func BeforeUpdateHandler(ctx context.Context, o *Orm) (*Orm, error) {
	err := o.CallMethod("BeforeUpdate")
	return o, err
}

func UpdateHandler(ctx context.Context, o *Orm) (*Orm, error) {
	o.model = GetModel(o.result)
	o.parseStatement()

	if !o.do {
		return o, nil
	}

	err := o.Exec(o.String())
	return o, err
}

func AfterUpdateHandler(ctx context.Context, o *Orm) (*Orm, error) {
	err := o.CallMethod("AfterUpdate")
	return o, err
}

type UpdateCallbackProcessor struct {
	CommonCallbackProcessor
}
