package dql

import (
	"fmt"
	"github.com/Cooooing/cutil/sql/base"
	"strings"
)

type conditionNode struct {
	expr string
	args []any
	op   string
}

type Condition struct {
	nodes []conditionNode
}

func NewCondition() base.ConditionBuilder {
	return &Condition{}
}

func (c *Condition) GetSql() string {
	sql, _ := c.Build()
	return sql
}

func (c *Condition) GetArgs() []any {
	_, args := c.Build()
	return args
}

func (c *Condition) append(expr string, args ...any) base.ConditionBuilder {
	op := ""
	if len(c.nodes) > 0 {
		// 默认逻辑符为 AND
		op = "AND"
		// 如果前面调用了 Or()，修改最后节点 op 为 OR
		if len(c.nodes) > 0 && c.nodes[len(c.nodes)-1].op == "OR_PENDING" {
			op = "OR"
			c.nodes[len(c.nodes)-1].op = "" // 重置标记
		}
	}
	c.nodes = append(c.nodes, conditionNode{
		expr: expr,
		args: args,
		op:   op,
	})
	return c
}

func (c *Condition) nextOp() string {
	if len(c.nodes) == 0 {
		return ""
	}
	return "AND"
}

func (c *Condition) Build() (string, []any) {
	if len(c.nodes) == 0 {
		return "", nil
	}

	var sqlParts []string
	var args []any

	for i, node := range c.nodes {
		part := ""
		if i > 0 {
			part += " " + node.op + " "
		}
		part += node.expr
		args = append(args, node.args...)
		sqlParts = append(sqlParts, part)
	}

	return strings.Join(sqlParts, ""), args
}

func (c *Condition) Where(cond string, args ...any) base.ConditionBuilder {
	return c
}

func (c *Condition) WhereIf(condition bool, cond string, args ...any) base.ConditionBuilder {
	if condition {
		c.Where(cond, args)
	}
	return c
}

func (c *Condition) WhereAlias(tableAlias string, column string, args ...any) base.ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) WhereAliasIf(condition bool, tableAlias string, column string, args ...any) base.ConditionBuilder {
	if condition {
		c.WhereAlias(tableAlias, column, args)
	}
	return c
}

func (c *Condition) And() base.ConditionBuilder {
	if len(c.nodes) > 0 && c.nodes[len(c.nodes)-1].op == "" {
		c.nodes[len(c.nodes)-1].op = "AND"
	}
	return c
}

func (c *Condition) AndIf(condition bool) base.ConditionBuilder {
	if condition {
		c.And()
	}
	return c
}

func (c *Condition) Or() base.ConditionBuilder {
	if len(c.nodes) > 0 {
		c.nodes[len(c.nodes)-1].op = "OR_PENDING"
	}
	return c
}

func (c *Condition) OrIf(condition bool) base.ConditionBuilder {
	if condition {
		c.Or()
	}
	return c
}

func (c *Condition) Eq(column string, args any) base.ConditionBuilder {
	return c.append(fmt.Sprintf("%s = ?", column), args)
}

func (c *Condition) EqIf(condition bool, column string, args any) base.ConditionBuilder {
	if condition {
		c.Eq(column, args)
	}
	return c
}

func (c *Condition) EqAlias(tableAlias string, column string, args any) base.ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) EqAliasIf(condition bool, tableAlias string, column string, args any) base.ConditionBuilder {
	if condition {
		c.EqAlias(tableAlias, column, args)
	}
	return c
}

func (c *Condition) Ne(column string, args any) base.ConditionBuilder {
	return c.append(fmt.Sprintf("%s <> ?", column), args)
}

func (c *Condition) NeIf(condition bool, column string, args any) base.ConditionBuilder {
	if condition {
		c.Ne(column, args)
	}
	return c
}

func (c *Condition) NeAlias(tableAlias string, column string, args any) base.ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) NeAliasIf(condition bool, tableAlias string, column string, args any) base.ConditionBuilder {
	if condition {
		c.NeAlias(tableAlias, column, args)
	}
	return c
}

