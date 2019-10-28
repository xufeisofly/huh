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

	// store col name to struct field name mapping
	ColToFieldNameMap map[string]string
}

// GetModel get model info from `in`
func GetModel(in interface{}) *Model {
	var name string
	var fields []*Field
	var reflectValue reflect.Value
	var tagMap map[string]string
	var primaryField *Field
	var isPrimaryKey bool

	colToFieldNameMap := make(map[string]string)

	reflectType := reflect.TypeOf(in)

	if reflectType.Kind() == reflect.Ptr {
		reflectValue = reflect.ValueOf(in).Elem()
	} else {
		reflectValue = reflect.ValueOf(in)
	}

	// deal with slice, etc. o.Where(...).Do(ctx, &users)
	if reflectValue.Kind() == reflect.Slice {
		// Set Slice to in to get item reflection info, but this will change `in` temperarily
		reflectValue.Set(reflect.MakeSlice(reflectValue.Type(), 1, 1))
		// get itemValue of slice reflection value
		itemValue := reflectValue.Index(0)
		// Clear reflect value in order not to change `in`
		// TODO 用 in 的 copy 就行了
		reflectValue.Set(reflect.MakeSlice(reflectValue.Type(), 0, 0))
		if itemValue.Kind() == reflect.Struct {
			itemIndirectValue := reflect.Indirect(itemValue)

			reflectValue = itemIndirectValue
			reflectType = reflectType.Elem()
		}
	}

	name, err := getTableName(reflectValue, reflectType)
	checkError(err)
	// get fields
	for i := 0; i < reflectValue.NumField(); i++ {
		isPrimaryKey = false
		tagMap = parseTagMap(reflectValue.Type().Field(i).Tag)
		fieldName := reflectValue.Type().Field(i).Name
		colName := strings.ToLower(fieldName)
		// check isPrimaryKey
		for k, v := range tagMap {
			if k == "PK" {
				isPrimaryKey = true
			}
			if k == "COL" {
				colName = v
			}
		}

		field := &Field{
			Name:         fieldName,
			Value:        reflectValue.Field(i).Interface(),
			TagMap:       tagMap,
			IsPrimaryKey: isPrimaryKey,
			ColName:      colName,
		}

		// store col to field name mapping
		colToFieldNameMap[colName] = fieldName

		if field.IsPrimaryKey {
			primaryField = field
		}
		fields = append(fields, field)
	}

	checkError(err)

	return &Model{
		TableName:         name,
		Fields:            fields,
		Value:             reflect.ValueOf(in),
		PrimaryField:      primaryField,
		ColToFieldNameMap: colToFieldNameMap,
	}
}

func getPrimaryKey(in interface{}) (string, error) {

	return defaultPrimaryKey, nil
}

func getTableName(reflectValue reflect.Value, reflectType reflect.Type) (string, error) {
	var tableName string

	// TODO 目前还不支持 *User 下面定义 TableName
	if methodValue := reflectValue.MethodByName("TableName"); methodValue.IsValid() {
		switch methodValue.Interface().(type) {
		case func() string:
			result := methodValue.Call([]reflect.Value{})
			tableName = result[0].Interface().(string)
		default:
			tableName = ""
		}
		return tableName, nil
	}

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

func (m *Model) Values() []interface{} {
	var values []interface{}
	for _, value := range m.Fields {
		values = append(values, value.Value)
	}
	return values
}
