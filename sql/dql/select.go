package dql

import (
	"fmt"
	"github.com/Cooooing/cutil/sql/base"
	"strings"
)

type Select struct {
	columns    []string
	table      string
	tableAlias string
	joins      []joinNode
	whereCond  base.ConditionBuilder
	groupBy    []string
	havingCond base.ConditionBuilder
	orderBy    []string
	limit      int
	offset     int
}

func NewSelect() *Select {
	return &Select{
		limit:  -1,
		offset: -1,
	}
}

func NewExistSelect() *Select {
	return &Select{
		columns: []string{"1"},
		limit:   -1,
		offset:  -1,
	}
}

func (s *Select) GetSql() string {
	sql, _ := s.Build()
	return sql
}

func (s *Select) GetArgs() []any {
	_, args := s.Build()
	return args
}

func (s *Select) Build() (string, []any) {
	if len(s.columns) == 0 {
		s.columns = []string{"*"}
	}
	var sqlParts []string
	var args []any

	// Columns
	if len(s.columns) == 0 {
		sqlParts = append(sqlParts, "SELECT *")
	} else {
		sqlParts = append(sqlParts, "SELECT "+strings.Join(s.columns, ", "))
	}

	// FROM
	if s.table != "" {
		if s.tableAlias != "" {
			sqlParts = append(sqlParts, fmt.Sprintf("FROM %s AS %s", s.table, s.tableAlias))
		} else {
			sqlParts = append(sqlParts, fmt.Sprintf("FROM %s", s.table))
		}
	}

	// JOIN
	for _, j := range s.joins {
		var joinSQL string
		if j.subQuery != nil {
			subSQL, subArgs := j.subQuery.Build()
			args = append(args, subArgs...)
			if j.alias != "" {
				joinSQL = fmt.Sprintf("%s (%s) AS %s ON %s", j.joinType, subSQL, j.alias, j.on.GetSql())
			} else {
				joinSQL = fmt.Sprintf("%s (%s) ON %s", j.joinType, subSQL, j.on.GetSql())
			}
		} else {
			if j.alias != "" {
				joinSQL = fmt.Sprintf("%s %s AS %s ON %s", j.joinType, j.table, j.alias, j.on.GetSql())
			} else {
				joinSQL = fmt.Sprintf("%s %s ON %s", j.joinType, j.table, j.on.GetSql())
			}
		}
		sqlParts = append(sqlParts, joinSQL)
		args = append(args, j.on.GetArgs()...)
	}

	// WHERE
	if s.whereCond != nil {
		whereSQL, whereArgs := s.whereCond.Build()
		if whereSQL != "" {
			sqlParts = append(sqlParts, "WHERE "+whereSQL)
			args = append(args, whereArgs...)
		}
	}

	// GROUP BY
	if len(s.groupBy) > 0 {
		sqlParts = append(sqlParts, "GROUP BY "+strings.Join(s.groupBy, ", "))
	}

	// HAVING
	if s.havingCond != nil {
		havingSQL, havingArgs := s.havingCond.Build()
		if havingSQL != "" {
			sqlParts = append(sqlParts, "HAVING "+havingSQL)
			args = append(args, havingArgs...)
		}
	}

	// ORDER BY
	if len(s.orderBy) > 0 {
		sqlParts = append(sqlParts, "ORDER BY "+strings.Join(s.orderBy, ", "))
	}

	// LIMIT / OFFSET
	if s.limit >= 0 {
		sqlParts = append(sqlParts, fmt.Sprintf("LIMIT %d", s.limit))
	}
	if s.offset >= 0 {
		sqlParts = append(sqlParts, fmt.Sprintf("OFFSET %d", s.offset))
	}

	return strings.Join(sqlParts, " "), args
}

func (s *Select) From(table string) base.SelectBuilder {
	s.table = table
	s.tableAlias = ""
	return s
}

func (s *Select) FromAlias(table string, alias string) base.SelectBuilder {
	s.table = table
	s.tableAlias = alias
	return s
}