func (c *Condition) Gt(column string, args any) base.ConditionBuilder {
	return c.append(fmt.Sprintf("%s > ?", column), args)
}

func (c *Condition) GtIf(condition bool, column string, args any) base.ConditionBuilder {
	if condition {
		c.Gt(column, args)
	}
	return c
}

func (c *Condition) GtAlias(tableAlias string, column string, args any) base.ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) GtAliasIf(condition bool, tableAlias string, column string, args any) base.ConditionBuilder {
	if condition {
		c.GtAlias(tableAlias, column, args)
	}
	return c
}

func (c *Condition) Ge(column string, args any) base.ConditionBuilder {
	return c.append(fmt.Sprintf("%s >= ?", column), args)
}

func (c *Condition) GeIf(condition bool, column string, args any) base.ConditionBuilder {
	if condition {
		c.Ge(column, args)
	}
	return c
}

func (c *Condition) GeAlias(tableAlias string, column string, args any) base.ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) GeAliasIf(condition bool, tableAlias string, column string, args any) base.ConditionBuilder {
	if condition {
		c.GeAlias(tableAlias, column, args)
	}
	return c
}

func (c *Condition) Lt(column string, args any) base.ConditionBuilder {
	return c.append(fmt.Sprintf("%s < ?", column), args)
}

func (c *Condition) LtIf(condition bool, column string, args any) base.ConditionBuilder {
	if condition {
		c.Lt(column, args)
	}
	return c
}

func (c *Condition) LtAlias(tableAlias string, column string, args any) base.ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) LtAliasIf(condition bool, tableAlias string, column string, args any) base.ConditionBuilder {
	if condition {
		c.LtAlias(tableAlias, column, args)
	}
	return c
}

func (c *Condition) Le(column string, args any) base.ConditionBuilder {
	return c.append(fmt.Sprintf("%s <= ?", column), args)
}

func (c *Condition) LeIf(condition bool, column string, args any) base.ConditionBuilder {
	if condition {
		c.Le(column, args)
	}
	return c
}

func (c *Condition) LeAlias(tableAlias string, column string, args any) base.ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) LeAliasIf(condition bool, tableAlias string, column string, args any) base.ConditionBuilder {
	if condition {
		c.LeAlias(tableAlias, column, args)
	}
	return c
}

func (c *Condition) Between(column string, min any, max any) base.ConditionBuilder {
	return c.append(fmt.Sprintf("%s >= ? and %s <= ?", column, column), min, max)
}

func (c *Condition) BetweenIf(condition bool, column string, min any, max any) base.ConditionBuilder {
	if condition {
		c.Between(column, min, max)
	}
	return c
}

func (c *Condition) BetweenAlias(tableAlias string, column string, min any, max any) base.ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) BetweenAliasIf(condition bool, tableAlias string, column string, min any, max any) base.ConditionBuilder {
	if condition {
		c.BetweenAlias(tableAlias, column, min, max)
	}
	return c
}

func (c *Condition) NotBetween(column string, min any, max any) base.ConditionBuilder {
	return c.append(fmt.Sprintf("(%s < ? or %s > ?)", column, column), min, max)
}

func (c *Condition) NotBetweenIf(condition bool, column string, min any, max any) base.ConditionBuilder {
	if condition {
		c.NotBetween(column, min, max)
	}
	return c
}

func (c *Condition) NotBetweenAlias(tableAlias string, column string, min any, max any) base.ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) NotBetweenAliasIf(condition bool, tableAlias string, column string, min any, max any) base.ConditionBuilder {
	if condition {
		c.NotBetweenAlias(tableAlias, column, min, max)
	}
	return c
}

func (c *Condition) Like(column string, args any) base.ConditionBuilder {
	return c.append(fmt.Sprintf("%s LIKE ?", column), fmt.Sprintf("%%%v%%", args))
}

func (c *Condition) LikeIf(condition bool, column string, args any) base.ConditionBuilder {
	if condition {
		c.Like(column, args)
	}
	return c
}

func (c *Condition) LikeAlias(tableAlias string, column string, args any) base.ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) LikeAliasIf(condition bool, tableAlias string, column string, args any) base.ConditionBuilder {
	if condition {
		c.LikeAlias(tableAlias, column, args)
	}
	return c
}

