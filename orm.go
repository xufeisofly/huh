package huh

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cast"
)

// Orm is the base struct
type Orm struct {
	masterDB *huhDB
	slaveDBs []*huhDB
	// transaction
	tx Tx

	callbacks []Callback
	model     *Model
	operator  Operator
	must      bool
	statement SQLStatement
	newValues map[string]interface{}
	// whether implement the sql
	do bool

	scope Scope
	// store the input interface
	result interface{}
}

// New initialize a Orm struct
func New() *Orm {
	return &Orm{
		masterDB: currentDB,
		must:     false,
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

	// Where("name", "sofly") => Where("name = ?", "sofly")
	if !strings.Contains(sqlStatement, "?") {
		sqlStatement += " = ?"
	}
	// default OperatorSelect
	c.operator = OperatorSelect
	c.scope.WSs = append(
		c.scope.WSs,
		WhereStatement{Condition: sqlStatement, Values: values},
	)
	return c
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

// Count *
func (o *Orm) Count() *Orm {
	c := o.clone()
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
	c.result = in

	c, err := c.callCallbacks(ctx)
	return c, err
}

// Begin is the begin of transaction
func (o *Orm) Begin() *Orm {
	c := o.clone()

	if c.inTransaction() {
		c.tx.parent = &c.tx
		sp := SavePoint{name: c.tx.parent.name}
		c.tx.AddSavePoint(sp)
	} else {
		tx, err := c.masterDB.Begin()
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

func (o *Orm) clone() *Orm {
	return &Orm{
		masterDB:  o.masterDB,
		slaveDBs:  o.slaveDBs,
		callbacks: o.callbacks,
		model:     o.model,
		tx:        o.tx,
		operator:  o.operator,
		must:      o.must,
		statement: o.statement,
		newValues: o.newValues,
		scope:     o.scope,
		result:    o.result,
		do:        o.do,
	}
}
