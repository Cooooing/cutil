package test

import (
	"encoding/json"
	sql2 "github.com/Cooooing/cutil/query"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"testing"
)

const query = `select * from "user" where id in $1 order by id`

var args = []any{[]int{1, 2}}

func TestPageQueryForStruct(t *testing.T) {
	Init()
	res, err := sql2.PageQueryForStruct[User](DB, nil, query, args...)
	if err != nil {
		t.Error(err)
	}
	marshal, _ := json.Marshal(res)
	t.Log(string(marshal))
}

func TestPageQueryForMap(t *testing.T) {
	Init()
	res, err := sql2.PageQueryForMap(DB, &sql2.PageReq{Page: 1, Size: 100}, query, args...)
	if err != nil {
		t.Error(err)
	}
	marshal, _ := json.Marshal(res)
	t.Log(string(marshal))
}

func TestPageQueryForStructWithLimitOffset(t *testing.T) {
	Init()
	res, err := sql2.PageQueryForStructWithLimitOffset[User](DB, &sql2.PageReq{Page: 1, Size: 100}, query, args...)
	if err != nil {
		t.Error(err)
	}
	marshal, _ := json.Marshal(res)
	t.Log(string(marshal))
}

func TestPageQueryForMapWithLimitOffset(t *testing.T) {
	Init()
	res, err := sql2.PageQueryForMapWithLimitOffset(DB, &sql2.PageReq{Page: 1, Size: 100}, query, args...)
	if err != nil {
		t.Error(err)
	}
	marshal, _ := json.Marshal(res)
	t.Log(string(marshal))
}

func TestPageQueryForStructWithRowNumber(t *testing.T) {
	Init()
	res, err := sql2.PageQueryForStructWithRowNumber[User](DB, &sql2.PageReq{Page: 1, Size: 100}, query, args...)
	if err != nil {
		t.Error(err)
	}
	marshal, _ := json.Marshal(res)
	t.Log(string(marshal))
}

func TestPageQueryForMapWithRowNumber(t *testing.T) {
	Init()
	res, err := sql2.PageQueryForMapWithRowNumber(DB, &sql2.PageReq{Page: 1, Size: 100}, query, args...)
	if err != nil {
		t.Error(err)
	}
	marshal, _ := json.Marshal(res)
	t.Log(string(marshal))
}

func TestPageQueryForStructWithFetchOffset(t *testing.T) {
	Init()
	res, err := sql2.PageQueryForStructWithFetchOffset[User](DB, &sql2.PageReq{Page: 1, Size: 100}, query, args...)
	if err != nil {
		t.Error(err)
	}
	marshal, _ := json.Marshal(res)
	t.Log(string(marshal))
}

func TestPageQueryForMapWithFetchOffset(t *testing.T) {
	Init()
	res, err := sql2.PageQueryForMapWithFetchOffset(DB, &sql2.PageReq{Page: 1, Size: 100}, query, args...)
	if err != nil {
		t.Error(err)
	}
	marshal, _ := json.Marshal(res)
	t.Log(string(marshal))
}

func TestPageQueryForStructWithDeclareCursor(t *testing.T) {
	Init()
	res, err := sql2.PageQueryForStructWithDeclareCursor[User](DB, &sql2.PageReq{Page: 1, Size: 100}, query, args...)
	if err != nil {
		t.Error(err)
	}
	marshal, _ := json.Marshal(res)
	t.Log(string(marshal))
}

func TestPageQueryForMapWithDeclareCursor(t *testing.T) {
	Init()
	res, err := sql2.PageQueryForMapWithDeclareCursor(DB, &sql2.PageReq{Page: 1, Size: 100}, query, args...)
	if err != nil {
		t.Error(err)
	}
	marshal, _ := json.Marshal(res)
	t.Log(string(marshal))
}
