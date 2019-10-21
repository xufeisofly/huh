package huh

import (
	"context"
	"reflect"
	"strings"
)

// Orm is the base struct
type Orm struct {
	masterDB *huhDB
	slaveDBs []*huhDB

	callbacks []Callback
	model     *Model
	operator  Operator
	statement SQLStatement
	newValues map[string]interface{}
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

func (o *Orm) Update(args ...interface{}) *Orm {
	mapArg := make(map[string]interface{})
	if len(args) != 1 {
		mapArg = multiArgsToMap(args...)
	} else {
		mapArg = args[0].(map[string]interface{})
	}
	return o.update(mapArg)
}

func (o *Orm) update(arg map[string]interface{}) *Orm {
	c := o.clone()
	c.operator = OperatorUpdate
	c.newValues = arg

	return c
}

func (o *Orm) Do(ctx context.Context, in interface{}) error {
	c := o.Of(ctx, in)
	err := c.callCallbacks(ctx)
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

func (o *Orm) clone() *Orm {
	return &Orm{
		masterDB:  o.masterDB,
		slaveDBs:  o.slaveDBs,
		callbacks: o.callbacks,
		model:     o.model,
		operator:  o.operator,
		statement: o.statement,
		newValues: o.newValues,
	}
}

func (o *Orm) callCallbacks(ctx context.Context) error {
	var cb *Callback
	switch o.operator {
	case OperatorCreate:
		cb = createCallback
	case OperatorUpdate:
		cb = updateCallback
	default:
		return ErrInvalidOperator
	}

	err := cb.processor.Process(ctx, o)
	return err
}

func (o *Orm) parseSQLStatement() (SQLStatement, error) {
	switch o.operator {
	case OperatorCreate:
		return InsertStatement{
			TableName: o.model.TableName,
			Columns:   o.model.Columns(),
			Values:    o.model.Values(),
		}, nil
	case OperatorUpdate:
		return UpdateStatement{
			TableName:    o.model.TableName,
			PrimaryKey:   strings.ToLower(o.model.PrimaryField.Name),
			PrimaryValue: o.model.PrimaryField.Value,
			Values:       o.newValues,
		}, nil
	default:
		return nil, ErrInvalidOperator
	}
}
