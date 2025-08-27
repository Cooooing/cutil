package base

// PageRespInterface 分页查询参数接口
type PageRespInterface[T any] interface {
	SetList(data []T)
	SetTotal(total int)
	SetPageReq(pageReq PageReqInterface)
	GetList() []T
	GetTotal() int
	GetPage() int
	GetSize() int
}

// PageReqInterface 分页查询返回值接口
type PageReqInterface interface {
	Validate() error
	GetPage() int
	GetSize() int
}

// BaseBuilder 通用方法接口
type BaseBuilder interface {
	GetSql() string
	GetArgs() []any
	Build() (string, []any)
}

type ConditionBuilder interface {
	BaseBuilder

	Where(column string, args ...any) ConditionBuilder
	WhereIf(condition bool, column string, args ...any) ConditionBuilder
	WhereAlias(tableAlias string, column string, args ...any) ConditionBuilder
	WhereAliasIf(condition bool, tableAlias string, column string, args ...any) ConditionBuilder

	And() ConditionBuilder
	AndIf(condition bool) ConditionBuilder

	Or() ConditionBuilder
	OrIf(condition bool) ConditionBuilder

	Eq(column string, args any) ConditionBuilder
	EqIf(condition bool, column string, args any) ConditionBuilder
	EqAlias(tableAlias string, column string, args any) ConditionBuilder
	EqAliasIf(condition bool, tableAlias string, column string, args any) ConditionBuilder

	Ne(column string, args any) ConditionBuilder
	NeIf(condition bool, column string, args any) ConditionBuilder
	NeAlias(tableAlias string, column string, args any) ConditionBuilder
	NeAliasIf(condition bool, tableAlias string, column string, args any) ConditionBuilder

	Gt(column string, args any) ConditionBuilder
	GtIf(condition bool, column string, args any) ConditionBuilder
	GtAlias(tableAlias string, column string, args any) ConditionBuilder
	GtAliasIf(condition bool, tableAlias string, column string, args any) ConditionBuilder

	Ge(column string, args any) ConditionBuilder
	GeIf(condition bool, column string, args any) ConditionBuilder
	GeAlias(tableAlias string, column string, args any) ConditionBuilder
	GeAliasIf(condition bool, tableAlias string, column string, args any) ConditionBuilder

	Lt(column string, args any) ConditionBuilder
	LtIf(condition bool, column string, args any) ConditionBuilder
	LtAlias(tableAlias string, column string, args any) ConditionBuilder
	LtAliasIf(condition bool, tableAlias string, column string, args any) ConditionBuilder

	Le(column string, args any) ConditionBuilder
	LeIf(condition bool, column string, args any) ConditionBuilder
	LeAlias(tableAlias string, column string, args any) ConditionBuilder
	LeAliasIf(condition bool, tableAlias string, column string, args any) ConditionBuilder

	Between(column string, min any, max any) ConditionBuilder
	BetweenIf(condition bool, column string, min any, max any) ConditionBuilder
	BetweenAlias(tableAlias string, column string, min any, max any) ConditionBuilder
	BetweenAliasIf(condition bool, tableAlias string, column string, min any, max any) ConditionBuilder

	NotBetween(column string, min any, max any) ConditionBuilder
	NotBetweenIf(condition bool, column string, min any, max any) ConditionBuilder
	NotBetweenAlias(tableAlias string, column string, min any, max any) ConditionBuilder
	NotBetweenAliasIf(condition bool, tableAlias string, column string, min any, max any) ConditionBuilder

	Like(column string, args any) ConditionBuilder
	LikeIf(condition bool, column string, args any) ConditionBuilder
	LikeAlias(tableAlias string, column string, args any) ConditionBuilder
	LikeAliasIf(condition bool, tableAlias string, column string, args any) ConditionBuilder

	LikeLeft(column string, args any) ConditionBuilder
	LikeLeftIf(condition bool, column string, args any) ConditionBuilder
	LikeLeftAlias(tableAlias string, column string, args any) ConditionBuilder
	LikeLeftAliasIf(condition bool, tableAlias string, column string, args any) ConditionBuilder

	LikeRight(column string, args any) ConditionBuilder
	LikeRightIf(condition bool, column string, args any) ConditionBuilder
	LikeRightAlias(tableAlias string, column string, args any) ConditionBuilder
	LikeRightAliasIf(condition bool, tableAlias string, column string, args any) ConditionBuilder

	NotLike(column string, args any) ConditionBuilder
	NotLikeIf(condition bool, column string, args any) ConditionBuilder
	NotLikeAlias(tableAlias string, column string, args any) ConditionBuilder
	NotLikeAliasIf(condition bool, tableAlias string, column string, args any) ConditionBuilder

	NotLikeLeft(column string, args any) ConditionBuilder
	NotLikeLeftIf(condition bool, column string, args any) ConditionBuilder
	NotLikeLeftAlias(tableAlias string, column string, args any) ConditionBuilder
	NotLikeLeftAliasIf(condition bool, tableAlias string, column string, args any) ConditionBuilder

	NotLikeRight(column string, args any) ConditionBuilder
	NotLikeRightIf(condition bool, column string, args any) ConditionBuilder
	NotLikeRightAlias(tableAlias string, column string, args any) ConditionBuilder
	NotLikeRightAliasIf(condition bool, tableAlias string, column string, args any) ConditionBuilder

	IsNull(column string) ConditionBuilder
	IsNullIf(condition bool, column string) ConditionBuilder
	IsNotNull(column string) ConditionBuilder
	IsNotNullIf(condition bool, column string) ConditionBuilder

	Exists(builder SelectBuilder) ConditionBuilder
	ExistsIf(condition bool, builder SelectBuilder) ConditionBuilder
	ExistsAlias(builder SelectBuilder, alias string) ConditionBuilder
	ExistsAliasIf(condition bool, builder SelectBuilder, alias string) ConditionBuilder

	NotExists(builder SelectBuilder) ConditionBuilder
	NotExistsIf(condition bool, builder SelectBuilder) ConditionBuilder
	NotExistsAlias(builder SelectBuilder, alias string) ConditionBuilder
	NotExistsAliasIf(condition bool, builder SelectBuilder, alias string) ConditionBuilder

	In(column string, args ...any) ConditionBuilder
	InIf(condition bool, column string, args ...any) ConditionBuilder
	InAlias(column string, alias string, args ...any) ConditionBuilder
	InAliasIf(condition bool, column string, alias string, args ...any) ConditionBuilder

	NotIn(column string, args ...any) ConditionBuilder
	NotInIf(condition bool, column string, args ...any) ConditionBuilder
	NotInAlias(column string, alias string, args ...any) ConditionBuilder
	NotInAliasIf(condition bool, column string, alias string, args ...any) ConditionBuilder

	Nested(cond ConditionBuilder) ConditionBuilder
	NestedIf(condition bool, cond ConditionBuilder) ConditionBuilder

	On(columnA string, columnB string) ConditionBuilder
	OnIf(condition bool, columnA string, columnB string) ConditionBuilder
	OnAlias(columnA string, aliasA string, columnB string, aliasB string) ConditionBuilder
	OnAliasIf(condition bool, columnA string, aliasA string, columnB string, aliasB string) ConditionBuilder
}

