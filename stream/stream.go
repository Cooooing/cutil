package stream

import (
	"context"
	"cutil"
)

type Stream[T any] interface {
	// 源操作（创建流）

	// 中间操作

	// Map 返回一个流，该流由该流中的每个元素通过指定的映射函数映射生成。转换后的元素类型不变。
	Map(mapper cutil.UnaryOperator[T]) Stream[T]
	// Peek 返回一个由该流的元素组成的流，并在从结果流中消耗元素时对每个元素执行提供的操作。
	Peek(action cutil.Consumer[T]) Stream[T]
	// Filter 返回一个流，该流由与给定条件函数匹配的元素组成。
	Filter(predicate cutil.Predicate[T]) Stream[T]
	// Skip 在丢弃流的前n个元素后，返回由该流的其余元素组成的流。如果此流包含的元素少于n个，则会阻塞。
	Skip(n int) Stream[T]
	// Limit 返回一个由该流的元素组成的流，该流被截断为长度不超过maxSize。如果此流包含的元素少于maxSize个，则会阻塞。
	Limit(maxSize int) Stream[T]
	// Distinct 返回一个去重后的流，元素顺序保持不变。
	Distinct() Stream[T]
	// Sorted 返回一个已排序的流，元素顺序保持不变。
	Sorted(comparator cutil.Comparator[T]) Stream[T]

	// 终止操作

	// ForEach 迭代流中的每个元素，按顺序执行提供的操作。
	ForEach(action cutil.Consumer[T]) error
	// ForEachOrdered 迭代流中的每个元素，按顺序执行提供的操作。
	ForEachOrdered(comparator cutil.Comparator[T], action cutil.Consumer[T]) error
	// AnyMatch 检查是否存在满足条件的元素
	AnyMatch(predicate cutil.Predicate[T]) (bool, error)
	// AllMatch 检查是否所有元素都满足条件
	AllMatch(predicate cutil.Predicate[T]) (bool, error)
	// NoneMatch 检查是否没有元素满足条件
	NoneMatch(predicate cutil.Predicate[T]) (bool, error)
	// ToArray 返回一个包含此流中所有元素的数组。
	ToArray() ([]T, error)
	// Count 返回此流中的元素数。
	Count() (int, error)
	// Min 返回此流中的最小元素。
	Min(comparator cutil.Comparator[T]) (T, error)
	// Max 返回此流中的最大元素。
	Max(comparator cutil.Comparator[T]) (T, error)
	// FindFirst 尝试返回此流中的第一个元素。
	FindFirst() (T, error)
	// FindAny 尝试返回此流中的任意元素。
	FindAny() (T, error)
	// Reduce 尝试将此流中元素归约为单个值。初始值为流中第一个元素。
	Reduce(accumulator cutil.BinaryOperator[T]) (T, error)
	// ReduceByDefault 尝试将此流中元素归约为单个值。给定初始值。
	ReduceByDefault(identity T, accumulator cutil.BinaryOperator[T]) (T, error)

	// Iterator 返回一个通道，用于迭代流中的元素。需要自行处理通道的生命周期。
	Iterator() chan T

	// 其他辅助函数

	// getCtx 返回上下文
	getCtx() context.Context
	// close 关闭流
	close()
	// IsParallel 返回是否是并行流
	IsParallel() bool
	// Parallel 将流转为并行流，设置并发协程数。
	Parallel(n int) Stream[T]
}

// 创建流（源操作）

// Of 从指定元素创建一个流（有限流）
func Of[T any](ctx context.Context, values ...T) Stream[T] {
	p := newPipeline[T](ctx)
	out := make(chan T, len(values)) // 使用缓冲通道优化性能
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

// OfChan 从指定通道创建一个流（无限流）
func OfChan[T any](ctx context.Context, ins ...chan T) Stream[T] {
	p := newPipeline[T](ctx)
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

// Generate 返回一个无限流，由 Supplier 提供的元素组成
func Generate[T any](ctx context.Context, s cutil.Supplier[T]) Stream[T] {
	p := newPipeline[T](ctx)
	out := make(chan T, 1)
	p.out <- out
	go func() {
		defer close(out)
		for {
			select {
			case <-p.ctx.Done():
				return
			case out <- s():
			}
		}
	}()
	return p
}

// Concat 返回一个流，该流由给定的多个流中的所有元素组成。
func Concat[T any](ctx context.Context, streams ...Stream[T]) Stream[T] {
	if len(streams) == 0 {
		return Of[T](ctx) // 返回空流
	}
	p := newPipeline[T](ctx)
	out := make(chan T, 1)
	p.out <- out
	go func() {
		defer close(out)
		for _, s := range streams {
			for v := range s.Iterator() {
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

// Empty 创建一个空流
func Empty[T any](ctx context.Context) Stream[T] {
	return newPipeline[T](ctx)
}

// 流操作，返回一个流或结果，但流中数据或结果类型发生改变。（解决go中方法不能增加泛型的问题）

// Map 返回一个流，该流由将给定函数应用于该流元素的结果组成。
func Map[T any, R any](stream Stream[T], mapper cutil.Function[T, R]) Stream[R] {
	in := stream.Iterator()
	p := newPipeline[R](stream.getCtx())
	out := make(chan R, 1)
	p.out <- out
	go func() {
		defer close(out)
		for v := range in {
			select {
			case <-p.ctx.Done():
				return
			case out <- mapper(v):
			}
		}
	}()
	return p
}

// FlatMap 返回一个流，该流由将此流的每个元素替换为映射流的内容的结果组成，该映射流是通过将提供的映射函数应用于每个元素而生成的。每个映射的流在其内容被放入该流后都会被关闭。（如果映射的流为null，则使用空流。）
func FlatMap[T any, R any, S Stream[R]](stream Stream[T], mapper cutil.Function[T, S]) Stream[R] {
	in := stream.Iterator()
	p := newPipeline[R](stream.getCtx())
	out := make(chan R, 1)
	p.out <- out
	go func() {
		defer close(out)
		for v := range in {
			mappedStream := mapper(v)
			if &mappedStream == nil {
				continue // 空流处理
			}
			for mappedValue := range mappedStream.Iterator() {
				select {
				case <-p.ctx.Done():
					return
				case out <- mappedValue:
				}
			}
		}
	}()
	return p
}
