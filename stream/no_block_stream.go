package stream

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"sync"

	"github.com/Cooooing/cutil/common"
)

type NoBlockStream[T any] struct {
	ctx       context.Context
	out       chan chan T
	cancel    context.CancelFunc // 用于取消所有操作
	err       error
	closeOnce sync.Once      // 确保只关闭流一次
	wg        sync.WaitGroup // 跟踪所有活动协程

	linkedOrConsumed   bool // 标记流是否已被链接（添加新操作）或消耗（执行终端操作）。用于确保流的一次性使用，防止重复操作。
	parallelGoroutines int  // 并行协程数，等于1时为顺序流，大于1时为并行流。
	hasOperations      bool // 标记是否已有操作，用于判断是否是第一次添加操作
}

func newNoBlockStream[T any](ctx context.Context) *NoBlockStream[T] {
	ctx, cancel := context.WithCancel(ctx)
	return &NoBlockStream[T]{
		ctx:                ctx,
		cancel:             cancel,
		out:                make(chan chan T, 1),
		linkedOrConsumed:   false,
		parallelGoroutines: 1,
	}
}

// 中间操作

func (p *NoBlockStream[T]) Map(action common.UnaryOperator[T]) Stream[T] {
	in, out := p.initOp()
	p.wg.Add(p.parallelGoroutines)
	var currentWg sync.WaitGroup
	currentWg.Add(p.parallelGoroutines)
	go func() {
		currentWg.Wait()
		close(out)
	}()
	for i := 0; i < p.parallelGoroutines; i++ {
		go func() {
			defer p.wg.Done()
			defer currentWg.Done()
			for {
				select {
				case <-p.ctx.Done():
					p.close(p.ctx.Err())
					return
				case v, ok := <-in:
					if !ok {
						return
					}
					out <- action(v)
				}
			}
		}()
	}
	return p
}

func (p *NoBlockStream[T]) Peek(action common.Consumer[T]) Stream[T] {
	in, out := p.initOp()
	p.wg.Add(p.parallelGoroutines)
	var currentWg sync.WaitGroup
	currentWg.Add(p.parallelGoroutines)
	go func() {
		currentWg.Wait()
		close(out)
	}()
	for i := 0; i < p.parallelGoroutines; i++ {
		go func() {
			defer p.wg.Done()
			defer currentWg.Done()
			for v := range in {
				action(v)
				select {
				case <-p.ctx.Done():
					p.close(p.ctx.Err())
					return
				case out <- v:
				}
			}
		}()
	}
	return p
}

func (p *NoBlockStream[T]) Filter(predicate common.Predicate[T]) Stream[T] {
	in, out := p.initOp()
	p.wg.Add(p.parallelGoroutines)
	var currentWg sync.WaitGroup
	currentWg.Add(p.parallelGoroutines)
	go func() {
		currentWg.Wait()
		close(out)
	}()
	for i := 0; i < p.parallelGoroutines; i++ {
		go func() {
			defer p.wg.Done()
			defer currentWg.Done()
			for v := range in {
				if predicate(v) {
					select {
					case <-p.ctx.Done():
						p.close(p.ctx.Err())
						return
					case out <- v:
					}
				}
			}
		}()
	}
	return p
}

func (p *NoBlockStream[T]) Skip(n int) Stream[T] {
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
				p.close(p.ctx.Err())
				return
			case out <- v:
			}
		}
	}()
	return p
}
func (p *NoBlockStream[T]) Limit(maxSize int) Stream[T] {
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
				p.close(p.ctx.Err())
				return
			case out <- v:
			}
		}
	}()
	return p
}

func (p *NoBlockStream[T]) Distinct() Stream[T] {
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
					p.close(p.ctx.Err())
					return
				case out <- v:
				}
			}
		}
	}()
	return p
}

