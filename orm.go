package huh

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/spf13/cast"
)

// Orm is the base struct
type Orm struct {
	pool Pool
	// transaction
	tx Tx

	callbacks     []Callback
	withCallbacks bool
	model         *Model
	operator      Operator
	must          bool
	statement     SQLStatement
	newValues     map[string]interface{}
	// whether implement the sql
	do bool

	scope Scope
	// store the input interface
	result interface{}
}

// New initialize a Orm struct
func New() *Orm {
	return &Orm{
		pool:          pool,
		must:          false,
		withCallbacks: false,
	}
}

// MasterDB get
func (o *Orm) MasterDB() *huhDB {
	return o.pool.masterDB
}

// SlaveDB get by roundbin
func (o *Orm) SlaveDB() *huhDB {
	o.pool.incrSlaveIdx()
	if len(o.pool.slaveDBs) != 0 {
		return o.pool.slaveDBs[pool.slaveIdx]
	}
	return nil
}

// ExecutorDB select master or slave DB by SQL type
func (o *Orm) ExecutorDB() *huhDB {
	if o.isSelect() && o.SlaveDB() != nil {
		return o.SlaveDB()
	}
	return o.MasterDB()
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

func (o *Orm) MustUpdate(args ...interface{}) *Orm {
	c := o.clone()
	c.must = true
	c.Update(args)
	return c
}

func (o *Orm) Destroy() *Orm {
	c := o.clone()
	c.operator = OperatorDelete
	return c
}

func (o *Orm) MustDestroy() *Orm {
	c := o.clone()
	c.operator = OperatorDelete
	c.must = true
	return c
}

// Get by pk, will panic if result if none
func (o *Orm) Get(pk interface{}) *Orm {
	c := o.clone()
	c.operator = OperatorSelect

	c.scope.WSs = append(
		c.scope.WSs,
		WhereStatement{Values: []interface{}{pk}, ByPK: true},
	)
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

	c.scope.WSs = append(c.scope.WSs, WhereStatement{
		Condition: strings.Join(conditionArr, " AND "),
		Values:    values,
		ByPK:      false,
	})
	c.scope.Limit = 1
	return c
}

// Where get multiple instances by raw sql
func (o *Orm) Where(sqlStatement string, values ...interface{}) *Orm {
	c := o.clone()
	c.where(sqlStatement, false, values...)
	return c
}

func (o *Orm) where(sqlStatement string, isOr bool, values ...interface{}) *Orm {
	sqlStatement, values = o.parseWhereArgs(sqlStatement, values)

	o.operator = OperatorSelect
	o.scope.WSs = append(
		o.scope.WSs,
		WhereStatement{Condition: sqlStatement, Values: values, isOr: isOr},
	)
	return o
}

// And is AND relation between where statement
// alias of Where func when a Where func has already been used
func (o *Orm) And(sqlStatement string, values ...interface{}) *Orm {
	c := o.clone()
	if len(c.scope.WSs) == 0 {
		panic("First where statement not found")
	}
	c.where(sqlStatement, false, values...)
	return c
}

// Or is OR relation between where statement
// alias of Where func when a Where func has already been used
func (o *Orm) Or(sqlStatement string, values ...interface{}) *Orm {
	c := o.clone()
	if len(c.scope.WSs) == 0 {
		panic("First where statement not found")
	}
	c.where(sqlStatement, true, values...)
	return c
}

func (o *Orm) parseWhereArgs(sqlStatement string, values []interface{}) (string, []interface{}) {
	// Where("name", "sofly") => Where("name = ?", "sofly")
	if len(values) == 1 && !strings.Contains(sqlStatement, "?") {
		sqlStatement += " = ?"
	}

	for _, value := range values {
		reflectValue := reflect.ValueOf(value)
		if reflectValue.Kind() == reflect.Slice {

		}
	}
	return sqlStatement, values
}

func (o *Orm) First() *Orm {
	c := o.clone()
	return c.Limit(1)
}

// Offset pagination offset
func (o *Orm) Offset(i uint) *Orm {
	c := o.clone()
	c.scope.Offset = i
	return c
}

// Limit pagination limit
func (o *Orm) Limit(i uint) *Orm {
	c := o.clone()
	c.scope.Limit = i
	return c
}

// Order By
func (o *Orm) Order(str string) *Orm {
	c := o.clone()
	c.scope.Order = str
	return c
}

// Select columns
func (o *Orm) Select(cols ...string) *Orm {
	c := o.clone()
	c.scope.Cols = append(c.scope.Cols, cols...)
	return c
}

// Do is usually the end of the orm schedule, assign result to in or get data from in
func (o *Orm) Do(ctx context.Context, in interface{}) error {
	c := o.clone()
	c.do = true
	c, err := c.Of(ctx, in)

	if c.must {
		checkError(err)
	}
	return err
}

// Of parse the sql statement without calling the it
func (o *Orm) Of(ctx context.Context, in interface{}) (*Orm, error) {
	c := o.clone()
	c.SetResult(ctx, in)

	c, err := c.CallCallbacks(ctx)
	return c, err
}

func (o *Orm) SetResult(ctx context.Context, in interface{}) {
	o.result = in
}

// WithCallBacks will not invoke hooks
func (o *Orm) WithCallbacks() *Orm {
	c := o.clone()
	c.withCallbacks = true
	return c
}

// Begin is the begin of transaction
func (o *Orm) Begin() *Orm {
	c := o.clone()

	if c.inTransaction() {
		c.tx.parent = &c.tx
		sp := SavePoint{name: c.tx.parent.name}
		c.tx.AddSavePoint(sp)
	} else {
		tx, err := c.ExecutorDB().Begin()
		checkError(err)
		c.tx = Tx{
			tx:   tx,
			name: cast.ToString(time.Now().Unix()),
		}
	}

	return c
}

func (o *Orm) inTransaction() bool {
	return o.tx.tx != nil
}

func (o *Orm) inNestedTransaction() bool {
	return len(o.tx.savePointStack.SavePoints) > 0
}

// Commit the transaction
func (o *Orm) Commit() error {
	c := o.clone()
	if c.inNestedTransaction() {
		c.tx.ReleaseSP(c.tx.NextSavePoint())
		return nil
	}
	return c.tx.Commit()
}

// Rollback the transaction
func (o *Orm) Rollback() error {
	c := o.clone()
	if c.inNestedTransaction() {
		c.tx.RollbackToSP(c.tx.NextSavePoint())
		return nil
	}
	return c.tx.Rollback()
}

// Transaction call the transaction with a callback function
func (o *Orm) Transaction(ctx context.Context, f func(o *Orm)) (err error) {
	c := o.Begin()
	defer func() {
		if r := recover(); r != nil {
			_ = c.Rollback()
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("unknow panic")
			}

		} else {
			err = c.Commit()
		}
	}()

	f(c)
	return
}

// Exec the raw SQL
func (o *Orm) Exec(rawSQL string) error {
	var err error
	if o.inTransaction() {
		_, err = o.tx.Exec(rawSQL)
	} else {
		_, err = o.ExecutorDB().Exec(rawSQL)
	}
	return err
}

// QueryRow wrap of the db QueryRow
func (o *Orm) QueryRow(rawSQL string) *sql.Row {
	return o.ExecutorDB().QueryRow(rawSQL)
}

// Query wrap of the db Query
func (o *Orm) Query(rawSQL string) (*sql.Rows, error) {
	return o.ExecutorDB().Query(rawSQL)
}

func (o *Orm) String() string {
	return o.statement.String()
}

func (o *Orm) clone() *Orm {
	return &Orm{
		pool:          pool,
		callbacks:     o.callbacks,
		withCallbacks: o.withCallbacks,
		model:         o.model,
		tx:            o.tx,
		operator:      o.operator,
		must:          o.must,
		statement:     o.statement,
		newValues:     o.newValues,
		scope:         o.scope,
		result:        o.result,
		do:            o.do,
	}
}
