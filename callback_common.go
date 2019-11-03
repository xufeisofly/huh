package huh

import "context"

func BeginTransactionHandler(ctx context.Context, o *Orm) (*Orm, error) {
	o = o.Begin()
	return o, nil
}

func CommitTransactionHandler(ctx context.Context, o *Orm) (*Orm, error) {
	err := o.Commit()
	return o, err
}

func RollbackTransactionHandler(ctx context.Context, o *Orm) (*Orm, error) {
	err := o.Rollback()
	return o, err
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

// Process calls the operation callback pipeline
// every operation is actually a transaction
// when error occurs in the callback pipeline, tx will roll back
func (cp *CommonCallbackProcessor) Process(ctx context.Context, o *Orm) (*Orm, error) {
	o, _ = BeginTransactionHandler(ctx, o)

	for _, handler := range cp.Handlers {
		c, err := handler(ctx, o)
		o = c
		if err != nil {
			o, _ = RollbackTransactionHandler(ctx, o)
			return o, err
		}
	}

	o, err := CommitTransactionHandler(ctx, o)
	return o, err
}
