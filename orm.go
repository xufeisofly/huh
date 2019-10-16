package huh

import "context"

// Orm is the base struct
type Orm struct {
	masterDB *huhDB
	slaveDBs []*huhDB

	callbacks []Callback
	model     *Model
	// operator  Operator
	ast *AST
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
	return c
}

func (o *Orm) clone() *Orm {
	return &Orm{
		masterDB:  o.masterDB,
		slaveDBs:  o.slaveDBs,
		callbacks: o.callbacks,
		model:     o.model,
		ast:       o.ast,
	}
}

func (o *Orm) Of(ctx context.Context, in interface{}) error {
	c := o.clone()
	c.model.Model = in
}

func (o *Orm) CallMethod(methodName string) error {
	return nil
}
