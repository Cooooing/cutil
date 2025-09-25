package stream

import (
	"context"
	"sync"

	"github.com/Cooooing/cutil/common"
)

type Stream[T any] interface {
	// 源操作（创建流）

	// 中间操作

	// Map 返回一个流，该流由该流中的每个元素通过指定的映射函数映射生成。转换后的元素类型不变。
	Map(mapper common.UnaryOperator[T]) Stream[T]
	// Peek 返回一个由该流的元素组成的流，并在从结果流中消耗元素时对每个元素执行提供的操作。
	Peek(action common.Consumer[T]) Stream[T]
	// Filter 返回一个流，该流由与给定条件函数匹配的元素组成。
	Filter(predicate common.Predicate[T]) Stream[T]
	// Skip 在丢弃流的前n个元素后，返回由该流的其余元素组成的流。如果此流包含的元素少于n个，则会阻塞。
	Skip(n int) Stream[T]
	// Limit 返回一个由该流的元素组成的流，该流被截断为长度不超过maxSize。如果此流包含的元素少于maxSize个，则会阻塞。
	Limit(maxSize int) Stream[T]
	// Distinct 返回一个去重后的流，元素顺序保持不变。
	Distinct() Stream[T]
	// Sorted 返回一个已排序的流，元素顺序保持不变。
	Sorted(comparator common.Comparator[T]) Stream[T]

	// 终止操作

	// ForEach 迭代流中的每个元素，按顺序执行提供的操作。
	ForEach(action common.Consumer[T]) error
	// ForEachOrdered 迭代流中的每个元素，按顺序执行提供的操作。
	ForEachOrdered(comparator common.Comparator[T], action common.Consumer[T]) error
	// AnyMatch 检查是否存在满足条件的元素
	AnyMatch(predicate common.Predicate[T]) (bool, error)
	// AllMatch 检查是否所有元素都满足条件
	AllMatch(predicate common.Predicate[T]) (bool, error)
	// NoneMatch 检查是否没有元素满足条件
	NoneMatch(predicate common.Predicate[T]) (bool, error)
	// ToArray 返回一个包含此流中所有元素的数组。
	ToArray() ([]T, error)
	// Count 返回此流中的元素数。
	Count() (int, error)
	// Min 返回此流中的最小元素。
	Min(comparator common.Comparator[T]) (T, error)
	// Max 返回此流中的最大元素。
	Max(comparator common.Comparator[T]) (T, error)
	// FindFirst 尝试返回此流中的第一个元素。
	FindFirst() (T, error)
	// FindAny 尝试返回此流中的任意元素。
	FindAny() (T, error)
	// Reduce 尝试将此流中元素归约为单个值。初始值为流中第一个元素。
	Reduce(accumulator common.BinaryOperator[T]) (T, error)
	// ReduceByDefault 尝试将此流中元素归约为单个值。给定初始值。
	ReduceByDefault(identity T, accumulator common.BinaryOperator[T]) (T, error)

	// Iterator 返回一个通道，用于迭代流中的元素。需要自行处理通道的生命周期。
	Iterator() chan T

	// 其他辅助函数

	// getCtx 返回上下文
	getCtx() context.Context
	// close 关闭流
	close(err error)
	// IsParallel 返回是否是并行流
	IsParallel() bool
	// GetParallelGoroutines 返回并行流并行线程数
	GetParallelGoroutines() int
	// Parallel 将流转为并行流，设置并发协程数。必须在创建流后紧接着调用。并行流不保证元素原始顺序。
	Parallel(n int) Stream[T]
}

// 创建流（源操作）

// Of 从指定元素创建一个流（有限流）
func Of[T any](ctx context.Context, values ...T) Stream[T] {
	if len(values) == 0 {
		return Empty[T](ctx) // 创建一个空流
	}
	p := newNoBlockStream[T](ctx)
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
	p := newNoBlockStream[T](ctx)
	out := make(chan T, 1)
	p.out <- out
	go func() {
		defer close(out)
		var wg sync.WaitGroup
		for _, in := range ins {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for v := range in {
					select {
					case <-p.ctx.Done():
						return
					case out <- v:
					}
				}
			}()
		}
		wg.Wait()
	}()
	return p
}

// Generate 返回一个无限流，由 Supplier 提供的元素组成
func Generate[T any](ctx context.Context, s common.Supplier[T]) Stream[T] {
	p := newNoBlockStream[T](ctx)
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
		return Empty[T](ctx) // 返回空流
	}
	p := newNoBlockStream[T](ctx)
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
	stream := newNoBlockStream[T](ctx)
	stream.out <- stream.closeChan()
	return stream
}

// 流操作，返回一个流或结果，但流中数据或结果类型发生改变。（解决go中方法不能增加泛型的问题）

// Map 返回一个流，该流由将给定函数应用于该流元素的结果组成。
func Map[T any, R any](stream Stream[T], mapper common.Function[T, R]) Stream[R] {
	in := stream.Iterator()
	p := newNoBlockStream[R](stream.getCtx())
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
func FlatMap[T any, R any](stream Stream[T], mapper common.Function[T, Stream[R]]) Stream[R] {
	in := stream.Iterator()
	p := newNoBlockStream[R](stream.getCtx())
	out := make(chan R, 1)
	p.out <- out
	go func() {
		defer close(out)
		for v := range in {
			mappedStream := mapper(v)
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

// Reduce 将流中的元素通过累积函数聚合为单一结果。
// identity: 归约的初始值，即使流为空也返回此值。
// accumulator: 累积函数，定义如何将流元素与中间结果合并。
// combiner: 合并函数，定义如何合并并行流中的子结果。
func Reduce[T any, R any](stream Stream[T], identity R, accumulator common.BiFunction[T, R], combiner common.BinaryOperator[R]) (R, error) {
	var err error
	result := identity
	if stream.IsParallel() {
		// 并行流逻辑
		subResults := make(chan R, stream.GetParallelGoroutines()) // 用于收集各子流的归约结果
		iterator := stream.Iterator()
		var currentWg sync.WaitGroup
		currentWg.Add(stream.GetParallelGoroutines())
		go func() {
			currentWg.Wait()
			stream.close(nil)
			close(subResults)
		}()
		for i := 0; i < stream.GetParallelGoroutines(); i++ {
			go func() {
				defer currentWg.Done()
				var localResult R
				for v := range iterator {
					localResult = accumulator(v, localResult)
				}
				subResults <- localResult
			}()
		}
		finalResult := identity
		for r := range subResults {
			finalResult = combiner(finalResult, r)
		}
		return finalResult, nil
	} else {
		// 顺序流逻辑
		err = stream.ForEach(func(item T) {
			result = accumulator(item, result)
		})
	}
	return result, err
}

// GroupBy 将流中的元素根据给定分类函数进行分组，并返回一个map，该map的键是分类函数的返回值，值是包含该键的元素组成的列表。
func GroupBy[T any, K comparable](stream Stream[T], classifier common.Function[T, K]) (map[K][]T, error) {
	var err error
	result := make(map[K][]T)
	err = stream.ForEach(func(item T) {
		key := classifier(item)
		result[key] = append(result[key], item)
	})
	return result, err
}
