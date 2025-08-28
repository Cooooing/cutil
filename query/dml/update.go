package dml

import (
	"fmt"
	"github.com/Cooooing/cutil/query/base"
	"strings"
)

type Update struct {
	table      string
	tableAlias string
	setCols    []string
	setArgs    []any
	whereCond  base.ConditionBuilder
}

func NewUpdate() base.UpdateBuilder {
	return &Update{}
}

func (u *Update) Table(table string) base.UpdateBuilder {
	u.table = table
	u.tableAlias = ""
	return u
}

func (u *Update) TableAlias(table, alias string) base.UpdateBuilder {
	u.table = table
	u.tableAlias = alias
	return u
}

func (u *Update) Set(column string, value any) base.UpdateBuilder {
	u.setCols = append(u.setCols, fmt.Sprintf("%s = ?", column))
	u.setArgs = append(u.setArgs, value)
	return u
}

func (u *Update) Where(cond base.ConditionBuilder) base.UpdateBuilder {
	u.whereCond = cond
	return u
}

func (u *Update) Build() (string, []any) {
	if u.table == "" || len(u.setCols) == 0 {
		panic("update must have table and set columns")
	}

	sqlParts := []string{fmt.Sprintf("UPDATE %s", u.table)}
	if u.tableAlias != "" {
		sqlParts[0] += " AS " + u.tableAlias
	}

	sqlParts = append(sqlParts, "SET "+strings.Join(u.setCols, ", "))

	args := make([]any, len(u.setArgs))
	copy(args, u.setArgs)

	if u.whereCond != nil {
		whereSQL, whereArgs := u.whereCond.Build()
		if whereSQL != "" {
			sqlParts = append(sqlParts, "WHERE "+whereSQL)
			args = append(args, whereArgs...)
		}
	}

	return strings.Join(sqlParts, " "), args
}

func (u *Update) GetSql() string {
	sql, _ := u.Build()
	return sql
}

func (u *Update) GetArgs() []any {
	_, args := u.Build()
	return args
}
