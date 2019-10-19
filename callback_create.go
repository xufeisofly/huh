package huh

import "context"

func init() {
	DefaultCallback.Create().Register(BeginTransactionHandler)
	DefaultCallback.Create().Register(BeforeCreateHandler)
	DefaultCallback.Create().Register(BeforeSaveHandler)
	DefaultCallback.Create().Register(CreateHandler)
	DefaultCallback.Create().Register(AfterSaveHandler)
	DefaultCallback.Create().Register(AfterCreateHandler)
	DefaultCallback.Create().Register(CommitOrRollbackTransactionHandler)
}

var BeforeCreateHandler = func(ctx context.Context, o *Orm) error {
	return o.CallMethod("BeforeCreate")
}

var BeginTransactionHandler = func(ctx context.Context, o *Orm) error {
	// TODO begin
	return nil
}

var BeforeSaveHandler = func(ctx context.Context, o *Orm) error {
	return o.CallMethod("BeforeSave")
}

var CreateHandler = func(ctx context.Context, o *Orm) error {
	// Create
	return nil
}

var AfterSaveHandler = func(ctx context.Context, o *Orm) error {
	return o.CallMethod("AfterSave")
}

var AfterCreateHandler = func(ctx context.Context, o *Orm) error {
	return o.CallMethod("AfterCreate")
}

var CommitOrRollbackTransactionHandler = func(ctx context.Context, o *Orm) error {
	// Commit
	return nil
}
