package huh

import (
	"reflect"
	"strings"
)

// Field ...
type Field struct {
	Name         string
	Value        reflect.StructField
	IsPrimaryKey bool
	TagMap       map[string]string
}

func parseTagMap(tags reflect.StructTag) map[string]string {
	var tagMap map[string]string
	tagStr := tags.Get("huh")

	tags = strings.Split(tagStr, ";")
	for _, tag := range tags {
		v := strings.Split(tag, ":")
		k := strings.TrimSpace(strings.ToUpper(v[0]))
		if len(v[1]) == 0 {
			tapMap[k] = k
		} else {
			tapMap[k] = v[1]
		}
	}
	return tagMap
}
