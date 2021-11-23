package services

import (
	"fmt"
	"github.com/sjeninfo/goconvert"
	"github.com/sjeninfo/sjmq"
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

func SendEvent(c *goconvert.Converter, sender sjmq.ISender, event interface{}, data interface{}) error {
	if err := c.Convert(data, &event); err != nil {
		return err
	}
	err := sender.SendEvent(event)
	return err
}
