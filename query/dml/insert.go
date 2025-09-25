package dml

import (
	"fmt"
	"strings"

	"github.com/Cooooing/cutil/query/base"
)

type Insert struct {
	table   string
	cols    []string
	values  [][]any
	selectQ base.SelectBuilder
}

func NewInsert() base.InsertBuilder {
	return &Insert{}
}

func (i *Insert) Into(table string) base.InsertBuilder {
	i.table = table
	return i
}

func (i *Insert) Columns(cols ...string) base.InsertBuilder {
	i.cols = append(i.cols, cols...)
	return i
}

func (i *Insert) Values(vals ...any) base.InsertBuilder {
	if len(vals) != len(i.cols) {
		panic("values count must match columns count")
	}
	i.values = append(i.values, vals)
	return i
}

func (i *Insert) Select(builder base.SelectBuilder) base.InsertBuilder {
	i.selectQ = builder
	return i
}

func (i *Insert) Build() (string, []any) {
	if i.table == "" || len(i.cols) == 0 {
		panic("insert must have table and columns")
	}

	sqlParts := []string{fmt.Sprintf("INSERT INTO %s (%s)", i.table, strings.Join(i.cols, ", "))}
	var args []any

	if i.selectQ != nil {
		sqlParts = append(sqlParts, "("+i.selectQ.GetSql()+")")
		args = append(args, i.selectQ.GetArgs()...)
	} else if len(i.values) > 0 {
		var valPlaceholders []string
		for _, row := range i.values {
			if len(row) != len(i.cols) {
				panic("values count must match columns count")
			}
			placeholders := make([]string, len(row))
			for j := range row {
				placeholders[j] = "?"
			}
			valPlaceholders = append(valPlaceholders, fmt.Sprintf("(%s)", strings.Join(placeholders, ", ")))
			args = append(args, row...)
		}
		sqlParts = append(sqlParts, "VALUES "+strings.Join(valPlaceholders, ", "))
	} else {
		panic("insert must have values or select")
	}

	return strings.Join(sqlParts, " "), args
}

func (i *Insert) GetSql() string {
	sql, _ := i.Build()
	return sql
}

func (i *Insert) GetArgs() []any {
	_, args := i.Build()
	return args
}
