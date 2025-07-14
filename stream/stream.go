package stream

import "cutil"

type Stream[T any] interface {
	Of(values ...T) Stream[T]

	Filter(predicate cutil.Predicate[T]) Stream[T]
	// IntStream mapToInt(ToIntFunction mapper)
	// LongStream mapToLong(ToLongFunction mapper)
	// DoubleStream mapToDouble(ToDoubleFunction mapper)
	Distinct() Stream[T]
	// Sorted() Stream[T]
	Sorted(comparator cutil.Comparator[T]) Stream[T]
	Peek(action cutil.Consumer[T]) Stream[T]
	Limit(maxSize int64) Stream[T]
	Skip(n int64) Stream[T]
	ForEach(action cutil.Consumer[T])
	ForEachOrdered(action cutil.Consumer[T])
	ToArray() []T
	// Object reduce(Object identity, BinaryOperator accumulator)
	// Optional reduce(BinaryOperator accumulator)
	// Optional min(Comparator comparator)
	// Optional max(Comparator comparator)
	Count() int64
	AnyMatch(predicate cutil.Predicate[T]) bool
	AllMatch(predicate cutil.Predicate[T]) bool
	NoneMatch(predicate cutil.Predicate[T]) bool
	// Optional findFirst()
	// Optional findAny()
	// Object collect(Collector collector)
	// Object collect(Supplier supplier, BiConsumer accumulator, BiConsumer combiner)
	// Object reduce(Object identity, BiFunction accumulator, BinaryOperator combiner)
	// Object[] toArray(IntFunction generator)
	// DoubleStream flatMapToDouble(Function mapper)
	// LongStream flatMapToLong(Function mapper)
	// IntStream flatMapToInt(Function mapper)
	// Iterator iterator()
	// Spliterator spliterator()
	// boolean isParallel()
	// BaseStream sequential()
	// BaseStream parallel()
	// BaseStream unordered()
	// BaseStream onClose(Runnable closeHandler)
	// void close()
}

func FlatMap[T any, R any](mapper cutil.Function[T, R]) Stream[R] {
	return nil
}

func Map[T any, R any](mapper cutil.Function[T, R]) Stream[R] {
	return nil
}

func Of[T any](values ...T) Stream[T] {
	return nil
}

func Concat[T any](a Stream[T], b Stream[T]) Stream[T] {
	return nil
}

func Generate[T any](s cutil.Supplier[T]) Stream[T] {
	return nil
}
