package huh

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/spf13/cast"
)

// CallMethod call a model method by methodName
func (o *Orm) CallMethod(methodName string) error {
	ctx := context.Background()
	var argsValue []reflect.Value
	var result []reflect.Value

	reflectValue := reflect.ValueOf(o.result)

	if methodValue := reflectValue.MethodByName(methodName); methodValue.IsValid() {
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

func (o *Orm) callCallbacks(ctx context.Context) (*Orm, error) {
	var cb *Callback
	switch o.operator {
	case OperatorCreate:
		cb = createCallback
	case OperatorUpdate:
		cb = updateCallback
	case OperatorSelect:
		cb = selectCallback
	default:
		return o, ErrInvalidOperator
	}

	o, err := cb.processor.Process(ctx, o)
	return o, err
}

func (o *Orm) parseStatement() {
	var s SQLStatement
	switch o.operator {
	case OperatorCreate:
		s = InsertStatement{
			TableName: o.model.TableName,
			Columns:   o.model.WritableColumns(),
			Values:    o.model.WritableValues(),
		}
	case OperatorUpdate:
		s = UpdateStatement{
			WS:           o.scope.WS,
			TableName:    o.model.TableName,
			PrimaryKey:   o.model.PrimaryField.ColName,
			PrimaryValue: o.model.PrimaryField.Value,
			Values:       o.newValues,
		}
	case OperatorSelect:
		primaryKey := o.model.PrimaryField.ColName

		if o.scope.WS.ByPK {
			o.scope.WS.Condition = fmt.Sprintf("%s = ?", primaryKey)
		}
		s = SelectStatement{
			WS:              o.scope.WS,
			Limit:           o.scope.Limit,
			Offset:          o.scope.Offset,
			Order:           o.scope.Order,
			TableName:       o.model.TableName,
			SelectedColumns: o.model.Columns(),
			PrimaryKey:      primaryKey,
			PrimaryValue:    o.model.PrimaryField.Value,
		}
	default:
		s = nil
	}
	o.statement = s
}

// setSelectResult assign the query result map to `&in` parameter of Do(ctx, &in)
func (o *Orm) setSelectResult(results []map[string]string) error {
	// no results, return directly
	if len(results) == 0 {
		return nil
	}

	v := reflect.ValueOf(o.result)

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
		case reflect.Struct: // time.Time
			layout := "2006-01-02 15:04:05"
			colTime, err := time.Parse(layout, col)
			if err != nil {
				return err
			}
			f.Set(reflect.ValueOf(colTime))
		default:
			return fmt.Errorf("unknow field type %v", f.Kind())
		}
	}
	return nil
}
