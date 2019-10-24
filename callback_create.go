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

func BeforeCreateHandler(ctx context.Context, o *Orm) error {
	return o.CallMethod("BeforeCreate")
}

func CreateHandler(ctx context.Context, o *Orm) error {
	err := o.Exec(o.String())
	return err
}

func AfterCreateHandler(ctx context.Context, o *Orm) error {
	return o.CallMethod("AfterCreate")
}

type CreateCallbackProcessor struct {
	CommonCallbackProcessor
}
