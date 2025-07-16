package stream

import (
	"context"
	"cutil"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"sync"
)

type Pipeline[T any] struct {
	ctx     context.Context
	out     chan chan T
	cancel  context.CancelFunc // 用于取消所有操作
	err     error
	errOnce sync.Once      // 确保只保留最初的错误
	wg      sync.WaitGroup // 跟踪所有活动协程

	linkedOrConsumed   bool // 标记流是否已被链接（添加新操作）或消耗（执行终端操作）。用于确保流的一次性使用，防止重复操作。
	parallelGoroutines int  // 并行协程数，等于1时为顺序流，大于1时为并行流。
	hasOperations      bool // 标记是否已有操作，用于判断是否是第一次添加操作
}

func newPipeline[T any](ctx context.Context) *Pipeline[T] {
	ctx, cancel := context.WithCancel(ctx)
	return &Pipeline[T]{
		ctx:                ctx,
		cancel:             cancel,
		out:                make(chan chan T, 1),
		linkedOrConsumed:   false,
		parallelGoroutines: 1,
	}
}

// 中间操作

func (p *Pipeline[T]) Map(action cutil.UnaryOperator[T]) Stream[T] {
	in, out := p.initOp()
	p.wg.Add(p.parallelGoroutines)
	for i := 0; i < p.parallelGoroutines; i++ {
		go func() {
			defer p.wg.Done()
			defer close(out)
			for v := range in {
				select {
				case <-p.ctx.Done():
					return
				case out <- action(v):
				}
			}
		}()
	}
	return p
}

func (p *Pipeline[T]) Peek(action cutil.Consumer[T]) Stream[T] {
	in, out := p.initOp()
	p.wg.Add(p.parallelGoroutines)
	for i := 0; i < p.parallelGoroutines; i++ {
		go func() {
			defer p.wg.Done()
			defer close(out)
			for v := range in {
				action(v)
				select {
				case <-p.ctx.Done():
					return
				case out <- v:
				}
			}
		}()
	}
	return p
}

func (p *Pipeline[T]) Filter(predicate cutil.Predicate[T]) Stream[T] {
	in, out := p.initOp()
	p.wg.Add(p.parallelGoroutines)
	for i := 0; i < p.parallelGoroutines; i++ {
		go func() {
			defer p.wg.Done()
			defer close(out)
			for v := range in {
				if predicate(v) {
					select {
					case <-p.ctx.Done():
						return
					case out <- v:
					}
				}
			}
		}()
	}
	return p
}

func (p *Pipeline[T]) Skip(n int) Stream[T] {
	in, out := p.initOp()
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		defer close(out)
		count := 0
		for v := range in {
			count++
			if count <= n {
				continue
			}
			select {
			case <-p.ctx.Done():
				return
			case out <- v:
			}
		}
	}()
	return p
}
func (p *Pipeline[T]) Limit(maxSize int) Stream[T] {
	in, out := p.initOp()
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		defer close(out)
		count := 0
		for v := range in {
			count++
			if count > maxSize {
				break
			}
			select {
			case <-p.ctx.Done():
				return
			case out <- v:
			}
		}
	}()
	return p
}

