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
	o.CallMethod("BeforeCreate")
	return nil
}

type User struct {
	ID uint
}

func (u *User) BeforeCreate(ctx context.Context) error {
	// xxx
	return nil
}
