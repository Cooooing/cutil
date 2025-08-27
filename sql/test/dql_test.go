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
