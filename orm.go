package huh

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

// Orm is the base struct
type Orm struct {
	masterDB *huhDB
	slaveDBs []*huhDB
	// transaction
	tx      Tx
	txCount int

	callbacks []Callback
	model     *Model
	operator  Operator
	must      bool
	statement SQLStatement
	newValues map[string]interface{}

	scope Scope
}

// New initialize a Orm struct
func New() *Orm {
	return &Orm{
		masterDB: currentDB,
		must:     false,
		txCount:  0,
	}
}

// Close current DB connection
func (o *Orm) Close() error {
	return o.masterDB.Close()
}

// Create for a model instance creation
func (o *Orm) Create() *Orm {
	c := o.clone()
	c.operator = OperatorCreate
	return c
}

// MustCreate Create with error panic
func (o *Orm) MustCreate() *Orm {
	c := o.clone()
	c.operator = OperatorCreate
	c.must = true
	return c
}

// Update for a model instance update
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

// Get by pk, will panic if result if none
func (o *Orm) Get(pk interface{}) *Orm {
	c := o.clone()
	c.operator = OperatorSelect

	c.scope.WS = WhereStatement{Values: []interface{}{pk}, ByPK: true}
	c.scope.Limit = 1
	c.must = true
	return c
}

// GetBy column value limit 1
func (o *Orm) GetBy(args ...interface{}) *Orm {
	mapArg := make(map[string]interface{})
	if len(args) != 1 {
		mapArg = multiArgsToMap(args...)
	} else {
		mapArg = args[0].(map[string]interface{})
	}
	return o.getBy(mapArg)
}

func (o *Orm) getBy(arg map[string]interface{}) *Orm {
	c := o.clone()
	c.operator = OperatorSelect

	var conditionArr []string
	var values []interface{}
	for k, v := range arg {
		conditionArr = append(conditionArr, fmt.Sprintf("`%s` = ?", k))
		values = append(values, v)
	}

	c.scope.WS = WhereStatement{
		Condition: strings.Join(conditionArr, " AND "),
		Values:    values,
		ByPK:      false,
	}
	c.scope.Limit = 1
	return c
}

// Where get multiple instances by raw sql
func (o *Orm) Where(sqlStatement string, values ...interface{}) *Orm {
	c := o.clone()
	// default OperatorSelect
	c.operator = OperatorSelect
	c.scope.WS = WhereStatement{Condition: sqlStatement, Values: values}
	return c
}

// Do is usually the end of the orm schedule, assign result to in or get data from in
func (o *Orm) Do(ctx context.Context, in interface{}) error {
	c := o.Of(ctx, in)
	err := c.callCallbacks(ctx)

	if c.must {
		checkError(err)
	}
	return err
}

// Of parse the sql statement without calling the it
func (o *Orm) Of(ctx context.Context, in interface{}) *Orm {
	c := o.clone()
	c.model = GetModel(in)

	statement, err := c.parseStatement(in)
	checkError(err)
	c.statement = statement

	return c
}

// Begin is the begin of transaction
func (o *Orm) Begin() *Orm {
	c := o.clone()

	// if already in transaction, just increment txCount
	if c.inTransaction() {
		c.txCount++
		c.tx.parent = &c.tx
	} else {
		// flatify embedded transaction, add function to parentTx deferedTasks
		// deferedTasks will be executed when the last commit called
		// the deferedTask will be cleared when the its rollback called
		tx, err := c.masterDB.Begin()
		checkError(err)
		c.tx.tx = tx
	}

	return c
}

func (o *Orm) inTransaction() bool {
	return o.txCount > 0
}

// Commit the transaction
func (o *Orm) Commit() error {
	if o.inTransaction() {
		o.txCount--
	}
	return o.tx.tx.Commit()
}

// Rollback the transaction
func (o *Orm) Rollback() error {
	if o.inTransaction() {
		o.txCount--
	}
	return o.tx.tx.Rollback()
}

// Transaction call the transaction with a callback function
func (o *Orm) Transaction(ctx context.Context, f func(o *Orm)) (err error) {
	c := o.Begin()
	defer func() {
		if err := recover(); err != nil {
			_ = c.Rollback()
		}
		err = c.Commit()
	}()

	f(c)
	return
}

// Exec the raw SQL
func (o *Orm) Exec(rawSQL string) error {
	var err error
	if o.tx.tx != nil {
		_, err = o.tx.tx.Exec(rawSQL)
	} else {
		_, err = o.masterDB.Exec(rawSQL)
	}
	return err
}

// QueryRow wrap of the db QueryRow
func (o *Orm) QueryRow(rawSQL string) *sql.Row {
	return o.masterDB.QueryRow(rawSQL)
}

// Query wrap of the db Query
func (o *Orm) Query(rawSQL string) (*sql.Rows, error) {
	return o.masterDB.Query(rawSQL)
}

func (o *Orm) String() string {
	return o.statement.String()
}

// CallMethod call a model method by methodName
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
		tx:        o.tx,
		txCount:   o.txCount,
		operator:  o.operator,
		must:      o.must,
		statement: o.statement,
		newValues: o.newValues,
		scope:     o.scope,
	}
}

func (o *Orm) callCallbacks(ctx context.Context) error {
	var cb *Callback
	switch o.operator {
	case OperatorCreate:
		cb = createCallback
	case OperatorUpdate:
		cb = updateCallback
	case OperatorSelect:
		cb = selectCallback
	default:
		return ErrInvalidOperator
	}

	err := cb.processor.Process(ctx, o)
	return err
}

func (o *Orm) parseStatement(in interface{}) (SQLStatement, error) {
	switch o.operator {
	case OperatorCreate:
		return InsertStatement{
			TableName: o.model.TableName,
			Columns:   o.model.Columns(),
			Values:    o.model.Values(),
		}, nil
	case OperatorUpdate:
		return UpdateStatement{
			WS:           o.scope.WS,
			TableName:    o.model.TableName,
			PrimaryKey:   o.model.PrimaryField.ColName,
			PrimaryValue: o.model.PrimaryField.Value,
			Values:       o.newValues,
		}, nil
	case OperatorSelect:
		primaryKey := o.model.PrimaryField.ColName

		if o.scope.WS.ByPK {
			o.scope.WS.Condition = fmt.Sprintf("%s = ?", primaryKey)
		}
		return SelectStatement{
			WS:              o.scope.WS,
			Limit:           o.scope.Limit,
			TableName:       o.model.TableName,
			SelectedColumns: o.model.Columns(),
			PrimaryKey:      primaryKey,
			PrimaryValue:    o.model.PrimaryField.Value,
			Result:          in,
		}, nil
	default:
		return nil, ErrInvalidOperator
	}
}
