package huh

import "context"

func BeginTransactionHandler(ctx context.Context, o *Orm) error {
	// TODO begin
	return nil
}

func CommitOrRollbackTransactionHandler(ctx context.Context, o *Orm) error {
	// Commit
	return nil
}

func BeforeSaveHandler(ctx context.Context, o *Orm) error {
	return o.CallMethod("BeforeSave")
}

func AfterSaveHandler(ctx context.Context, o *Orm) error {
	return o.CallMethod("AfterSave")
}

type CommonCallbackProcessor struct {
	Handlers []CallbackHandler
}

func (cp *CommonCallbackProcessor) Register(handler CallbackHandler) {
	cp.Handlers = append(cp.Handlers, handler)
}

func (cp *CommonCallbackProcessor) Process(ctx context.Context, o *Orm) error {
	for _, handler := range cp.Handlers {
		err := handler(ctx, o)
		if err != nil {
			return err
		}
	}
	return nil
}
