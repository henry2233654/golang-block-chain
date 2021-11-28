package services

import (
	"fmt"
	"golang-block-chain/services/contexts"
	"strings"
)

func getOrderBy(columnsStr string, desc bool) (orderBy []string) {
	columns := strings.Split(columnsStr, ",")
	for _, column := range columns {
		if desc {
			orderBy = append(orderBy, fmt.Sprintf("%s %s", column, "desc"))
		} else {
			orderBy = append(orderBy, fmt.Sprintf("%s %s", column, "asc"))
		}
	}
	return
}

func StartTransaction(c contexts.IContext, fn func() error) (err error) {
	err = c.Begin()
	defer func() {
		if err == nil {
			err = c.Commit()
		} else {
			c.Rollback()
		}
	}()
	if err != nil {
		return
	}
	err = fn()
	return
}
