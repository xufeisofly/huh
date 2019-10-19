package huh

import (
	"fmt"
	"strings"
)

type SQLStatement interface {
	String() string
}

type InsertStatement struct {
	TableName string
	Columns   []string
	Values    []interface{}
}

func (is InsertStatement) String() string {
	var strValues []string
	for _, v := range is.Values {
		strValues = append(strValues, v.(string))
	}

	return fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		is.TableName,
		strings.Join(is.Columns, ","),
		strings.Join(strValues, ","),
	)
}

type UpdateStatement struct {
	TableName  string
	Columns    []string
	Conditions []string
}
