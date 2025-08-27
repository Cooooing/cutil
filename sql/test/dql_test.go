package test

import (
	"encoding/json"
	"github.com/Cooooing/cutil/common/logger"
	"github.com/Cooooing/cutil/sql"
	"github.com/Cooooing/cutil/sql/dql"
	"testing"
)

func TestSimpleSelect(t *testing.T) {
	Init()

	/*
		SELECT id, name, age, email, created_at
		FROM users
		WHERE age > 20
	*/
	users, err := sql.WithExecutor[User](
		DB,
		dql.NewSelect().
			Columns("id", "name", "age", "email", "created_at").
			From("users").
			Where(dql.NewCondition().Gt("age", 20)),
	).Debug().List()
	if err != nil {
		t.Error(err)
	}
	logger.Info("users: %+v", len(users))
	bytes, _ := json.Marshal(users)
	logger.Info("users: %s", string(bytes))

}

func TestJoinSelect(t *testing.T) {
	Init()
	var (
		err   error
		res   []PostTitle
		bytes []byte
	)
	/*
		SELECT u.name, p.title
		FROM users u
				 LEFT JOIN posts p ON u.id = p.user_id
	*/
	res, err = sql.WithExecutor[PostTitle](
		DB,
		dql.NewSelect().
			Columns("u.name", "p.title").
			From("users u").
			LeftJoin("posts", "p", dql.NewCondition().On("u.id", "p.user_id")),
	).Debug().List()
	if err != nil {
		t.Error(err)
	}
	logger.Info("users: %+v", len(res))
	bytes, _ = json.Marshal(res)
	logger.Info("users: %s", string(bytes))
	/*
		SELECT u.id, u.name, p.title
		FROM users u
				 LEFT JOIN (SELECT user_id, title FROM posts) AS p
						   ON u.id = p.user_id
	*/
	res, err = sql.WithExecutor[PostTitle](
		DB,
		dql.NewSelect().
			Columns("u.id", "u.name", "p.title").
			From("users u").
			LeftJoinSelect(dql.NewSelect().Columns("user_id", "title").From("posts"), "p", dql.NewCondition().On("u.id", "p.user_id")),
	).Debug().List()
	if err != nil {
		t.Error(err)
	}
	logger.Info("users: %+v", len(res))
	bytes, _ = json.Marshal(res)
	logger.Info("users: %s", string(bytes))

}

func TestNestedSelect(t *testing.T) {
	Init()
	var (
		err   error
		res   []PostTitle
		bytes []byte
	)
	/*
		SELECT *
		FROM users
		WHERE (name = 'Alice' OR name = 'Bob')
	*/
	res, err = sql.WithExecutor[PostTitle](
		DB,
		dql.NewSelect().
			From("users").
			Where(
				dql.NewCondition().
					Nested(
						dql.NewCondition().Eq("name", "Alice").Or().Eq("name", "Bob"),
					),
			),
	).Debug().List()
	if err != nil {
		t.Error(err)
	}
	logger.Info("users: %+v", len(res))
	bytes, _ = json.Marshal(res)
	logger.Info("users: %s", string(bytes))
}
