package huh

import "context"

var createCallback *Callback

func init() {
	createCallback = DefaultCallback.Create()
	createCallback.processor.Register(BeginTransactionHandler)
	createCallback.processor.Register(BeforeCreateHandler)
	createCallback.processor.Register(BeforeSaveHandler)
	createCallback.processor.Register(CreateHandler)
	createCallback.processor.Register(AfterSaveHandler)
	createCallback.processor.Register(AfterCreateHandler)
	createCallback.processor.Register(CommitOrRollbackTransactionHandler)
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
