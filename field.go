package huh

import (
	"reflect"
	"strings"
)

// Field ...
type Field struct {
	Name   string
	Value  interface{}
	TagMap map[string]string
}

func parseTagMap(tags reflect.StructTag) map[string]string {
	var tagMap map[string]string
	tagStr := tags.Get("huh")

	tagItems := strings.Split(tagStr, ";")
	for _, tag := range tagItems {
		v := strings.Split(tag, ":")
		k := strings.TrimSpace(strings.ToUpper(v[0]))

		if len(v[1]) == 0 {
			tagMap[k] = k
		} else {
			tagMap[k] = v[1]
		}
	}
	return tagMap
}
