package huh

import "context"

var UpdateCallback *Callback

func init() {
	UpdateCallback = DefaultCallback.Update()
	UpdateCallback.Processor.Register(BeforeUpdateHandler)
	UpdateCallback.Processor.Register(BeforeSaveHandler)
	UpdateCallback.Processor.Register(UpdateHandler)
	UpdateCallback.Processor.Register(AfterSaveHandler)
	UpdateCallback.Processor.Register(AfterUpdateHandler)
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
