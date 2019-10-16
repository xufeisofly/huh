package huh

import "reflect"

type Model struct {
	DefaultTableName string
	Fields           []*Field
	PrimaryField     Field
}

func GetModel(in interface{}) *Model {
	reflectType := reflect.TypeOf(in)
	if reflectType.Kind() == reflect.Ptr {
		name := reflectType.Elem().Name()
	} else {
		name := reflectType.Name()
	}

	var fields []*Field

	for i := 0; i < reflectType.NumField(); i++ {
		structField := reflectType.Field(i)
		field := &Field{
			Name:   structField.Name,
			Value:  structField,
			TagMap: parseTagMap(structField.Tag),
		}
		fields = append(fields, field)
	}

	return &Model{
		DefaultTableName: name,
		Fields:           fields,
	}
}
