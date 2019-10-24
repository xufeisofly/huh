package huh

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/spf13/cast"
)

var selectCallback *Callback

func init() {
	selectCallback = DefaultCallback.Select()
	selectCallback.processor.Register(SelectHandler)
}

func SelectHandler(ctx context.Context, o *Orm) error {
	statement := o.statement.(SelectStatement)
	var results []map[string]string

	if statement.WS.Limit == 1 {
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

		err := o.setSelectResult(results[0])
		return err
	}
	return nil
}

func (o *Orm) setSelectResult(result map[string]string) error {
	s := reflect.ValueOf(o.statement.(SelectStatement).Result).Elem()

	if s.Kind() == reflect.Struct {
		for colName, col := range result {
			fName := o.model.ColToFieldNameMap[colName]
			f := s.FieldByName(fName)

			if f.IsValid() && f.CanSet() {
				switch f.Kind() {
				case reflect.Int:
				case reflect.Uint:
				case reflect.Uint32:
					colUint64 := cast.ToUint64(col)
					f.SetUint(colUint64)
				case reflect.String:
					f.SetString(col)
				case reflect.Bool:
				default:
					return errors.New(fmt.Sprintf("unknow field type %v", f.Kind()))
				}
			}
		}

	}
	return nil
}

type SelectCallbackProcessor struct {
	CommonCallbackProcessor
}
