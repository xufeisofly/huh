package huh

import (
	"reflect"
	"strings"
)

// Field ...
type Field struct {
	Name         string
	Value        interface{}
	TagMap       map[string]string
	IsPrimaryKey bool
	ColName      string
}

func parseTagMap(tags reflect.StructTag) map[string]string {
	tagMap := make(map[string]string)
	tagStr := tags.Get("huh")
	if tagStr == "" {
		return tagMap
	}

	tagItems := strings.Split(tagStr, ";")
	for _, tag := range tagItems {
		v := strings.Split(tag, ":")
		k := strings.TrimSpace(strings.ToUpper(v[0]))

		if len(v) == 1 {
			tagMap[k] = k
		} else {
			tagMap[k] = v[1]
		}
	}
	return tagMap
}
