package huh

import "context"

func BeginTransactionHandler(ctx context.Context, o *Orm) (*Orm, error) {
	// c := o.Begin()
	// TODO begin
	return o, nil
}

func CommitOrRollbackTransactionHandler(ctx context.Context, o *Orm) (*Orm, error) {
	// err := o.Commit()
	// Commit
	return o, nil
}

func BeforeSaveHandler(ctx context.Context, o *Orm) (*Orm, error) {
	err := o.CallMethod("BeforeSave")
	return o, err
}

func AfterSaveHandler(ctx context.Context, o *Orm) (*Orm, error) {
	err := o.CallMethod("AfterSave")
	return o, err
}

type CommonCallbackProcessor struct {
	Handlers []CallbackHandler
}

func (cp *CommonCallbackProcessor) Register(handler CallbackHandler) {
	cp.Handlers = append(cp.Handlers, handler)
}

func (cp *CommonCallbackProcessor) Process(ctx context.Context, o *Orm) (*Orm, error) {
	for _, handler := range cp.Handlers {
		c, err := handler(ctx, o)
		o = c
		if err != nil {
			return o, err
		}
	}
	return o, nil
}
