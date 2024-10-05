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