// Sorted 排序操作
func (p *NoBlockStream[T]) Sorted(comparator common.Comparator[T]) Stream[T] {
	in, out := p.initOp()
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		defer close(out)
		elements := make([]T, 0)
		over := false
		for {
			select {
			case <-p.ctx.Done():
				p.close(p.ctx.Err())
				return
			case v, ok := <-in:
				if !ok {
					over = true
					break
				}
				elements = append(elements, v)
			}
			if over {
				break
			}
		}
		// 排序
		sort.Slice(elements, func(i, j int) bool {
			return comparator(elements[i], elements[j]) < 0
		})
		for _, v := range elements {
			select {
			case <-p.ctx.Done():
				p.close(p.ctx.Err())
				return
			case out <- v:
			}
		}
	}()
	return p
}

// 终止操作

func (p *NoBlockStream[T]) ForEach(action common.Consumer[T]) error {
	in := p.initTerminalOp()
	if p.err != nil {
		return p.err
	}
	for {
		select {
		case <-p.ctx.Done():
			p.close(p.ctx.Err())
			return p.ctx.Err()
		case v, ok := <-in:
			if !ok {
				p.close(p.err)
				return p.err
			}
			action(v)
		}
	}
}

func (p *NoBlockStream[T]) ForEachOrdered(comparator common.Comparator[T], action common.Consumer[T]) error {
	p.Sorted(comparator)
	return p.ForEach(action)
}

func (p *NoBlockStream[T]) AnyMatch(predicate common.Predicate[T]) (bool, error) {
	in := p.initTerminalOp()
	if p.err != nil {
		return false, p.err
	}
	for {
		select {
		case <-p.ctx.Done():
			p.close(p.ctx.Err())
			return false, p.ctx.Err()
		case v, ok := <-in:
			if !ok {
				p.close(p.err)
				return false, p.err
			}
			return predicate(v), p.err
		}
	}
}

func (p *NoBlockStream[T]) AllMatch(predicate common.Predicate[T]) (bool, error) {
	in := p.initTerminalOp()
	if p.err != nil {
		return false, p.err
	}
	for {
		select {
		case <-p.ctx.Done():
			p.close(p.ctx.Err())
			return true, p.ctx.Err()
		case v, ok := <-in:
			if !ok {
				p.close(p.err)
				return true, p.err
			}
			if !predicate(v) {
				p.close(p.err)
				return false, p.err
			}
		}
	}
}

func (p *NoBlockStream[T]) NoneMatch(predicate common.Predicate[T]) (bool, error) {
	in := p.initTerminalOp()
	if p.err != nil {
		return false, p.err
	}
	for {
		select {
		case <-p.ctx.Done():
			p.close(p.ctx.Err())
			return true, p.ctx.Err()
		case v, ok := <-in:
			if !ok {
				p.close(p.err)
				return true, p.err
			}
			if predicate(v) {
				p.close(p.err)
				return false, p.err
			}
		}
	}
}

func (p *NoBlockStream[T]) ToArray() ([]T, error) {
	in := p.initTerminalOp()
	if p.err != nil {
		return nil, p.err
	}
	array := make([]T, 0)
	for {
		select {
		case <-p.ctx.Done():
			p.close(p.ctx.Err())
			return array, p.ctx.Err()
		case v, ok := <-in:
			if !ok {
				p.close(p.err)
				return array, p.err
			}
			array = append(array, v)
		}
	}
}

func (p *NoBlockStream[T]) Count() (int, error) {
	in := p.initTerminalOp()
	if p.err != nil {
		return 0, p.err
	}
	count := 0
	for {
		select {
		case <-p.ctx.Done():
			p.close(p.ctx.Err())
			return count, p.ctx.Err()
		case _, ok := <-in:
			if !ok {
				p.close(p.err)
				return count, p.err
			}
			count++
		}
	}
}

func (p *NoBlockStream[T]) Min(comparator common.Comparator[T]) (T, error) {
	in := p.initTerminalOp()
	var zero T
	if p.err != nil {
		return zero, p.err
	}
	var m T
	for {
		select {
		case <-p.ctx.Done():
			p.close(p.ctx.Err())
			return zero, p.ctx.Err()
		case v, ok := <-in:
			if !ok {
				p.close(p.err)
				return m, p.err
			}
			if comparator(v, m) < 0 {
				m = v
			}
		}
	}
}

