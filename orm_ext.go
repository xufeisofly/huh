package huh

import (
	"fmt"
	"reflect"

	"github.com/spf13/cast"
)

// setSelectResult assign the query result map to `&in` parameter of Do(ctx, &in)
func (o *Orm) setSelectResult(results []map[string]string) error {
	// no results, return directly
	if len(results) == 0 {
		return nil
	}

	v := reflect.ValueOf(o.statement.(SelectStatement).Result)

	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("non-pointer type %v", v.Type())
	}

	v = v.Elem()

	if v.Kind() == reflect.Struct {
		err := o.setOutputResult(v, results[0])
		if err != nil {
			return err
		}
	} else if v.Kind() == reflect.Slice {
		v.Set(reflect.MakeSlice(v.Type(), len(results), len(results)))
		if !canAssign(v.Index(0)) {
			return fmt.Errorf("can't assign non-struct to slice")
		}

		for i, result := range results {
			err := o.setOutputResult(v.Index(i), result)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (o *Orm) setOutputResult(output reflect.Value, data map[string]string) error {
	for colName, col := range data {
		fName := o.model.ColToFieldNameMap[colName]
		f := output.FieldByName(fName)

		if !f.IsValid() || !f.CanSet() {
			return fmt.Errorf("result field can't be set")
		}

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
		case reflect.Float32, reflect.Float64:
			colFloat, err := cast.ToFloat64E(col)
			if err != nil {
				return err
			}
			f.SetFloat(colFloat)
		default:
			return fmt.Errorf("unknow field type %v", f.Kind())
		}
	}
	return nil
}