package dml

import (
	"fmt"
	"strings"

	"github.com/Cooooing/cutil/query/base"
)

type Delete struct {
	table      string
	tableAlias string
	whereCond  base.ConditionBuilder
}

func NewDelete() *Delete {
	return &Delete{}
}

func (d *Delete) From(table string) base.DeleteBuilder {
	d.table = table
	d.tableAlias = ""
	return d
}

func (d *Delete) FromAlias(table, alias string) base.DeleteBuilder {
	d.table = table
	d.tableAlias = alias
	return d
}

func (d *Delete) Where(cond base.ConditionBuilder) base.DeleteBuilder {
	d.whereCond = cond
	return d
}

func (d *Delete) Build() (string, []any) {
	if d.table == "" {
		panic("delete must have table")
	}

	sqlParts := []string{fmt.Sprintf("DELETE FROM %s", d.table)}
	if d.tableAlias != "" {
		sqlParts[0] += " AS " + d.tableAlias
	}

	args := []any{}
	if d.whereCond != nil {
		whereSQL, whereArgs := d.whereCond.Build()
		if whereSQL != "" {
			sqlParts = append(sqlParts, "WHERE "+whereSQL)
			args = append(args, whereArgs...)
		}
	}

	return strings.Join(sqlParts, " "), args
}

func (d *Delete) GetSql() string {
	sql, _ := d.Build()
	return sql
}

func (d *Delete) GetArgs() []any {
	_, args := d.Build()
	return args
}
