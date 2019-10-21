package huh

import (
	"reflect"
	"strings"
)

var defaultPrimaryKey = "id"

type Model struct {
	TableName    string
	Fields       []*Field
	Value        reflect.Value
	PrimaryField *Field
}

func GetModel(in interface{}) *Model {
	var name string
	var fields []*Field
	var reflectValue reflect.Value
	var tagMap map[string]string
	var primaryField *Field
	var isPrimaryKey bool

	reflectType := reflect.TypeOf(in)
	name, err := getTableName(in)
	checkError(err)

	if reflectType.Kind() == reflect.Ptr {
		reflectValue = reflect.ValueOf(in).Elem()
	} else {
		reflectValue = reflect.ValueOf(in)
	}

	// get fields
	for i := 0; i < reflectValue.NumField(); i++ {
		isPrimaryKey = false
		tagMap = parseTagMap(reflectValue.Type().Field(i).Tag)
		// check isPrimaryKey
		for k, v := range tagMap {
			if k == "pk" || k == "primary_key" {
				isPrimaryKey = true
			}
		}

		field := &Field{
			Name:         reflectValue.Type().Field(i).Name,
			Value:        reflectValue.Field(i).Interface(),
			TagMap:       tagMap,
			IsPrimaryKey: isPrimaryKey,
		}
		if field.IsPrimaryKey {
			primaryField = field
		}
		fields = append(fields, field)
	}

	// get primary key
	primaryKey, err := getPrimaryKey(in)
	checkError(err)

	return &Model{
		TableName:    name,
		Fields:       fields,
		Value:        reflect.ValueOf(in),
		PrimaryField: primaryField,
	}
}

func getPrimaryKey(in interface{}) (string, error) {

	return defaultPrimaryKey, nil
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
