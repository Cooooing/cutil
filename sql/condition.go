package sql

import "strings"

type Condition struct {
	Sql  strings.Builder
	Args []any
}

func NewCondition() ConditionBuilder {
	return &Condition{}
}

func (c *Condition) GetSql() string {
	return c.Sql.String()
}

func (c *Condition) GetArgs() []any {
	return c.Args
}

func (c *Condition) Build() (string, []any) {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) Where(cond string, args ...any) ConditionBuilder {
	return c
}

func (c *Condition) WhereIf(condition bool, cond string, args ...any) ConditionBuilder {
	if condition {
		c.Where(cond, args)
	}
	return c
}

func (c *Condition) WhereAlias(tableAlias string, column string, args ...any) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) WhereAliasIf(condition bool, tableAlias string, column string, args ...any) ConditionBuilder {
	if condition {
		c.WhereAlias(tableAlias, column, args)
	}
	return c
}

func (c *Condition) And(column string, args ...any) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) AndIf(condition bool, column string, args ...any) ConditionBuilder {
	if condition {
		c.And(column, args)
	}
	return c
}

func (c *Condition) AndAlias(tableAlias string, column string, args ...any) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) AndAliasIf(condition bool, tableAlias string, column string, args ...any) ConditionBuilder {
	if condition {
		c.AndAlias(tableAlias, column, args)
	}
	return c
}

func (c *Condition) Or(column string, args ...any) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) OrIf(condition bool, column string, args ...any) ConditionBuilder {
	if condition {
		c.Or(column, args)
	}
	return c
}

func (c *Condition) OrAlias(tableAlias string, column string, args ...any) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) OrAliasIf(condition bool, tableAlias string, column string, args ...any) ConditionBuilder {
	if condition {
		c.OrAlias(tableAlias, column, args)
	}
	return c
}

func (c *Condition) Eq(column string, args ...any) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) EqIf(condition bool, column string, args ...any) ConditionBuilder {
	if condition {
		c.Eq(column, args)
	}
	return c
}

func (c *Condition) EqAlias(tableAlias string, column string, args ...any) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) EqAliasIf(condition bool, tableAlias string, column string, args ...any) ConditionBuilder {
	if condition {
		c.EqAlias(tableAlias, column, args)
	}
	return c
}

func (c *Condition) Ne(column string, args ...any) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) NeIf(condition bool, column string, args ...any) ConditionBuilder {
	if condition {
		c.Ne(column, args)
	}
	return c
}

func (c *Condition) NeAlias(tableAlias string, column string, args ...any) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) NeAliasIf(condition bool, tableAlias string, column string, args ...any) ConditionBuilder {
	if condition {
		c.NeAlias(tableAlias, column, args)
	}
	return c
}

func (c *Condition) Gt(column string, args ...any) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) GtIf(condition bool, column string, args ...any) ConditionBuilder {
	if condition {
		c.Gt(column, args)
	}
	return c
}

func (c *Condition) GtAlias(tableAlias string, column string, args ...any) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) GtAliasIf(condition bool, tableAlias string, column string, args ...any) ConditionBuilder {
	if condition {
		c.GtAlias(tableAlias, column, args)
	}
	return c
}

func (c *Condition) Ge(column string, args ...any) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) GeIf(condition bool, column string, args ...any) ConditionBuilder {
	if condition {
		c.Ge(column, args)
	}
	return c
}

func (c *Condition) GeAlias(tableAlias string, column string, args ...any) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) GeAliasIf(condition bool, tableAlias string, column string, args ...any) ConditionBuilder {
	if condition {
		c.GeAlias(tableAlias, column, args)
	}
	return c
}

func (c *Condition) Lt(column string, args ...any) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) LtIf(condition bool, column string, args ...any) ConditionBuilder {
	if condition {
		c.Lt(column, args)
	}
	return c
}

func (c *Condition) LtAlias(tableAlias string, column string, args ...any) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) LtAliasIf(condition bool, tableAlias string, column string, args ...any) ConditionBuilder {
	if condition {
		c.LtAlias(tableAlias, column, args)
	}
	return c
}

func (c *Condition) Le(column string, args ...any) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) LeIf(condition bool, column string, args ...any) ConditionBuilder {
	if condition {
		c.Le(column, args)
	}
	return c
}

func (c *Condition) LeAlias(tableAlias string, column string, args ...any) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) LeAliasIf(condition bool, tableAlias string, column string, args ...any) ConditionBuilder {
	if condition {
		c.LeAlias(tableAlias, column, args)
	}
	return c
}

func (c *Condition) Between(column string, min any, max any) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) BetweenIf(condition bool, column string, min any, max any) ConditionBuilder {
	if condition {
		c.Between(column, min, max)
	}
	return c
}

func (c *Condition) BetweenAlias(tableAlias string, column string, min any, max any) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) BetweenAliasIf(condition bool, tableAlias string, column string, min any, max any) ConditionBuilder {
	if condition {
		c.BetweenAlias(tableAlias, column, min, max)
	}
	return c
}

func (c *Condition) NotBetween(column string, min any, max any) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) NotBetweenIf(condition bool, column string, min any, max any) ConditionBuilder {
	if condition {
		c.NotBetween(column, min, max)
	}
	return c
}

func (c *Condition) NotBetweenAlias(tableAlias string, column string, min any, max any) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) NotBetweenAliasIf(condition bool, tableAlias string, column string, min any, max any) ConditionBuilder {
	if condition {
		c.NotBetweenAlias(tableAlias, column, min, max)
	}
	return c
}

