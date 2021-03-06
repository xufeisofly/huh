package huh

import (
	"fmt"
	"reflect"
)

var defaultPrimaryKey = "id"

type Model struct {
	TableName    string
	Fields       []*Field
	Value        reflect.Value
	PrimaryField *Field

	// store col name to struct field name mapping
	ColToFieldNameMap map[string]string
}

// GetModel get model info from `in`
func GetModel(in interface{}) (*Model, error) {
	var name string
	var fields []*Field
	var reflectValue reflect.Value
	var tagMap map[string]string
	var primaryField *Field
	var isPrimaryKey, isReadOnly bool

	// clone an `in`
	inC := cloneInterface(in)

	colToFieldNameMap := make(map[string]string)

	reflectType := reflect.TypeOf(inC)

	if reflectType.Kind() == reflect.Ptr {
		reflectValue = reflect.ValueOf(inC).Elem()
	} else {
		reflectValue = reflect.ValueOf(inC)
	}

	// deal with slice, etc. o.Where(...).Do(ctx, &users)
	if reflectValue.Kind() == reflect.Slice {
		if !reflectValue.CanSet() {
			return nil, fmt.Errorf("%w", ErrResultUnassignable)
		}
		reflectValue.Set(reflect.MakeSlice(reflectValue.Type(), 1, 1))
		itemValue := reflectValue.Index(0)

		if itemValue.Kind() == reflect.Struct {
			itemIndirectValue := reflect.Indirect(itemValue)

			reflectValue = itemIndirectValue
			reflectType = reflectType.Elem()
		}
	}

	name, err := getTableName(reflectValue, reflectType)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	// get fields
	for i := 0; i < reflectValue.NumField(); i++ {
		isPrimaryKey = false
		tagMap = parseTagMap(reflectValue.Type().Field(i).Tag)
		fieldName := reflectValue.Type().Field(i).Name
		colName := camelCaseToUnderline(fieldName)
		// check isPrimaryKey
		for k, v := range tagMap {
			if k == "PK" {
				isPrimaryKey = true
			}
			if k == "COL" {
				colName = v
			}
			if k == "READONLY" {
				isReadOnly = true
			}
		}

		field := &Field{
			Name:         fieldName,
			Value:        reflectValue.Field(i).Interface(),
			TagMap:       tagMap,
			IsPrimaryKey: isPrimaryKey,
			ColName:      colName,
			IsReadOnly:   isReadOnly,
		}

		// store col to field name mapping
		colToFieldNameMap[colName] = fieldName

		if field.IsPrimaryKey {
			primaryField = field
		}
		fields = append(fields, field)
	}

	return &Model{
		TableName:         name,
		Fields:            fields,
		Value:             reflect.ValueOf(in),
		PrimaryField:      primaryField,
		ColToFieldNameMap: colToFieldNameMap,
	}, nil
}

func getPrimaryKey(in interface{}) (string, error) {
	return defaultPrimaryKey, nil
}

func getTableName(reflectValue reflect.Value, reflectType reflect.Type) (string, error) {
	var tableName string

	methodValue := reflectValue.MethodByName("TableName")

	// if no method under in, then check if it is under in's pointer
	if !methodValue.IsValid() {

		inPtr := reflect.New(reflectValue.Type()).Interface()
		reflectValue = reflect.ValueOf(inPtr)
		methodValue = reflectValue.MethodByName("TableName")
	}

	if methodValue.IsValid() {
		switch methodValue.Interface().(type) {
		case func() string:
			result := methodValue.Call([]reflect.Value{})
			tableName = result[0].Interface().(string)
		default:
			tableName = ""
		}
		return tableName, nil
	}

	// otherwise take default name from in
	if tableName == "" {
		if reflectType.Kind() == reflect.Ptr {
			tableName = reflectType.Elem().Name()
		} else {
			tableName = reflectType.Name()
		}
	}
	return tableName, nil
}

func (m *Model) Columns() []string {
	var columns []string
	for _, field := range m.Fields {
		columns = append(columns, field.ColName)
	}
	return columns
}

func (m *Model) WritableColumns() []string {
	var columns []string
	for _, field := range m.Fields {
		if field.IsReadOnly {
			continue
		}
		columns = append(columns, field.ColName)
	}
	return columns
}

func (m *Model) Values() []interface{} {
	var values []interface{}
	for _, value := range m.Fields {
		values = append(values, value.Value)
	}
	return values
}

func (m *Model) WritableValues() []interface{} {
	var values []interface{}
	for _, value := range m.Fields {
		if value.IsReadOnly {
			continue
		}
		values = append(values, value.Value)
	}
	return values
}
