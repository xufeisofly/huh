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
	for _, v := range ws.Values {
		str = strings.Replace(str, "?", fmt.Sprintf("'%v'", v), 1)
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

type SelectStatement struct {
	WS              WhereStatement
	TableName       string
	SelectedColumns []string
	PrimaryKey      string
	PrimaryValue    interface{}
	Limit           uint

	// the input interface pointer need to be assigned by the query scan
	Result interface{}
}

func (ss SelectStatement) String() string {
	// SELECT * FROM `users` WHERE id = 1
	// TODO * 这里需要改成指定顺序的 columns，防止 model 字段和数据库字段顺序不一致
	rawSQL := fmt.Sprintf(
		"SELECT %s FROM `%s` WHERE %s",
		strings.Join(ss.SelectedColumns, ","),
		ss.TableName,
		ss.WS.String(),
	)
	if ss.Limit == 1 {
		rawSQL += " LIMIT 1"
	}
	return rawSQL
}
