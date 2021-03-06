package huh

import (
	"context"
	"reflect"
)

var SelectCallback *Callback

func init() {
	SelectCallback = DefaultCallback.Select()
	SelectCallback.Processor.Register(SelectHandler)
}

func SelectHandler(ctx context.Context, o *Orm) (*Orm, error) {
	var results []map[string]string

	o, err := o.ParseStatement()
	if err != nil {
		return o, err
	}

	if !o.do {
		return o, nil
	}

	rows, _ := o.Query(o.String())
	defer rows.Close()

	colNames, _ := rows.Columns()

	cols := make([]interface{}, len(colNames))
	colPtrs := make([]interface{}, len(colNames))
	for i, _ := range cols {
		colPtrs[i] = &cols[i]
	}

	for rows.Next() {
		rows.Scan(colPtrs...)
		ret := make(map[string]string)

		for i, col := range cols {
			colName := colNames[i]
			colValueStr := string(col.([]uint8))
			ret[colName] = colValueStr
		}
		results = append(results, ret)
	}

	err = o.SetSelectResult(results, o.result)
	return o, err
}

func canAssign(v reflect.Value) bool {
	return v.Kind() == reflect.Struct
}

type SelectCallbackProcessor struct {
	CommonCallbackProcessor
}
