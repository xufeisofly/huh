package huh

type Scope struct {
	WSs    []WhereStatement
	Cols   []string
	Limit  uint
	Offset uint
	Order  string
}

const AND = " AND "
const OR = " OR "

func (s *Scope) parseWhereStatement() WhereStatement {
	type condition struct {
		conStr string
		isOr   bool
	}
	var byPK bool
	var conditions []condition
	var values []interface{}
	var joinedCon, joinedStr string

	if len(s.WSs) == 1 {
		byPK = s.WSs[0].ByPK
	} else {
		byPK = false
	}

	for _, ws := range s.WSs {
		conditions = append(conditions, condition{conStr: ws.Condition, isOr: ws.isOr})
		values = append(values, ws.Values...)
	}

	// join ws.Conditions with AND or OR
	for i, condition := range conditions {
		if i == 0 {
			joinedCon = condition.conStr
		} else {
			if condition.isOr {
				joinedStr = OR
			} else {
				joinedStr = AND
			}
			joinedCon += joinedStr + condition.conStr
		}
	}

	return WhereStatement{
		ByPK:      byPK,
		Condition: joinedCon,
		Values:    values,
	}
}
