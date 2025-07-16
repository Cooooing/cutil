package cutil

type (
	Function[T any, R any] func(T) R
	Consumer[T any]        func(T)
	Supplier[T any]        func() T
	Predicate[T any]       func(T) bool
	Comparator[T any]      func(T, T) int
	UnaryOperator[T any]   func(T) T

	BinaryOperator[T any]    func(T, T) T
	BiFunction[T any, R any] func(T, R) R

	ToIntFunction[T any]        func(T) int
	ToInt8Function[T any]       func(T) int8
	ToInt16Function[T any]      func(T) int16
	ToInt32Function[T any]      func(T) int32
	ToInt64Function[T any]      func(T) int64
	ToUintFunction[T any]       func(T) uint
	ToUint8Function[T any]      func(T) uint8
	ToUint16Function[T any]     func(T) uint16
	ToUint32Function[T any]     func(T) uint32
	ToUint64Function[T any]     func(T) uint64
	ToFloat32Function[T any]    func(T) float32
	ToFloat64Function[T any]    func(T) float64
	ToComplex64Function[T any]  func(T) complex64
	ToComplex128Function[T any] func(T) complex128
	ToStringFunction[T any]     func(T) string
	ToBytesFunction[T any]      func(T) []byte
	ToRuneFunction[T any]       func(T) rune
)
