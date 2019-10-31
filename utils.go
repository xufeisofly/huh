package huh

import (
	"reflect"
	"strings"
	"unicode"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func multiArgsToMap(args ...interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for i := 0; i < len(args); i += 2 {
		result[args[i].(string)] = args[i+1]
	}
	return result
}

func uint8ToString(bs []uint8) string {
	b := make([]byte, len(bs))
	for i, v := range bs {
		b[i] = byte(v)
	}
	return string(b)
}

func camelCaseToUnderline(camelStr string) string {
	var result string
	for i, c := range camelStr {
		if unicode.IsUpper(c) {
			if i == 0 || unicode.IsUpper([]rune(camelStr)[i-1]) {
				result += strings.ToLower(string(c))
				continue
			} else {
				result += "_"
			}
			result += strings.ToLower(string(c))
		} else {
			result += string(c)
		}
	}
	return result
}

// cloneInterface deep clone an interface
func cloneInterface(in interface{}) interface{} {
	reflectType := reflect.TypeOf(in)

	if reflectType.Kind() == reflect.Ptr {
		nInValue := reflect.New(reflectType.Elem())
		reflectValue := reflect.ValueOf(in).Elem()

		switch reflectValue.Kind() {
		case reflect.Slice:
			for i := 0; i < reflectValue.Len(); i++ {
				nElem := nInValue.Elem().Index(i)
				nElem.Set(reflectValue.Index(i))
			}
		case reflect.Struct:
			for i := 0; i < reflectValue.NumField(); i++ {
				nField := nInValue.Elem().Field(i)
				nField.Set(reflectValue.Field(i))
			}
		default:
			nInValue.Elem().Set(reflectValue)
		}

		return nInValue.Interface()
	}

	inCopy := in
	return inCopy
}
