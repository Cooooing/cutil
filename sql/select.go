package sql

import "strings"

type Select struct {
	where strings.Builder
	args  []any
}

func (s Select) GetSql() string {
	// TODO implement me
	panic("implement me")
}

func (s Select) GetArgs() []any {
	// TODO implement me
	panic("implement me")
}

func (s Select) Build() (string, []any) {
	// TODO implement me
	panic("implement me")
}

func (s Select) From(table string) SelectBuilder {
	// TODO implement me
	panic("implement me")
}

func (s Select) FromAlias(table string, alias string) SelectBuilder {
	// TODO implement me
	panic("implement me")
}

func (s Select) Columns(columns ...string) SelectBuilder {
	// TODO implement me
	panic("implement me")
}

func (s Select) InnerJoin(table string, on ConditionBuilder) SelectBuilder {
	// TODO implement me
	panic("implement me")
}

func (s Select) InnerJoinAlias(table string, alias string, on ConditionBuilder) SelectBuilder {
	// TODO implement me
	panic("implement me")
}

func (s Select) InnerJoinSelect(builder SelectBuilder, alias string, on ConditionBuilder) SelectBuilder {
	// TODO implement me
	panic("implement me")
}

func (s Select) FullJoin(table string, on ConditionBuilder) SelectBuilder {
	// TODO implement me
	panic("implement me")
}

func (s Select) FullJoinAlias(table string, alias string, on ConditionBuilder) SelectBuilder {
	// TODO implement me
	panic("implement me")
}

func (s Select) FullJoinSelect(builder SelectBuilder, alias string, on ConditionBuilder) SelectBuilder {
	// TODO implement me
	panic("implement me")
}

func (s Select) LeftJoin(table string, on ConditionBuilder) SelectBuilder {
	// TODO implement me
	panic("implement me")
}

func (s Select) LeftJoinAlias(table string, alias string, on ConditionBuilder) SelectBuilder {
	// TODO implement me
	panic("implement me")
}

func (s Select) LeftJoinSelect(builder SelectBuilder, alias string, on ConditionBuilder) SelectBuilder {
	// TODO implement me
	panic("implement me")
}

func (s Select) RightJoin(table string, on ConditionBuilder) SelectBuilder {
	// TODO implement me
	panic("implement me")
}

func (s Select) RightJoinAlias(table string, alias string, on ConditionBuilder) SelectBuilder {
	// TODO implement me
	panic("implement me")
}

func (s Select) RightJoinSelect(builder SelectBuilder, alias string, on ConditionBuilder) SelectBuilder {
	// TODO implement me
	panic("implement me")
}

func (s Select) Where(cond ConditionBuilder) SelectBuilder {
	// TODO implement me
	panic("implement me")
}

func (s Select) GroupBy(columns ...string) SelectBuilder {
	// TODO implement me
	panic("implement me")
}

func (s Select) Having(cond ConditionBuilder) SelectBuilder {
	// TODO implement me
	panic("implement me")
}

func (s Select) OrderBy(columns ...string) SelectBuilder {
	// TODO implement me
	panic("implement me")
}

func (s Select) OrderByDesc(columns ...string) SelectBuilder {
	// TODO implement me
	panic("implement me")
}

func (s Select) Limit(limit int) SelectBuilder {
	// TODO implement me
	panic("implement me")
}

func (s Select) Offset(offset int) SelectBuilder {
	// TODO implement me
	panic("implement me")
}