func (p *NoBlockStream[T]) Max(comparator common.Comparator[T]) (T, error) {
	in := p.initTerminalOp()
	var zero T
	if p.err != nil {
		return zero, p.err
	}
	var m T
	for {
		select {
		case <-p.ctx.Done():
			p.close(p.ctx.Err())
			return zero, p.ctx.Err()
		case v, ok := <-in:
			if !ok {
				p.close(p.err)
				return m, p.err
			}
			if comparator(v, m) > 0 {
				m = v
			}
		}
	}
}

func (p *NoBlockStream[T]) FindFirst() (T, error) {
	in := p.initTerminalOp()
	var m T
	if p.err != nil {
		return m, p.err
	}
	for {
		select {
		case <-p.ctx.Done():
			p.close(p.ctx.Err())
			return m, p.ctx.Err()
		case v, ok := <-in:
			if !ok {
				p.close(p.err)
				return m, p.err
			}
			return v, p.err
		}
	}
}
func (p *NoBlockStream[T]) FindAny() (T, error) {
	in := p.initTerminalOp()
	var m T
	if p.err != nil {
		return m, p.err
	}
	for {
		select {
		case <-p.ctx.Done():
			p.close(p.ctx.Err())
			return m, p.ctx.Err()
		case v, ok := <-in:
			if !ok {
				p.close(p.err)
				return m, p.err
			}
			return v, p.err
		}
	}
}

func (p *NoBlockStream[T]) Reduce(accumulator common.BinaryOperator[T]) (T, error) {
	in := p.initTerminalOp()
	var result T
	for {
		select {
		case <-p.ctx.Done():
			var zero T
			return zero, p.ctx.Err()
		case v, ok := <-in:
			if !ok {
				p.close(p.err)
				return result, p.err
			}
			result = accumulator(result, v)
		}
	}
}

func (p *NoBlockStream[T]) ReduceByDefault(identity T, accumulator common.BinaryOperator[T]) (T, error) {
	in := p.initTerminalOp()
	result := identity
	for {
		select {
		case <-p.ctx.Done():
			var zero T
			return zero, p.ctx.Err()
		case v, ok := <-in:
			if !ok {
				p.close(p.err)
				return result, p.err
			}
			result = accumulator(result, v)
		}
	}
}

func (p *NoBlockStream[T]) Iterator() chan T {
	p.linkedOrConsumed = true
	return <-p.out
}

func (p *NoBlockStream[T]) IsParallel() bool {
	return p.parallelGoroutines != 1
}

func (p *NoBlockStream[T]) GetParallelGoroutines() int {
	return p.parallelGoroutines
}

func (p *NoBlockStream[T]) Parallel(n int) Stream[T] {
	if n <= 0 {
		p.close(errors.New(fmt.Sprintf("parallelism must be positive,but now is %s", strconv.Itoa(n))))
	}
	if p.hasOperations {
		p.close(errors.New("parallel operation must be the first operation"))
	}
	p.hasOperations = true
	p.parallelGoroutines = n
	return p
}

// 其他辅助函数

func (p *NoBlockStream[T]) getCtx() context.Context {
	return p.ctx
}

// initOp 初始化中间操作
func (p *NoBlockStream[T]) initOp() (chan T, chan T) {
	if p.linkedOrConsumed {
		p.close(errors.New("stream already operated upon or closed"))
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
func (p *NoBlockStream[T]) initTerminalOp() chan T {
	if p.linkedOrConsumed {
		p.close(errors.New("stream already operated upon or closed"))
	}
	p.hasOperations = true
	in, ok := <-p.out
	if !ok {
		return p.closeChan()
	}
	return in
}

func (p *NoBlockStream[T]) closeChan() chan T {
	ch := make(chan T)
	close(ch)
	return ch
}

func (p *NoBlockStream[T]) close(err error) {
	go func() {
		p.closeOnce.Do(func() {
			if err != nil {
				p.err = err
				p.linkedOrConsumed = true
				p.cancel()  // 取消所有操作
				p.wg.Wait() // 等待所有协程结束
				close(p.out)
			}
		})
	}()
}
