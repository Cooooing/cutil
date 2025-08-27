package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Cooooing/cutil/common/logger"
	"github.com/Cooooing/cutil/sql/base"
)

type Executor[T any] struct {
	db      *sql.DB
	builder base.BaseBuilder
	debug   bool
}

func WithExecutor[T any](db *sql.DB, builder base.BaseBuilder) *Executor[T] {
	return &Executor[T]{
		db:      db,
		builder: builder,
		debug:   false,
	}
}

func (e *Executor[T]) Debug() *Executor[T] {
	e.debug = true
	return e
}

func (e *Executor[T]) Log() {
	s, args := e.builder.Build()
	logger.Info("\nSQL: %s\nArgs:%+v", s, args)
}

func (e *Executor[T]) log(s string, args ...any) {
	if e.debug {
		logger.Info("\nSQL: %s\nArgs:%+v", s, args)
	}
}

func (e *Executor[T]) Exec() (sql.Result, error) {
	s, args := e.builder.Build()
	e.log(s, args...)
	return e.db.Exec(s, args...)
}

func (e *Executor[T]) Raw() (*sql.Rows, error) {
	s, args := e.builder.Build()
	e.log(s, args...)
	return e.db.Query(s, args...)
}

func (e *Executor[T]) First() (*T, error) {
	if _, ok := e.builder.(base.SelectBuilder); ok {
		s, args := e.builder.Build()
		e.log(s, args...)
		s = fmt.Sprintf(`SELECT t.* FROM (%s) AS t LIMIT %d`, s, 1)
		t, err := base.Raw2Struct[T](e.db, s, args...)
		if err != nil {
			return nil, err
		}
		if len(t) == 0 {
			return nil, errors.New("no data")
		}
		return &t[0], err
	}
	return nil, base.ErrorExecutorNotSupportSelect
}

func (e *Executor[T]) List() ([]T, error) {
	if _, ok := e.builder.(base.SelectBuilder); ok {
		s, args := e.builder.Build()
		e.log(s, args...)
		return base.Raw2Struct[T](e.db, s, args...)
	}
	return nil, base.ErrorExecutorNotSupportSelect
}

func (e *Executor[T]) Count() (int, error) {
	if _, ok := e.builder.(base.SelectBuilder); ok {
		s, args := e.builder.Build()
		e.log(s, args...)
		return QueryCount(e.db, s, args...)
	}
	return 0, base.ErrorExecutorNotSupportSelect
}

func (e *Executor[T]) Page(page base.PageReqInterface) (base.PageRespInterface[T], error) {
	if _, ok := e.builder.(base.SelectBuilder); ok {
		s, args := e.builder.Build()
		e.log(s, args...)
		return PageQueryForStructWithLimitOffset[T](e.db, page, s, args...)
	}
	return nil, base.ErrorExecutorNotSupportSelect
}

func (e *Executor[T]) Delete() (int64, error) {
	if _, ok := e.builder.(base.DeleteBuilder); ok {
		exec, err := e.Exec()
		if err != nil {
			return 0, err
		}
		return exec.RowsAffected()
	}
	return 0, base.ErrorExecutorNotSupportDelete
}