func (p *Pipeline[T]) Distinct() Stream[T] {
	in, out := p.initOp()
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		defer close(out)
		seen := make(map[any]struct{}) // 使用 map 去重
		for v := range in {
			if _, exists := seen[v]; !exists {
				seen[v] = struct{}{}
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

// Sorted 排序操作
func (p *Pipeline[T]) Sorted(comparator cutil.Comparator[T]) Stream[T] {
	in, out := p.initOp()
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		defer close(out)
		elements := make([]T, 0)
		for v := range in {
			elements = append(elements, v)
		}
		// 排序
		sort.Slice(elements, func(i, j int) bool {
			return comparator(elements[i], elements[j]) < 0
		})
		for _, v := range elements {
			select {
			case <-p.ctx.Done():
				return
			case out <- v:
			}
		}
	}()
	return p
}

// 终止操作

func (p *Pipeline[T]) ForEach(action cutil.Consumer[T]) error {
	in := p.initTerminalOp()
	if p.err != nil {
		return p.err
	}
	for v := range in {
		action(v)
	}
	p.close()
	return p.err
}

func (p *Pipeline[T]) ForEachOrdered(comparator cutil.Comparator[T], action cutil.Consumer[T]) error {
	p.Sorted(comparator)
	return p.ForEach(action)
}

func (p *Pipeline[T]) AnyMatch(predicate cutil.Predicate[T]) (bool, error) {
	in := p.initTerminalOp()
	if p.err != nil {
		return false, p.err
	}
	for v := range in {
		if predicate(v) {
			p.close()
			return true, p.err
		}
	}
	p.close()
	return false, p.err
}

func (p *Pipeline[T]) AllMatch(predicate cutil.Predicate[T]) (bool, error) {
	in := p.initTerminalOp()
	if p.err != nil {
		return false, p.err
	}
	for v := range in {
		if !predicate(v) {
			p.close()
			return false, p.err
		}
	}
	p.close()
	return true, p.err
}

func (p *Pipeline[T]) NoneMatch(predicate cutil.Predicate[T]) (bool, error) {
	in := p.initTerminalOp()
	if p.err != nil {
		return false, p.err
	}
	for v := range in {
		if predicate(v) {
			p.close()
			return false, p.err
		}
	}
	p.close()
	return true, p.err
}

func (p *Pipeline[T]) ToArray() ([]T, error) {
	in := p.initTerminalOp()
	if p.err != nil {
		return nil, p.err
	}
	array := make([]T, 0)
	for v := range in {
		array = append(array, v)
	}
	p.close()
	return array, p.err
}

func (p *Pipeline[T]) Count() (int, error) {
	in := p.initTerminalOp()
	if p.err != nil {
		return 0, p.err
	}
	count := 0
	for range in {
		count++
	}
	p.close()
	return count, p.err
}

func (p *Pipeline[T]) Min(comparator cutil.Comparator[T]) (T, error) {
	in := p.initTerminalOp()
	var zero T
	if p.err != nil {
		return zero, p.err
	}
	var m T
	for v := range in {
		select {
		case <-p.ctx.Done():
			p.close()
			return zero, p.err
		default:
			if comparator(v, m) < 0 {
				m = v
			}
		}
	}
	p.close()
	return m, p.err
}

func (p *Pipeline[T]) Max(comparator cutil.Comparator[T]) (T, error) {
	in := p.initTerminalOp()
	var zero T
	if p.err != nil {
		return zero, p.err
	}
	var m T
	for v := range in {
		select {
		case <-p.ctx.Done():
			p.close()
			return zero, p.err
		default:
			if comparator(v, m) > 0 {
				m = v
			}
		}
	}
	p.close()
	return m, p.err
}

func (p *Pipeline[T]) FindFirst() (T, error) {
	in := p.initTerminalOp()
	var m T
	if p.err != nil {
		return m, p.err
	}
	for v := range in {
		select {
		case <-p.ctx.Done():
			p.close()
			return m, p.err
		default:
			m = v
		}
	}
	p.close()
	return m, p.err
}
func (p *Pipeline[T]) FindAny() (T, error) {
	in := p.initTerminalOp()
	var m T
	if p.err != nil {
		return m, p.err
	}
	for v := range in {
		select {
		case <-p.ctx.Done():
			p.close()
			return m, p.err
		default:
			m = v
		}
	}
	p.close()
	return m, p.err
}

func (p *Pipeline[T]) Reduce(accumulator cutil.BinaryOperator[T]) (T, error) {
	in := p.initTerminalOp()
	var result T
	for v := range in {
		select {
		case <-p.ctx.Done():
			var zero T
			return zero, p.ctx.Err()
		default:
			result = accumulator(result, v)
		}
	}
	return result, nil
}

func (p *Pipeline[T]) ReduceByDefault(identity T, accumulator cutil.BinaryOperator[T]) (T, error) {
	in := p.initTerminalOp()
	result := identity
	for v := range in {
		select {
		case <-p.ctx.Done():
			var zero T
			return zero, p.ctx.Err()
		default:
			result = accumulator(result, v)
		}
	}
	return result, nil
}

func (p *Pipeline[T]) Iterator() chan T {
	p.linkedOrConsumed = true
	return <-p.out
}

func (p *Pipeline[T]) IsParallel() bool {
	return p.parallelGoroutines != 1
}

func (p *Pipeline[T]) Parallel(n int) Stream[T] {
	if n <= 0 {
		p.sendErr(errors.New(fmt.Sprintf("parallelism must be positive,but now is %s", strconv.Itoa(n))))
	}
	if p.hasOperations {
		p.sendErr(errors.New("parallel operation must be the first operation"))
	}
	p.hasOperations = true
	p.parallelGoroutines = n
	return p
}

// 其他辅助函数

func (p *Pipeline[T]) getCtx() context.Context {
	return p.ctx
}

// initOp 初始化中间操作
func (p *Pipeline[T]) initOp() (chan T, chan T) {
	if p.linkedOrConsumed {
		p.sendErr(errors.New("stream already operated upon or closed"))
	}
	p.hasOperations = true
	in, ok := <-p.out
	if !ok {
		return p.closeChan(), p.closeChan()
	}
	out := make(chan T, p.parallelGoroutines)
	p.out <- out
	return in, out
}

// initTerminalOp 初始化终端操作
func (p *Pipeline[T]) initTerminalOp() chan T {
	if p.linkedOrConsumed {
		p.sendErr(errors.New("stream already operated upon or closed"))
	}
	p.hasOperations = true
	in, ok := <-p.out
	if !ok {
		return p.closeChan()
	}
	return in
}

func (p *Pipeline[T]) closeChan() chan T {
	ch := make(chan T)
	close(ch)
	return ch
}

func (p *Pipeline[T]) close() {
	p.linkedOrConsumed = true
	p.cancel()  // 取消所有操作
	p.wg.Wait() // 等待所有协程结束

	close(p.out)
}

func (p *Pipeline[T]) sendErr(err error) {
	p.errOnce.Do(func() {
		p.err = err
		p.close()
	})
}
