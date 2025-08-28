package test

import (
	"encoding/json"
	"github.com/Cooooing/cutil/common/logger"
	"github.com/Cooooing/cutil/query"
	"github.com/Cooooing/cutil/query/dml"
	"testing"
)

func TestInsert(t *testing.T) {
	Init()
	var (
		err   error
		res   []*PostTitle
		bytes []byte
	)
	/*
		SELECT u.name, p.title
		FROM users u
				 LEFT JOIN posts p ON u.id = p.user_id
	*/
	res, err = sql.WithExecutor[PostTitle](
		DB,
		dml.NewInsert().
			Into("users").
			Columns("name", "age", "email").
			Values("David", 28, "david@example.com"),
	).Debug().List()
	if err != nil {
		t.Error(err)
	}
	logger.Info("users: %+v", len(res))
	bytes, _ = json.Marshal(res)
	logger.Info("users: %s", string(bytes))

}
