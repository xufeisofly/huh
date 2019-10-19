package huh

import "reflect"

type Model struct {
	TableName string
	Fields    []*Field
	Value     reflect.Value
}

func GetModel(in interface{}) *Model {
	reflectValue := reflect.ValueOf(in)
	reflectType := reflectValue.Type()

	var name string
	if reflectType.Kind() == reflect.Ptr {
		name = reflectType.Elem().Name()
	} else {
		name = reflectType.Name()
	}

	var fields []*Field

	for i := 0; i < reflectValue.NumField(); i++ {
		field := &Field{
			Name:   reflectType.Field(i).Name,
			Value:  reflectValue.Field(i).Interface(),
			TagMap: parseTagMap(reflectType.Field(i).Tag),
		}
		fields = append(fields, field)
	}

	return &Model{
		TableName: name,
		Fields:    fields,
		Value:     reflectValue,
	}
}

func (m *Model) Columns() []string {
	var columns []string
	for _, field := range m.Fields {
		columns = append(columns, field.Name)
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
