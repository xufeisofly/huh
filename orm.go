package huh

import (
	"context"
	"reflect"
)

// Orm is the base struct
type Orm struct {
	masterDB *huhDB
	slaveDBs []*huhDB

	callbacks []Callback
	model     *Model
	operator  Operator
	statement SQLStatement
}

// New initialize a Orm struct
func New() *Orm {
	return &Orm{
		masterDB: currentDB,
	}
}

// Close current DB connection
func (o *Orm) Close() error {
	return o.masterDB.Close()
}

func (o *Orm) Create() *Orm {
	c := o.clone()
	c.operator = OperatorCreate
	return c
}

func (o *Orm) clone() *Orm {
	return &Orm{
		masterDB:  o.masterDB,
		slaveDBs:  o.slaveDBs,
		callbacks: o.callbacks,
		model:     o.model,
		operator:  o.operator,
		statement: o.statement,
	}
}

func (o *Orm) Do(ctx context.Context, in interface{}) error {
	c := o.Of(ctx, in)
	err := c.callCallbacks(ctx)
	// err := o.Exec(c.String())
	return err
}

func (o *Orm) callCallbacks(ctx context.Context) error {
	var cb *Callback
	switch o.operator {
	case OperatorCreate:
		cb = createCallback
	default:
		return ErrInvalidOperator
	}

	err := cb.processor.Process(ctx, o)
	return err
}

func (o *Orm) Of(ctx context.Context, in interface{}) *Orm {
	c := o.clone()
	c.model = GetModel(in)

	statement, err := c.parseSQLStatement()
	checkError(err)
	c.statement = statement

	return c
}

func (o *Orm) Exec(rawSQL string) error {
	_, err := o.masterDB.Exec(rawSQL)
	return err
}

func (o *Orm) String() string {
	return o.statement.String()
}

func (o *Orm) CallMethod(methodName string) error {
	ctx := context.Background()
	var argsValue []reflect.Value
	var result []reflect.Value

	if methodValue := o.model.Value.MethodByName(methodName); methodValue.IsValid() {
		switch methodValue.Interface().(type) {
		case func(context.Context) error: // BeforeCreate
			argsValue = []reflect.Value{reflect.ValueOf(ctx)}
			result = methodValue.Call(argsValue)

			if result[0].Interface() == interface{}(nil) {
				return nil
			}
			return result[0].Interface().(error)
		default:
			return ErrMethodNotFound
		}
	}
	return nil
}

func (o *Orm) parseSQLStatement() (SQLStatement, error) {
	switch o.operator {
	case OperatorCreate:
		return InsertStatement{
			TableName: o.model.TableName,
			Columns:   o.model.Columns(),
			Values:    o.model.Values(),
		}, nil
	default:
		return nil, ErrInvalidOperator
	}
}