func (c *Condition) Like(column string, args ...any) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) LikeIf(condition bool, column string, args ...any) ConditionBuilder {
	if condition {
		c.Like(column, args)
	}
	return c
}

func (c *Condition) LikeAlias(tableAlias string, column string, args ...any) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) LikeAliasIf(condition bool, tableAlias string, column string, args ...any) ConditionBuilder {
	if condition {
		c.LikeAlias(tableAlias, column, args)
	}
	return c
}

func (c *Condition) LikeLeft(column string, args ...any) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) LikeLeftIf(condition bool, column string, args ...any) ConditionBuilder {
	if condition {
		c.LikeLeft(column, args)
	}
	return c
}

func (c *Condition) LikeLeftAlias(tableAlias string, column string, args ...any) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) LikeLeftAliasIf(condition bool, tableAlias string, column string, args ...any) ConditionBuilder {
	if condition {
		c.LikeLeftAlias(tableAlias, column, args)
	}
	return c
}

func (c *Condition) LikeRight(column string, args ...any) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) LikeRightIf(condition bool, column string, args ...any) ConditionBuilder {
	if condition {
		c.LikeRight(column, args)
	}
	return c
}

func (c *Condition) LikeRightAlias(tableAlias string, column string, args ...any) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) LikeRightAliasIf(condition bool, tableAlias string, column string, args ...any) ConditionBuilder {
	if condition {
		c.LikeRightAlias(tableAlias, column, args)
	}
	return c
}

func (c *Condition) NotLike(column string, args ...any) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) NotLikeIf(condition bool, column string, args ...any) ConditionBuilder {
	if condition {
		c.NotLike(column, args)
	}
	return c
}

func (c *Condition) NotLikeAlias(tableAlias string, column string, args ...any) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) NotLikeAliasIf(condition bool, tableAlias string, column string, args ...any) ConditionBuilder {
	if condition {
		c.NotLikeAlias(tableAlias, column, args)
	}
	return c
}

func (c *Condition) NotLikeLeft(column string, args ...any) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) NotLikeLeftIf(condition bool, column string, args ...any) ConditionBuilder {
	if condition {
		c.NotLikeLeft(column, args)
	}
	return c
}

func (c *Condition) NotLikeLeftAlias(tableAlias string, column string, args ...any) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) NotLikeLeftAliasIf(condition bool, tableAlias string, column string, args ...any) ConditionBuilder {
	if condition {
		c.NotLikeLeftAlias(tableAlias, column, args)
	}
	return c
}

func (c *Condition) NotLikeRight(column string, args ...any) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) NotLikeRightIf(condition bool, column string, args ...any) ConditionBuilder {
	if condition {
		c.NotLikeRight(column, args)
	}
	return c
}

func (c *Condition) NotLikeRightAlias(tableAlias string, column string, args ...any) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) NotLikeRightAliasIf(condition bool, tableAlias string, column string, args ...any) ConditionBuilder {
	if condition {
		c.NotLikeRightAlias(tableAlias, column, args)
	}
	return c
}

func (c *Condition) IsNull(column string) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) IsNullIf(condition bool, column string) ConditionBuilder {
	if condition {
		c.IsNull(column)
	}
	return c
}

func (c *Condition) IsNotNull(column string) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) IsNotNullIf(condition bool, column string) ConditionBuilder {
	if condition {
		c.IsNotNull(column)
	}
	return c
}

func (c *Condition) Exists(builder SelectBuilder) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) ExistsIf(condition bool, builder SelectBuilder) ConditionBuilder {
	if condition {
		c.Exists(builder)
	}
	return c
}

func (c *Condition) ExistsAlias(builder SelectBuilder, alias string) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) ExistsAliasIf(condition bool, builder SelectBuilder, alias string) ConditionBuilder {
	if condition {
		c.ExistsAlias(builder, alias)
	}
	return c
}

func (c *Condition) NotExists(builder SelectBuilder) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) NotExistsIf(condition bool, builder SelectBuilder) ConditionBuilder {
	if condition {
		c.NotExists(builder)
	}
	return c
}

func (c *Condition) NotExistsAlias(builder SelectBuilder, alias string) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) NotExistsAliasIf(condition bool, builder SelectBuilder, alias string) ConditionBuilder {
	if condition {
		c.NotExistsAlias(builder, alias)
	}
	return c
}

func (c *Condition) In(column string, args ...any) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) InIf(condition bool, column string, args ...any) ConditionBuilder {
	if condition {
		c.In(column, args)
	}
	return c
}

func (c *Condition) InAlias(column string, alias string, args ...any) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) InAliasIf(condition bool, column string, alias string, args ...any) ConditionBuilder {
	if condition {
		c.InAlias(column, alias, args)
	}
	return c
}

func (c *Condition) NotIn(column string, args ...any) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) NotInIf(condition bool, column string, args ...any) ConditionBuilder {
	if condition {
		c.NotIn(column, args)
	}
	return c
}

func (c *Condition) NotInAlias(column string, alias string, args ...any) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) NotInAliasIf(condition bool, column string, alias string, args ...any) ConditionBuilder {
	if condition {
		c.NotInAlias(column, alias, args)
	}
	return c
}

func (c *Condition) Nested(cond ConditionBuilder) ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) NestedIf(condition bool, cond ConditionBuilder) ConditionBuilder {
	if condition {
		c.Nested(cond)
	}
	return c
}