func (s *Select) Columns(columns ...string) base.SelectBuilder {
	s.columns = append(s.columns, columns...)
	return s
}

func (s *Select) InnerJoin(table string, on base.ConditionBuilder) base.SelectBuilder {
	s.joins = append(s.joins, joinNode{"INNER JOIN", table, "", on, nil})
	return s
}

func (s *Select) InnerJoinAlias(table string, alias string, on base.ConditionBuilder) base.SelectBuilder {
	s.joins = append(s.joins, joinNode{"INNER JOIN", table, alias, on, nil})
	return s
}

func (s *Select) InnerJoinSelect(builder base.SelectBuilder, alias string, on base.ConditionBuilder) base.SelectBuilder {
	s.joins = append(s.joins, joinNode{"INNER JOIN", "", alias, on, builder})
	return s
}

func (s *Select) FullJoin(table string, on base.ConditionBuilder) base.SelectBuilder {
	s.joins = append(s.joins, joinNode{"FULL JOIN", table, "", on, nil})
	return s
}

func (s *Select) FullJoinAlias(table string, alias string, on base.ConditionBuilder) base.SelectBuilder {
	s.joins = append(s.joins, joinNode{"FULL JOIN", table, alias, on, nil})
	return s
}

func (s *Select) FullJoinSelect(builder base.SelectBuilder, alias string, on base.ConditionBuilder) base.SelectBuilder {
	s.joins = append(s.joins, joinNode{"FULL JOIN", "", alias, on, builder})
	return s
}

func (s *Select) LeftJoin(table string, on base.ConditionBuilder) base.SelectBuilder {
	s.joins = append(s.joins, joinNode{"LEFT JOIN", table, "", on, nil})
	return s
}

func (s *Select) LeftJoinAlias(table string, alias string, on base.ConditionBuilder) base.SelectBuilder {
	s.joins = append(s.joins, joinNode{"LEFT JOIN", table, alias, on, nil})
	return s
}

func (s *Select) LeftJoinSelect(builder base.SelectBuilder, alias string, on base.ConditionBuilder) base.SelectBuilder {
	s.joins = append(s.joins, joinNode{"LEFT JOIN", "", alias, on, builder})
	return s
}

func (s *Select) RightJoin(table string, on base.ConditionBuilder) base.SelectBuilder {
	s.joins = append(s.joins, joinNode{"RIGHT JOIN", table, "", on, nil})
	return s
}

func (s *Select) RightJoinAlias(table string, alias string, on base.ConditionBuilder) base.SelectBuilder {
	s.joins = append(s.joins, joinNode{"RIGHT JOIN", table, alias, on, nil})
	return s
}

func (s *Select) RightJoinSelect(builder base.SelectBuilder, alias string, on base.ConditionBuilder) base.SelectBuilder {
	s.joins = append(s.joins, joinNode{"RIGHT JOIN", "", alias, on, builder})
	return s
}

func (s *Select) Where(cond base.ConditionBuilder) base.SelectBuilder {
	s.whereCond = cond
	return s
}

func (s *Select) GroupBy(columns ...string) base.SelectBuilder {
	s.groupBy = append(s.groupBy, columns...)
	return s
}

func (s *Select) Having(cond base.ConditionBuilder) base.SelectBuilder {
	s.havingCond = cond
	return s
}

func (s *Select) OrderBy(columns ...string) base.SelectBuilder {
	s.orderBy = append(s.orderBy, columns...)
	return s
}

func (s *Select) OrderByDesc(columns ...string) base.SelectBuilder {
	for _, col := range columns {
		s.orderBy = append(s.orderBy, fmt.Sprintf("%s DESC", col))
	}
	return s
}

func (s *Select) Limit(limit int) base.SelectBuilder {
	s.limit = limit
	return s
}

func (s *Select) Offset(offset int) base.SelectBuilder {
	s.offset = offset
	return s
}

type joinNode struct {
	joinType string
	table    string
	alias    string
	on       base.ConditionBuilder
	subQuery base.SelectBuilder
}
