package base

import "errors"

var (
	ErrorExecutorNotSupportSelect = errors.New("this executor not support select")
	ErrorExecutorNotSupportUpdate = errors.New("this executor not support update")
	ErrorExecutorNotSupportDelete = errors.New("this executor not support delete")
	ErrorExecutorNotSupportInsert = errors.New("this executor not support insert")
)

const ()
