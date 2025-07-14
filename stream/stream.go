package stream

import (
	"context"
	"cutil"
)

type Stream[T any] interface {
	// 源操作

	// Of 从指定元素创建一个流（有限流）
	// Of(values ...T) Stream[T]
	// OfChan 从指定通道创建一个流（无限流）
	// OfChan(chan T) Stream[T]

	// 中间操作

	// Peek 返回一个由该流的元素组成的流，并在从结果流中消耗元素时对每个元素执行提供的操作。
	Peek(action cutil.Consumer[T]) Stream[T]
	// Filter 返回一个流，该流由与给定条件函数匹配的元素组成。
	Filter(predicate cutil.Predicate[T]) Stream[T]
	// Skip 在丢弃流的前n个元素后，返回由该流的其余元素组成的流。如果此流包含的元素少于n个，则将返回一个空流。
	Skip(n int64) Stream[T]

	// 终止操作

	// ForEach 迭代流中的每个元素，按顺序执行提供的操作。
	ForEach(action cutil.Consumer[T])
	// ToArray 返回一个包含此流中所有元素的数组。
	ToArray() []T
	// Count 返回此流中的元素数。
	Count() int64
	// done 返回一个通道，用于接收流中的元素。内部方法
	done() chan T

	// 其他辅助函数

	// getCtx 返回上下文
	getCtx() context.Context

	// Todo 待实现

	// IntStream mapToInt(ToIntFunction mapper)
	// LongStream mapToLong(ToLongFunction mapper)
	// DoubleStream mapToDouble(ToDoubleFunction mapper)
	Distinct() Stream[T]
	// Sorted() Stream[T]
	Sorted(comparator cutil.Comparator[T]) Stream[T]
	Limit(maxSize int64) Stream[T]
	ForEachOrdered(action cutil.Consumer[T])
	// Object reduce(Object identity, BinaryOperator accumulator)
	// Optional reduce(BinaryOperator accumulator)
	// Optional min(Comparator comparator)
	// Optional max(Comparator comparator)
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

func Of[T any](ctx context.Context, values ...T) Stream[T] {
	p := newReferencePipeline[T](ctx)
	out := make(chan T, 1)
	p.out <- out
	go func() {
		defer close(out)
		for _, v := range values {
			select {
			case <-p.ctx.Done():
				return
			case out <- v:
			}
		}
	}()
	return p
}

func OfChan[T any](ctx context.Context, ins ...chan T) Stream[T] {
	p := newReferencePipeline[T](ctx)
	out := make(chan T, 1)
	p.out <- out
	go func() {
		defer close(out)
		for _, in := range ins {
			for v := range in {
				select {
				case <-p.ctx.Done():
					return
				case out <- v:
				}
			}
		}
	}()
	return p
}

// Map 返回一个流，该流由将给定函数应用于该流元素的结果组成。
func Map[T any, R any](stream Stream[T], mapper cutil.Function[T, R]) Stream[R] {
	in := stream.done()
	out := make(chan R, 1)
	go func() {
		defer close(out)
		for v := range in {
			result := mapper(v) // 应用映射函数
			out <- result
		}
	}()
	return OfChan(stream.getCtx(), out)
}

// FlatMap 返回一个流，该流由将此流的每个元素替换为映射流的内容的结果组成，该映射流是通过将提供的映射函数应用于每个元素而生成的。每个映射的流在其内容被放入该流后都会被关闭。（如果映射的流为null，则使用空流。）
func FlatMap[T any, R any, S Stream[R]](stream Stream[T], mapper cutil.Function[T, S]) Stream[R] {
	in := stream.done()
	out := make(chan R)
	go func() {
		defer close(out)
		for v := range in {
			mappedStream := mapper(v)
			for mappedValue := range mappedStream.done() {
				out <- mappedValue
			}
		}
	}()
	return OfChan(stream.getCtx(), out)
}

// Concat 返回一个流，该流由给定的多个流中的所有元素组成。
func Concat[T any](streams ...Stream[T]) Stream[T] {
	// Todo 返回空流
	if len(streams) == 0 {
		return nil
	}
	ctx := streams[0].getCtx()
	var chs []chan T
	for _, s := range streams {
		chs = append(chs, s.done())
	}
	return OfChan(ctx, chs...)
}

func Generate[T any](s cutil.Supplier[T]) Stream[T] {
	return nil
}
