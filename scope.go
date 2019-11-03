package huh

import "strings"

type Scope struct {
	WSs    []WhereStatement
	Limit  uint
	Offset uint
	Order  string
}

func (s *Scope) parseWhereStatement() WhereStatement {
	var byPK bool
	var conStrs []string
	var values []interface{}

	if len(s.WSs) == 1 {
		byPK = s.WSs[0].ByPK
	} else {
		byPK = false
	}

	for _, ws := range s.WSs {
		conStrs = append(conStrs, ws.Condition)
		values = append(values, ws.Values...)
	}

	return WhereStatement{
		ByPK:      byPK,
		Condition: strings.Join(conStrs, " AND "),
		Values:    values,
	}
}