type SelectBuilder interface {
	BaseBuilder

	From(table string) SelectBuilder
	FromAlias(table string, alias string) SelectBuilder
	Columns(columns ...string) SelectBuilder

	InnerJoin(table string, on ConditionBuilder) SelectBuilder
	InnerJoinAlias(table string, alias string, on ConditionBuilder) SelectBuilder
	InnerJoinSelect(builder SelectBuilder, alias string, on ConditionBuilder) SelectBuilder

	FullJoin(table string, on ConditionBuilder) SelectBuilder
	FullJoinAlias(table string, alias string, on ConditionBuilder) SelectBuilder
	FullJoinSelect(builder SelectBuilder, alias string, on ConditionBuilder) SelectBuilder

	LeftJoin(table string, on ConditionBuilder) SelectBuilder
	LeftJoinAlias(table string, alias string, on ConditionBuilder) SelectBuilder
	LeftJoinSelect(builder SelectBuilder, alias string, on ConditionBuilder) SelectBuilder

	RightJoin(table string, on ConditionBuilder) SelectBuilder
	RightJoinAlias(table string, alias string, on ConditionBuilder) SelectBuilder
	RightJoinSelect(builder SelectBuilder, alias string, on ConditionBuilder) SelectBuilder

	Where(cond ConditionBuilder) SelectBuilder
	GroupBy(columns ...string) SelectBuilder
	Having(cond ConditionBuilder) SelectBuilder

	OrderBy(columns ...string) SelectBuilder
	OrderByDesc(columns ...string) SelectBuilder

	Limit(limit int) SelectBuilder
	Offset(offset int) SelectBuilder
}

type UpdateBuilder interface {
	BaseBuilder

	Table(table string) UpdateBuilder
	TableAlias(table string, alias string) UpdateBuilder
	Set(column string, value any) UpdateBuilder
	Where(cond ConditionBuilder) UpdateBuilder
}

type InsertBuilder interface {
	BaseBuilder

	Into(table string) InsertBuilder
	Columns(cols ...string) InsertBuilder
	Values(values ...any) InsertBuilder
	Select(builder SelectBuilder) InsertBuilder
}

type DeleteBuilder interface {
	BaseBuilder

	From(table string) DeleteBuilder
	FromAlias(table string, alias string) DeleteBuilder
	Where(cond ConditionBuilder) DeleteBuilder
}
