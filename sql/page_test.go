package sql

import (
	"database/sql"
	"encoding/json"
	"github.com/Cooooing/cutil/common/logger"
	"log"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

const query = `select * from "user" order by id`

var args = []any{}

var db *sql.DB

type User struct {
	Id         int        `json:"id"`
	Username   string     `json:"username"`
	Email      string     `json:"email"`
	CreateTime *time.Time `json:"create_time"`
}

func init() {
	var err error
	db, err = sql.Open("postgres", "host=127.0.0.1 user=root password=123456 dbname=public port=5432 sslmode=disable TimeZone=Asia/Shanghai")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	logger.Info("connect to database success")
}

func TestPageQueryForStruct(t *testing.T) {
	res, err := PageQueryForStruct[User](db, PageReq{Page: 1, Size: 100}, query, args...)
	if err != nil {
		t.Error(err)
	}
	marshal, _ := json.Marshal(res)
	t.Log(string(marshal))
}

func TestPageQueryForMap(t *testing.T) {
	res, err := PageQueryForMap(db, PageReq{Page: 1, Size: 100}, query, args...)
	if err != nil {
		t.Error(err)
	}
	marshal, _ := json.Marshal(res)
	t.Log(string(marshal))
}

func TestPageQueryForStructWithLimitOffset(t *testing.T) {
	res, err := PageQueryForStructWithLimitOffset[User](db, PageReq{Page: 1, Size: 100}, query, args...)
	if err != nil {
		t.Error(err)
	}
	marshal, _ := json.Marshal(res)
	t.Log(string(marshal))
}

func TestPageQueryForMapWithLimitOffset(t *testing.T) {
	res, err := PageQueryForMapWithLimitOffset(db, PageReq{Page: 1, Size: 100}, query, args...)
	if err != nil {
		t.Error(err)
	}
	marshal, _ := json.Marshal(res)
	t.Log(string(marshal))
}

func TestPageQueryForStructWithRowNumber(t *testing.T) {
	res, err := PageQueryForStructWithRowNumber[User](db, PageReq{Page: 1, Size: 100}, query, args...)
	if err != nil {
		t.Error(err)
	}
	marshal, _ := json.Marshal(res)
	t.Log(string(marshal))
}

func TestPageQueryForMapWithRowNumber(t *testing.T) {
	res, err := PageQueryForMapWithRowNumber(db, PageReq{Page: 1, Size: 100}, query, args...)
	if err != nil {
		t.Error(err)
	}
	marshal, _ := json.Marshal(res)
	t.Log(string(marshal))
}

func TestPageQueryForStructWithFetchOffset(t *testing.T) {
	res, err := PageQueryForStructWithFetchOffset[User](db, PageReq{Page: 1, Size: 100}, query, args...)
	if err != nil {
		t.Error(err)
	}
	marshal, _ := json.Marshal(res)
	t.Log(string(marshal))
}

func TestPageQueryForMapWithFetchOffset(t *testing.T) {
	res, err := PageQueryForMapWithFetchOffset(db, PageReq{Page: 1, Size: 100}, query, args...)
	if err != nil {
		t.Error(err)
	}
	marshal, _ := json.Marshal(res)
	t.Log(string(marshal))
}

func TestPageQueryForStructWithDeclareCursor(t *testing.T) {
	res, err := PageQueryForStructWithDeclareCursor[User](db, PageReq{Page: 1, Size: 100}, query, args...)
	if err != nil {
		t.Error(err)
	}
	marshal, _ := json.Marshal(res)
	t.Log(string(marshal))
}

func TestPageQueryForMapWithDeclareCursor(t *testing.T) {
	res, err := PageQueryForMapWithDeclareCursor(db, PageReq{Page: 1, Size: 100}, query, args...)
	if err != nil {
		t.Error(err)
	}
	marshal, _ := json.Marshal(res)
	t.Log(string(marshal))
}
