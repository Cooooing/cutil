package stream

import (
	"context"
	"cutil"
	"errors"
	"sort"
	"sync"
)

type ReferencePipeline[T any] struct {
	ctx     context.Context
	out     chan chan T
	cancel  context.CancelFunc // 用于取消所有操作
	err     error
	errOnce sync.Once      // 确保错误只关闭一次
	wg      sync.WaitGroup // 跟踪所有活动协程

	// 从源到当前阶段的所有标志的组合，反映整个管道的特性。
	combinedFlags int
	// 当前阶段的标志（源标志或操作标志），受 StreamOpFlag 约束。
	sourceOrOpFlags int
	// 表示流是否为并行流，仅对源阶段有效。用于控制流的执行模式（顺序或并行）。
	parallel bool
	// 标记流是否已被链接（添加新操作）或消耗（执行终端操作）。用于确保流的一次性使用，防止重复操作。
	linkedOrConsumed bool
	// 表示当前阶段到源阶段的中间操作数量（顺序流）或到上一个有状态操作的距离（并行流）。用于支持并行流的阶段切分和优化。
	depth int
}

func newReferencePipeline[T any](ctx context.Context, sourceFlags int) *ReferencePipeline[T] {
	ctx, cancel := context.WithCancel(ctx)
	return &ReferencePipeline[T]{
		ctx:             ctx,
		cancel:          cancel,
		out:             make(chan chan T, 1),
		combinedFlags:   CombineOpFlags(sourceFlags, InitialOpsValue),
		sourceOrOpFlags: sourceFlags & StreamMask,
	}
}

// 中间操作

func (p *ReferencePipeline[T]) Map(action cutil.UnaryOperator[T]) Stream[T] {
	in, out := p.initOp(0)
	p.wg.Add(1)
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
	return p
}

func (p *ReferencePipeline[T]) Peek(action cutil.Consumer[T]) Stream[T] {
	in, out := p.initOp(0)
	p.wg.Add(1)
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
	return p
}

func (p *ReferencePipeline[T]) Filter(predicate cutil.Predicate[T]) Stream[T] {
	in, out := p.initOp(NotSized) // 过滤可能改变大小
	p.wg.Add(1)
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
	return p
}

func (p *ReferencePipeline[T]) Skip(n int64) Stream[T] {
	in, out := p.initOp(IsSizeAdjusting) // 跳过调整大小
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		defer close(out)
		count := int64(0)
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
func (p *ReferencePipeline[T]) Limit(maxSize int64) Stream[T] {
	in, out := p.initOp(IsSizeAdjusting | IsShortCircuit) // 限制调整大小并支持短路
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		defer close(out)
		count := int64(0)
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

func (p *ReferencePipeline[T]) Distinct() Stream[T] {
	if DISTINCT.IsKnown(p.combinedFlags) {
		return p // 流已去重，直接返回
	}
	in, out := p.initOp(IsDistinct) // 设置 DISTINCT 标志
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
func (p *ReferencePipeline[T]) Sorted(comparator cutil.Comparator[T]) Stream[T] {
	in, out := p.initOp(IsSorted | NotOrdered | NotSized) // 设置 SORTED，清除 ORDERED 和 SIZED
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

func (p *ReferencePipeline[T]) ForEach(action cutil.Consumer[T]) error {
	in := p.initTerminalOp(0)
	if p.err != nil {
		return p.err
	}
	for v := range in {
		action(v)
	}
	p.close()
	return nil
}

func (p *ReferencePipeline[T]) ForEachOrdered(comparator cutil.Comparator[T], action cutil.Consumer[T]) error {
	if !SORTED.IsKnown(p.combinedFlags) {
		// 如果流无序，先排序
		p.Sorted(comparator)
	}
	return p.ForEach(action)
}

func (p *ReferencePipeline[T]) AnyMatch(predicate cutil.Predicate[T]) (bool, error) {
	in := p.initTerminalOp(IsShortCircuit)
	if p.err != nil {
		return false, p.err
	}
	for v := range in {
		if predicate(v) {
			p.close()
			return true, nil
		}
	}
	p.close()
	return false, nil
}

func (p *ReferencePipeline[T]) AllMatch(predicate cutil.Predicate[T]) (bool, error) {
	in := p.initTerminalOp(IsShortCircuit)
	if p.err != nil {
		return false, p.err
	}
	for v := range in {
		if !predicate(v) {
			p.close()
			return false, nil
		}
	}
	p.close()
	return true, nil
}

func (p *ReferencePipeline[T]) NoneMatch(predicate cutil.Predicate[T]) (bool, error) {
	in := p.initTerminalOp(IsShortCircuit)
	if p.err != nil {
		return false, p.err
	}
	for v := range in {
		if predicate(v) {
			p.close()
			return false, nil
		}
	}
	p.close()
	return true, nil
}

func (p *ReferencePipeline[T]) ToArray() ([]T, error) {
	in := p.initTerminalOp(0)
	if p.err != nil {
		return nil, p.err
	}
	array := make([]T, 0)
	for v := range in {
		array = append(array, v)
	}
	p.close()
	return array, nil
}

func (p *ReferencePipeline[T]) Count() (int64, error) {
	in := p.initTerminalOp(0)
	if p.err != nil {
		return 0, p.err
	}
	count := int64(0)
	for range in {
		count++
	}
	p.close()
	return count, nil
}

func (p *ReferencePipeline[T]) done() chan T {
	p.linkedOrConsumed = true
	return <-p.out
}

// 其他辅助函数

func (p *ReferencePipeline[T]) Parallel() Stream[T] {
	p.parallel = true
	return p
}

func (p *ReferencePipeline[T]) Sequential() Stream[T] {
	p.parallel = false
	return p
}

func (p *ReferencePipeline[T]) getCtx() context.Context {
	return p.ctx
}

// initOp 初始化中间操作
func (p *ReferencePipeline[T]) initOp(opFlags int) (chan T, chan T) {
	if p.linkedOrConsumed {
		p.sendErr(errors.New("stream already operated upon or closed"))
	}
	in, ok := <-p.out
	if !ok {
		return p.closeChan(), p.closeChan()
	}
	p.sourceOrOpFlags = opFlags & OpMask
	p.combinedFlags = CombineOpFlags(opFlags, p.combinedFlags)
	p.depth++
	out := make(chan T, 1)
	p.out <- out
	return in, out
}

// initTerminalOp 初始化终端操作
func (p *ReferencePipeline[T]) initTerminalOp(terminalFlags int) chan T {
	if p.linkedOrConsumed {
		p.sendErr(errors.New("stream already operated upon or closed"))
	}
	in, ok := <-p.out
	if !ok {
		return p.closeChan()
	}
	p.combinedFlags = CombineOpFlags(terminalFlags, p.combinedFlags)
	return in
}

func (p *ReferencePipeline[T]) closeChan() chan T {
	ch := make(chan T)
	close(ch)
	return ch
}

func (p *ReferencePipeline[T]) close() {
	p.linkedOrConsumed = true
	p.cancel()  // 取消所有操作
	p.wg.Wait() // 等待所有协程结束

	// 安全关闭通道
	close(p.out)
	for ch := range p.out {
		close(ch)
	}
}

func (p *ReferencePipeline[T]) sendErr(err error) {
	p.errOnce.Do(func() {
		p.err = err
	})
	p.close()
}
