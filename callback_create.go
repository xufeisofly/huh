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

func BeginTransactionHandler(ctx context.Context, o *Orm) error {
	// TODO begin
	return nil
}

func BeforeSaveHandler(ctx context.Context, o *Orm) error {
	return o.CallMethod("BeforeSave")
}

func CreateHandler(ctx context.Context, o *Orm) error {
	err := o.Exec(o.String())
	return err
}

func AfterSaveHandler(ctx context.Context, o *Orm) error {
	return o.CallMethod("AfterSave")
}

func AfterCreateHandler(ctx context.Context, o *Orm) error {
	return o.CallMethod("AfterCreate")
}

func CommitOrRollbackTransactionHandler(ctx context.Context, o *Orm) error {
	// Commit
	return nil
}

type CreateCallbackProcessor struct {
	Handlers []CallbackHandler
}

func (cp *CreateCallbackProcessor) Register(handler CallbackHandler) {
	cp.Handlers = append(cp.Handlers, handler)
}

func (cp *CreateCallbackProcessor) Process(ctx context.Context, o *Orm) error {
	for _, handler := range cp.Handlers {
		err := handler(ctx, o)
		if err != nil {
			return err
		}
	}
	return nil
}
