package dao

import (
	"fmt"
	"strings"
)

type Query struct {
	Sort    string
	SortCol string
	Limit   int
	Offset  int
}

func checkWhereClause(stmt string) string {
	if !strings.HasSuffix(stmt, "WHERE") {
		return fmt.Sprintf("%s AND", stmt)
	}

	return stmt
}

func addDefaultQuery(stmt string, q Query) string {
	if q.SortCol != "" {
		stmt = fmt.Sprintf("%s ORDER BY %s", stmt, q.SortCol)
	} else {
		stmt = fmt.Sprintf("%s ORDER BY id", stmt)
	}
	if q.Sort != "" {
		stmt = fmt.Sprintf("%s %s", stmt, q.Sort)
	} else {
		stmt = fmt.Sprintf("%s ASC", stmt)
	}
	if q.Limit > 0 {
		stmt = fmt.Sprintf("%s LIMIT %d", stmt, q.Limit)
	} else {
		stmt = fmt.Sprintf("%s LIMIT 100", stmt)
	}
	if q.Offset > 0 {
		stmt = fmt.Sprintf("%s OFFSET %d", stmt, q.Offset)
	} else {
		stmt = fmt.Sprintf("%s OFFSET 0", stmt)
	}

	return stmt
}
