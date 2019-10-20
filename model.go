package huh

import (
	"reflect"
	"strings"
)

type Model struct {
	TableName string
	Fields    []*Field
	Value     reflect.Value
}

func GetModel(in interface{}) *Model {
	reflectType := reflect.TypeOf(in)

	var name string
	name, err := getTableName(in)
	checkError(err)

	var fields []*Field
	var reflectValue reflect.Value
	if reflectType.Kind() == reflect.Ptr {
		reflectValue = reflect.ValueOf(in).Elem()
	} else {
		reflectValue = reflect.ValueOf(in)
	}

	for i := 0; i < reflectValue.NumField(); i++ {
		field := &Field{
			Name:   reflectValue.Type().Field(i).Name,
			Value:  reflectValue.Field(i).Interface(),
			TagMap: parseTagMap(reflectValue.Type().Field(i).Tag),
		}
		fields = append(fields, field)
	}

	return &Model{
		TableName: name,
		Fields:    fields,
		Value:     reflect.ValueOf(in),
	}
}

func getTableName(in interface{}) (string, error) {
	var tableName string

	reflectValue := reflect.ValueOf(in)
	if methodValue := reflectValue.MethodByName("TableName"); methodValue.IsValid() {
		switch methodValue.Interface().(type) {
		case func() string:
			result := methodValue.Call([]reflect.Value{})
			tableName = result[0].Interface().(string)
		default:
			tableName = ""
		}
	}

	reflectType := reflect.TypeOf(in)
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
		columns = append(columns, strings.ToLower(field.Name))
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
