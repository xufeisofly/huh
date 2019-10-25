package huh

import (
	"context"
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

	if statement.Limit == 1 {
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

// setSelectResult assign the query result map to `&in` parameter of Do(ctx, &in)
func (o *Orm) setSelectResult(result map[string]string) error {
	s := reflect.ValueOf(o.statement.(SelectStatement).Result).Elem()

	if s.Kind() == reflect.Struct {
		for colName, col := range result {
			fName := o.model.ColToFieldNameMap[colName]
			f := s.FieldByName(fName)

			if f.IsValid() && f.CanSet() {
				switch f.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					colInt64, err := cast.ToInt64E(col)
					if err != nil {
						return err
					}
					f.SetInt(colInt64)
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					colUint64, err := cast.ToUint64E(col)
					if err != nil {
						return err
					}
					f.SetUint(colUint64)
				case reflect.String:
					f.SetString(col)
				case reflect.Bool:
					colBool, err := cast.ToBoolE(col)
					if err != nil {
						return err
					}
					f.SetBool(colBool)
				default:
					return fmt.Errorf("unknow field type %v", f.Kind())
				}
			}
		}

	}
	return nil
}

type SelectCallbackProcessor struct {
	CommonCallbackProcessor
}
