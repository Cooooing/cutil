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
	Skip(n int64) Stream[T]
	// Limit 返回一个由该流的元素组成的流，该流被截断为长度不超过maxSize。如果此流包含的元素少于maxSize个，则会阻塞。
	Limit(maxSize int64) Stream[T]
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
	Count() (int64, error)

	// done 返回一个通道，用于接收流中的元素。内部方法
	done() chan T

	// 其他辅助函数

	// Parallel 设置并行流
	Parallel() Stream[T]
	// Sequential 设置顺序流
	Sequential() Stream[T]
	// getCtx 返回上下文
	getCtx() context.Context

	// Todo 待实现

	// IntStream mapToInt(ToIntFunction mapper)
	// LongStream mapToLong(ToLongFunction mapper)
	// DoubleStream mapToDouble(ToDoubleFunction mapper)
	// Sorted() FlagStream[T]

	// Object reduce(Object identity, BinaryOperator accumulator)
	// Optional reduce(BinaryOperator accumulator)
	// Optional min(Comparator comparator)
	// Optional max(Comparator comparator)

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
	// FlagSpliterator spliterator()
	// boolean isParallel()
	// BaseStream sequential()
	// BaseStream parallel()
	// BaseStream unordered()
	// BaseStream onClose(Runnable closeHandler)
	// void close()
}

// Of 从指定元素创建一个流（有限流）
func Of[T any](ctx context.Context, values ...T) Stream[T] {
	sourceFlags := IsSized | IsOrdered // 设置 SIZED 和 ORDERED 标志，因为 values 是固定大小的切片
	p := newReferencePipeline[T](ctx, sourceFlags)
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
	sourceFlags := 0 // 通道流通常大小未知，默认无 SIZED、ORDERED 等标志
	p := newReferencePipeline[T](ctx, sourceFlags)
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
	// Map 不改变 SIZED 或 ORDERED 标志
	opFlags := 0
	in := stream.done()
	p := newReferencePipeline[R](stream.getCtx(), opFlags)
	out := make(chan R, 1)
	p.out <- out
	p.sourceOrOpFlags = opFlags & OpMask
	p.combinedFlags = CombineOpFlags(opFlags, stream.(*ReferencePipeline[T]).combinedFlags)
	p.depth = stream.(*ReferencePipeline[T]).depth + 1
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
	// FlatMap 可能改变大小，设置为 NOT_SIZED
	opFlags := NotSized
	in := stream.done()
	p := newReferencePipeline[R](stream.getCtx(), opFlags)
	out := make(chan R, 1)
	p.out <- out
	p.sourceOrOpFlags = opFlags & OpMask
	p.combinedFlags = CombineOpFlags(opFlags, stream.(*ReferencePipeline[T]).combinedFlags)
	p.depth = stream.(*ReferencePipeline[T]).depth + 1
	go func() {
		defer close(out)
		for v := range in {
			mappedStream := mapper(v)
			if &mappedStream == nil {
				continue // 空流处理
			}
			for mappedValue := range mappedStream.done() {
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

// Concat 返回一个流，该流由给定的多个流中的所有元素组成。
func Concat[T any](ctx context.Context, streams ...Stream[T]) Stream[T] {
	if len(streams) == 0 {
		return Of[T](ctx) // 返回空流
	}
	// 计算合并后的标志
	combinedFlags := InitialOpsValue
	sourceFlags := 0 // 默认无标志
	for _, s := range streams {
		rp := s.(*ReferencePipeline[T])
		combinedFlags = CombineOpFlags(rp.sourceOrOpFlags, combinedFlags)
		// 如果任一输入流不是 SIZED，则结果流也不是 SIZED
		if !SIZED.IsKnown(rp.combinedFlags) {
			sourceFlags |= NotSized
		}
	}
	p := newReferencePipeline[T](ctx, sourceFlags)
	out := make(chan T, 1)
	p.out <- out
	go func() {
		defer close(out)
		for _, s := range streams {
			for v := range s.done() {
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
	// Generate 流是无限的，设置为 NOT_SIZED
	sourceFlags := NotSized
	p := newReferencePipeline[T](ctx, sourceFlags)
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

// Empty 创建一个空流
func Empty[T any](ctx context.Context) Stream[T] {
	return newReferencePipeline[T](ctx, IsSized|IsOrdered)
}

// estimateSize 估计流大小
// 标志使用：根据 SIZED 标志返回估计大小
func estimateSize[T any](p *ReferencePipeline[T]) int {
	if SIZED.IsKnown(p.combinedFlags) {
		// TODO: 实现实际大小估计，当前返回默认值
		return 16 // 假设默认容量
	}
	return 16 // 默认容量
}
