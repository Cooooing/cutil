package cutil

type (
	Function[T any, R any] func(T) R
	Consumer[T any]        func(T)
	Supplier[T any]        func() T
	Predicate[T any]       func(T) bool
	Comparator[T any]      func(T, T) int
)
