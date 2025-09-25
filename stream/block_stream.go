package stream

import (
	"context"
	"errors"
	"sort"

	"github.com/Cooooing/cutil/common"
)

// BlockingStream 阻塞流实现，内部直接存储在 slice 中。
type BlockingStream[T any] struct {
	ctx      context.Context
	elements []T
	parallel bool
	workers  int
}

func newBlockingStream[T any](ctx context.Context, elems []T) *BlockingStream[T] {
	return &BlockingStream[T]{ctx: ctx, elements: elems, parallel: false, workers: 1}
}

// OfBlocking 从元素创建阻塞流
func OfBlocking[T any](ctx context.Context, values ...T) Stream[T] {
	return newBlockingStream(ctx, values)
}

// --------------------- 中间操作 ---------------------

func (s *BlockingStream[T]) Map(mapper common.UnaryOperator[T]) Stream[T] {
	newElems := make([]T, 0, len(s.elements))
	for _, v := range s.elements {
		newElems = append(newElems, mapper(v))
	}
	return newBlockingStream(s.ctx, newElems)
}

func (s *BlockingStream[T]) Peek(action common.Consumer[T]) Stream[T] {
	for _, v := range s.elements {
		action(v)
	}
	return newBlockingStream(s.ctx, s.elements)
}

func (s *BlockingStream[T]) Filter(predicate common.Predicate[T]) Stream[T] {
	newElems := make([]T, 0)
	for _, v := range s.elements {
		if predicate(v) {
			newElems = append(newElems, v)
		}
	}
	return newBlockingStream(s.ctx, newElems)
}

func (s *BlockingStream[T]) Skip(n int) Stream[T] {
	if n >= len(s.elements) {
		return newBlockingStream(s.ctx, []T{})
	}
	return newBlockingStream(s.ctx, s.elements[n:])
}

func (s *BlockingStream[T]) Limit(maxSize int) Stream[T] {
	if maxSize >= len(s.elements) {
		return newBlockingStream(s.ctx, s.elements)
	}
	return newBlockingStream(s.ctx, s.elements[:maxSize])
}

func (s *BlockingStream[T]) Distinct() Stream[T] {
	seen := make(map[any]struct{})
	newElems := make([]T, 0, len(s.elements))
	for _, v := range s.elements {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			newElems = append(newElems, v)
		}
	}
	return newBlockingStream(s.ctx, newElems)
}

func (s *BlockingStream[T]) Sorted(comparator common.Comparator[T]) Stream[T] {
	newElems := append([]T(nil), s.elements...)
	sort.Slice(newElems, func(i, j int) bool {
		return comparator(newElems[i], newElems[j]) < 0
	})
	return newBlockingStream(s.ctx, newElems)
}

// --------------------- 终止操作 ---------------------

func (s *BlockingStream[T]) ForEach(action common.Consumer[T]) error {
	for _, v := range s.elements {
		select {
		case <-s.ctx.Done():
			return s.ctx.Err()
		default:
			action(v)
		}
	}
	return nil
}

func (s *BlockingStream[T]) ForEachOrdered(comparator common.Comparator[T], action common.Consumer[T]) error {
	newElems := append([]T(nil), s.elements...)
	sort.Slice(newElems, func(i, j int) bool {
		return comparator(newElems[i], newElems[j]) < 0
	})
	for _, v := range newElems {
		select {
		case <-s.ctx.Done():
			return s.ctx.Err()
		default:
			action(v)
		}
	}
	return nil
}

func (s *BlockingStream[T]) AnyMatch(predicate common.Predicate[T]) (bool, error) {
	for _, v := range s.elements {
		if predicate(v) {
			return true, nil
		}
	}
	return false, nil
}

func (s *BlockingStream[T]) AllMatch(predicate common.Predicate[T]) (bool, error) {
	for _, v := range s.elements {
		if !predicate(v) {
			return false, nil
		}
	}
	return true, nil
}

func (s *BlockingStream[T]) NoneMatch(predicate common.Predicate[T]) (bool, error) {
	for _, v := range s.elements {
		if predicate(v) {
			return false, nil
		}
	}
	return true, nil
}

func (s *BlockingStream[T]) ToArray() ([]T, error) {
	return s.elements, nil
}

func (s *BlockingStream[T]) Count() (int, error) {
	return len(s.elements), nil
}

func (s *BlockingStream[T]) Min(comparator common.Comparator[T]) (T, error) {
	if len(s.elements) == 0 {
		var zero T
		return zero, errors.New("stream is empty")
	}
	min := s.elements[0]
	for _, v := range s.elements[1:] {
		if comparator(v, min) < 0 {
			min = v
		}
	}
	return min, nil
}

func (s *BlockingStream[T]) Max(comparator common.Comparator[T]) (T, error) {
	if len(s.elements) == 0 {
		var zero T
		return zero, errors.New("stream is empty")
	}
	max := s.elements[0]
	for _, v := range s.elements[1:] {
		if comparator(v, max) > 0 {
			max = v
		}
	}
	return max, nil
}

func (s *BlockingStream[T]) FindFirst() (T, error) {
	if len(s.elements) == 0 {
		var zero T
		return zero, errors.New("stream is empty")
	}
	return s.elements[0], nil
}

func (s *BlockingStream[T]) FindAny() (T, error) {
	if len(s.elements) == 0 {
		var zero T
		return zero, errors.New("stream is empty")
	}
	// 阻塞流里 FindAny 就返回第一个
	return s.elements[0], nil
}

func (s *BlockingStream[T]) Reduce(accumulator common.BinaryOperator[T]) (T, error) {
	if len(s.elements) == 0 {
		var zero T
		return zero, errors.New("stream is empty")
	}
	result := s.elements[0]
	for _, v := range s.elements[1:] {
		result = accumulator(result, v)
	}
	return result, nil
}

func (s *BlockingStream[T]) ReduceByDefault(identity T, accumulator common.BinaryOperator[T]) (T, error) {
	result := identity
	for _, v := range s.elements {
		result = accumulator(result, v)
	}
	return result, nil
}

// --------------------- 辅助函数 ---------------------

func (s *BlockingStream[T]) Iterator() chan T {
	ch := make(chan T)
	go func() {
		defer close(ch)
		for _, v := range s.elements {
			ch <- v
		}
	}()
	return ch
}

func (s *BlockingStream[T]) getCtx() context.Context { return s.ctx }
func (s *BlockingStream[T]) close(err error)         {}
func (s *BlockingStream[T]) IsParallel() bool        { return s.parallel }
func (s *BlockingStream[T]) GetParallelGoroutines() int {
	return s.workers
}
func (s *BlockingStream[T]) Parallel(n int) Stream[T] {
	// 阻塞流即使设置并行，也还是顺序执行
	newS := *s
	newS.parallel = true
	newS.workers = n
	return &newS
}
