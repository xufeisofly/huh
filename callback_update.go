package huh

import "context"

var updateCallback *Callback

func init() {
	updateCallback = DefaultCallback.Update()
	updateCallback.processor.Register(BeginTransactionHandler)
	updateCallback.processor.Register(BeforeUpdateHandler)
	updateCallback.processor.Register(BeforeSaveHandler)
	updateCallback.processor.Register(UpdateHandler)
	updateCallback.processor.Register(AfterSaveHandler)
	updateCallback.processor.Register(AfterUpdateHandler)
	updateCallback.processor.Register(CommitOrRollbackTransactionHandler)
}

func BeforeUpdateHandler(ctx context.Context, o *Orm) error {
	return o.CallMethod("BeforeUpdate")
}

func UpdateHandler(ctx context.Context, o *Orm) error {
	err := o.Exec(o.String())
	return err
}

func AfterUpdateHandler(ctx context.Context, o *Orm) error {
	return o.CallMethod("AfterUpdate")
}

type UpdateCallbackProcessor struct {
	CommonCallbackProcessor
	// Handlers []CallbackHandler
}

// func (cp *UpdateCallbackProcessor) Register(handler CallbackHandler) {
// 	cp.Handlers = append(cp.Handlers, handler)
// }

// func (cp *UpdateCallbackProcessor) Process(ctx context.Context, o *Orm) error {
// 	for _, handler := range cp.Handlers {
// 		err := handler(ctx, o)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
