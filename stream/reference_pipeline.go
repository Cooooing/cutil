package stream

import (
	"context"
	"cutil"
)

type ReferencePipeline[T any] struct {
	ctx context.Context
	out chan chan T
}

func newReferencePipeline[T any](ctx context.Context) *ReferencePipeline[T] {
	return &ReferencePipeline[T]{
		ctx: ctx,
		out: make(chan chan T, 1),
	}
}

func (p *ReferencePipeline[T]) Filter(predicate cutil.Predicate[T]) Stream[T] {
	in, ok := <-p.out
	if !ok {
		p.out <- p.closeChan()
		return p
	}
	out := make(chan T, 1)
	p.out <- out

	go func() {
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

func (p *ReferencePipeline[T]) Distinct() Stream[T] {
	// TODO implement me
	panic("implement me")
}

func (p *ReferencePipeline[T]) Sorted(comparator cutil.Comparator[T]) Stream[T] {
	// TODO implement me
	panic("implement me")
}

func (p *ReferencePipeline[T]) Peek(action cutil.Consumer[T]) Stream[T] {
	in, ok := <-p.out
	if !ok {
		p.out <- p.closeChan()
		return p
	}
	out := make(chan T, 1)
	p.out <- out

	go func() {
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

func (p *ReferencePipeline[T]) Limit(maxSize int64) Stream[T] {
	in, ok := <-p.out
	if !ok {
		p.out <- p.closeChan()
		return p
	}
	out := make(chan T, 1)
	p.out <- out

	go func() {
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

func (p *ReferencePipeline[T]) Skip(n int64) Stream[T] {
	in, ok := <-p.out
	if !ok {
		p.out <- p.closeChan()
		return p
	}
	out := make(chan T, 1)
	p.out <- out

	go func() {
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

func (p *ReferencePipeline[T]) ForEach(action cutil.Consumer[T]) {
	in, ok := <-p.out
	if !ok {
		p.out <- p.closeChan()
		return
	}
	for v := range in {
		action(v)
	}
}

func (p *ReferencePipeline[T]) ForEachOrdered(action cutil.Consumer[T]) {
	// TODO implement me
	panic("implement me")
}

func (p *ReferencePipeline[T]) ToArray() []T {
	in := <-p.out
	array := make([]T, 0)
	for v := range in {
		array = append(array, v)
	}
	return array
}

func (p *ReferencePipeline[T]) Count() int64 {
	in := <-p.out
	count := int64(0)
	for _ = range in {
		count++
	}
	return count
}

func (p *ReferencePipeline[T]) AnyMatch(predicate cutil.Predicate[T]) bool {
	// TODO implement me
	panic("implement me")
}

func (p *ReferencePipeline[T]) AllMatch(predicate cutil.Predicate[T]) bool {
	// TODO implement me
	panic("implement me")
}

func (p *ReferencePipeline[T]) NoneMatch(predicate cutil.Predicate[T]) bool {
	// TODO implement me
	panic("implement me")
}

func (p *ReferencePipeline[T]) closeChan() chan T {
	ch := make(chan T)
	close(ch)
	return ch
}

func (p *ReferencePipeline[T]) done() chan T {
	select {
	case ch, ok := <-p.out:
		if !ok {
			return p.closeChan()
		}
		return ch
	case <-p.ctx.Done():
		return p.closeChan()
	}
}

func (p *ReferencePipeline[T]) getCtx() context.Context {
	return p.ctx
}
