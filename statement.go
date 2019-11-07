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
		strValues = append(strValues, fmt.Sprintf("'%v'", v))
	}

	return fmt.Sprintf(
		"INSERT INTO %s (%v) VALUES (%v)",
		is.TableName,
		strings.Join(is.Columns, ","),
		strings.Join(strValues, ","),
	)
}

type WhereStatement struct {
	ByPK      bool
	Condition string
	Values    []interface{}
}

func (ws WhereStatement) String() string {
	// column1 = 1 AND column2 = 2
	var str = ws.Condition
	var valueStr string
	for _, v := range ws.Values {
		valueStr = toQuotedStr(v)
		str = strings.Replace(str, "?", valueStr, 1)
	}
	return str
}

type UpdateStatement struct {
	WS           WhereStatement
	TableName    string
	PrimaryKey   string
	PrimaryValue interface{}
	Values       map[string]interface{}
}

func (us UpdateStatement) String() string {
	// UPDATE `users` SET column = 1, column = 2 WHERE column1 = 1 AND column2 = 2
	var columnValueStrs []string
	for k, v := range us.Values {
		columnValueStrs = append(columnValueStrs, fmt.Sprintf("%s = '%v'", k, v))
	}

	if len(us.WS.Values) != 0 { // Use where first
		return fmt.Sprintf(
			"UPDATE `%s` SET %s WHERE %s",
			us.TableName,
			strings.Join(columnValueStrs, ","),
			us.WS.String(),
		)
	} else if us.PrimaryValue != nil { // Use model primary key second
		return fmt.Sprintf(
			"UPDATE `%s` SET %s WHERE %s",
			us.TableName,
			strings.Join(columnValueStrs, ","),
			fmt.Sprintf("%s = '%v'", us.PrimaryKey, us.PrimaryValue),
		)
	} else {
		return ""
	}
}

type DeleteStatement struct {
	WS           WhereStatement
	TableName    string
	PrimaryKey   string
	PrimaryValue interface{}
}

func (ds DeleteStatement) String() string {
	// DELETE FROM `users` WHERE id = 1
	if len(ds.WS.Values) != 0 {
		return fmt.Sprintf(
			"DELETE FROM `%s` WHERE %s",
			ds.TableName,
			ds.WS.String(),
		)
	} else if ds.PrimaryValue != nil {
		return fmt.Sprintf(
			"DELETE FROM `%s` WHERE %s",
			ds.TableName,
			fmt.Sprintf("%s = '%v'", ds.PrimaryKey, ds.PrimaryValue),
		)
	} else {
		return ""
	}
}

type SelectStatement struct {
	WS              WhereStatement
	TableName       string
	SelectedColumns []string
	PrimaryKey      string
	PrimaryValue    interface{}

	Limit  uint
	Offset uint
	Order  string
}

func (ss SelectStatement) String() string {
	// SELECT * FROM `users` WHERE id = 1
	rawSQL := fmt.Sprintf(
		"SELECT %s FROM `%s` WHERE %s",
		strings.Join(ss.SelectedColumns, ","),
		ss.TableName,
		ss.WS.String(),
	)
	if ss.Order != "" {
		rawSQL += fmt.Sprintf(" ORDER BY %s", ss.Order)
	}
	if ss.Limit != 0 {
		rawSQL += fmt.Sprintf(" LIMIT %d", ss.Limit)
	}
	if ss.Offset != 0 {
		rawSQL += fmt.Sprintf(" OFFSET %d", ss.Offset)
	}
	return rawSQL
}
