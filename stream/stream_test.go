package stream

import (
	"context"
	"testing"
)

func TestName(t *testing.T) {
	ints := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	ctx, cancel := context.WithCancelCause(context.Background())
	New[int](ctx).
		Of(ints...).
		Peek(func(i int) {
			t.Logf("peek: %d", i)
		}).
		Filter(func(i int) bool {
			return i%2 == 0
		}).
		Peek(func(i int) {
			t.Logf("peek after filter: %d", i)
		}).
		Skip(2).
		ForEach(func(i int) {
			t.Logf("forEach: %d", i)
		})
	cancel(nil)
}