func (c *Condition) LikeLeft(column string, args any) base.ConditionBuilder {
	return c.append(fmt.Sprintf("%s LIKE ?", column), fmt.Sprintf("%%%v", args))
}

func (c *Condition) LikeLeftIf(condition bool, column string, args any) base.ConditionBuilder {
	if condition {
		c.LikeLeft(column, args)
	}
	return c
}

func (c *Condition) LikeLeftAlias(tableAlias string, column string, args any) base.ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) LikeLeftAliasIf(condition bool, tableAlias string, column string, args any) base.ConditionBuilder {
	if condition {
		c.LikeLeftAlias(tableAlias, column, args)
	}
	return c
}

func (c *Condition) LikeRight(column string, args any) base.ConditionBuilder {
	return c.append(fmt.Sprintf("%s LIKE ?", column), fmt.Sprintf("%v%%", args))
}

func (c *Condition) LikeRightIf(condition bool, column string, args any) base.ConditionBuilder {
	if condition {
		c.LikeRight(column, args)
	}
	return c
}

func (c *Condition) LikeRightAlias(tableAlias string, column string, args any) base.ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) LikeRightAliasIf(condition bool, tableAlias string, column string, args any) base.ConditionBuilder {
	if condition {
		c.LikeRightAlias(tableAlias, column, args)
	}
	return c
}

func (c *Condition) NotLike(column string, args any) base.ConditionBuilder {
	return c.append(fmt.Sprintf("%s NOT LIKE ?", column), fmt.Sprintf("%%%v%%", args))
}

func (c *Condition) NotLikeIf(condition bool, column string, args any) base.ConditionBuilder {
	if condition {
		c.NotLike(column, args)
	}
	return c
}

func (c *Condition) NotLikeAlias(tableAlias string, column string, args any) base.ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) NotLikeAliasIf(condition bool, tableAlias string, column string, args any) base.ConditionBuilder {
	if condition {
		c.NotLikeAlias(tableAlias, column, args)
	}
	return c
}

func (c *Condition) NotLikeLeft(column string, args any) base.ConditionBuilder {
	return c.append(fmt.Sprintf("%s NOT LIKE ?", column), fmt.Sprintf("%%%v", args))
}

func (c *Condition) NotLikeLeftIf(condition bool, column string, args any) base.ConditionBuilder {
	if condition {
		c.NotLikeLeft(column, args)
	}
	return c
}

func (c *Condition) NotLikeLeftAlias(tableAlias string, column string, args any) base.ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) NotLikeLeftAliasIf(condition bool, tableAlias string, column string, args any) base.ConditionBuilder {
	if condition {
		c.NotLikeLeftAlias(tableAlias, column, args)
	}
	return c
}

func (c *Condition) NotLikeRight(column string, args any) base.ConditionBuilder {
	return c.append(fmt.Sprintf("%s NOT LIKE ?", column), fmt.Sprintf("%v%%", args))
}

func (c *Condition) NotLikeRightIf(condition bool, column string, args any) base.ConditionBuilder {
	if condition {
		c.NotLikeRight(column, args)
	}
	return c
}

func (c *Condition) NotLikeRightAlias(tableAlias string, column string, args any) base.ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) NotLikeRightAliasIf(condition bool, tableAlias string, column string, args any) base.ConditionBuilder {
	if condition {
		c.NotLikeRightAlias(tableAlias, column, args)
	}
	return c
}

func (c *Condition) IsNull(column string) base.ConditionBuilder {
	return c.append(fmt.Sprintf("%s IS NULL", column))
}

func (c *Condition) IsNullIf(condition bool, column string) base.ConditionBuilder {
	if condition {
		c.IsNull(column)
	}
	return c
}

func (c *Condition) IsNotNull(column string) base.ConditionBuilder {
	return c.append(fmt.Sprintf("%s IS NOT NULL", column))
}

func (c *Condition) IsNotNullIf(condition bool, column string) base.ConditionBuilder {
	if condition {
		c.IsNotNull(column)
	}
	return c
}

