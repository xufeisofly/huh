package huh

import "database/sql"

type Tx struct {
	tx             *sql.Tx
	name           string
	parent         *Tx
	savePointStack SavePointStack
}

type SavePointStack struct {
	SavePoints []SavePoint
}

type SavePoint struct {
	name string
}

func (t *Tx) Exec(query string, args ...interface{}) (sql.Result, error) {
	return t.tx.Exec(query, args...)
}

func (t *Tx) Commit() error {
	return t.tx.Commit()
}

func (t *Tx) Rollback() error {
	return t.tx.Rollback()
}

func (t *Tx) AddSavePoint(name string) {
	sp := SavePoint{name: name}
	t.savePointStack.Push(sp)
}

func (t *Tx) NextSavePoint() string {
	return t.savePointStack.Pop().name
}

func (sps *SavePointStack) curSavePoint() SavePoint {
	return sps.SavePoints[len(sps.SavePoints)-1]
}

func (sps *SavePointStack) Push(sp SavePoint) {
	sps.SavePoints = append(sps.SavePoints, sp)
}

func (sps *SavePointStack) Pop() SavePoint {
	newestSP := sps.curSavePoint()
	sps.SavePoints = sps.SavePoints[:len(sps.SavePoints)-1]
	return newestSP
}
