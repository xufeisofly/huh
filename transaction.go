package huh

import "database/sql"

type Tx struct {
	tx           *sql.Tx
	parent       *Tx
	deferedTasks []deferedTask
}

type deferedTask func(*Orm)

func (t *Tx) callDeferedTasks() error {
	return nil
}

func (t *Tx) clearDeferedTasks() error {
	return nil
}