func (c *Condition) Exists(builder base.SelectBuilder) base.ConditionBuilder {
	sql, args := builder.Build()
	return c.append(fmt.Sprintf("EXISTS (%s)", sql), args...)
}

func (c *Condition) ExistsIf(condition bool, builder base.SelectBuilder) base.ConditionBuilder {
	if condition {
		c.Exists(builder)
	}
	return c
}

func (c *Condition) ExistsAlias(builder base.SelectBuilder, alias string) base.ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) ExistsAliasIf(condition bool, builder base.SelectBuilder, alias string) base.ConditionBuilder {
	if condition {
		c.ExistsAlias(builder, alias)
	}
	return c
}

func (c *Condition) NotExists(builder base.SelectBuilder) base.ConditionBuilder {
	sql, args := builder.Build()
	return c.append(fmt.Sprintf("NOT EXISTS (%s)", sql), args...)
}

func (c *Condition) NotExistsIf(condition bool, builder base.SelectBuilder) base.ConditionBuilder {
	if condition {
		c.NotExists(builder)
	}
	return c
}

func (c *Condition) NotExistsAlias(builder base.SelectBuilder, alias string) base.ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) NotExistsAliasIf(condition bool, builder base.SelectBuilder, alias string) base.ConditionBuilder {
	if condition {
		c.NotExistsAlias(builder, alias)
	}
	return c
}

func (c *Condition) In(column string, args ...any) base.ConditionBuilder {
	var placeholder string
	for _, _ = range args {
		placeholder += " ?,"
	}
	placeholder = placeholder[:len(placeholder)-1]
	return c.append(fmt.Sprintf("%s IN (%s)", column, placeholder), args...)
}

func (c *Condition) InIf(condition bool, column string, args ...any) base.ConditionBuilder {
	if condition {
		c.In(column, args)
	}
	return c
}

func (c *Condition) InAlias(column string, alias string, args ...any) base.ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) InAliasIf(condition bool, column string, alias string, args ...any) base.ConditionBuilder {
	if condition {
		c.InAlias(column, alias, args)
	}
	return c
}

func (c *Condition) NotIn(column string, args ...any) base.ConditionBuilder {
	var placeholder string
	for _, _ = range args {
		placeholder += " ?,"
	}
	placeholder = placeholder[:len(placeholder)-1]
	return c.append(fmt.Sprintf("%s NOT IN (%s)", column, placeholder), args...)
}

func (c *Condition) NotInIf(condition bool, column string, args ...any) base.ConditionBuilder {
	if condition {
		c.NotIn(column, args)
	}
	return c
}

func (c *Condition) NotInAlias(column string, alias string, args ...any) base.ConditionBuilder {
	// TODO implement me
	panic("implement me")
}

func (c *Condition) NotInAliasIf(condition bool, column string, alias string, args ...any) base.ConditionBuilder {
	if condition {
		c.NotInAlias(column, alias, args)
	}
	return c
}

func (c *Condition) Nested(cond base.ConditionBuilder) base.ConditionBuilder {
	sql, args := cond.Build()
	return c.append(fmt.Sprintf("(%s)", sql), args...)
}

func (c *Condition) NestedIf(condition bool, cond base.ConditionBuilder) base.ConditionBuilder {
	if condition {
		c.Nested(cond)
	}
	return c
}

func (c *Condition) On(columnA string, columnB string) base.ConditionBuilder {
	return c.append(fmt.Sprintf("%s = %s", columnA, columnB))
}
func (c *Condition) OnIf(condition bool, columnA string, columnB string) base.ConditionBuilder {
	if condition {
		c.On(columnA, columnB)
	}
	return c
}

func (c *Condition) OnAlias(aliasA string, columnA string, aliasB string, columnB string) base.ConditionBuilder {
	return c.append(fmt.Sprintf("%s.%s = %s.%s", aliasA, columnA, aliasB, columnB))
}

func (c *Condition) OnAliasIf(condition bool, aliasA string, columnA string, aliasB string, columnB string) base.ConditionBuilder {
	if condition {
		c.OnAlias(aliasA, columnA, aliasB, columnB)
	}
	return c
}
