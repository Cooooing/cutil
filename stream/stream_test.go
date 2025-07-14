package stream

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	ints := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	strs := []string{"Go", "Java", "Python", "C", "C++", "Rust", "C#", "PHP", "JavaScript", "TypeScript"}

	ctx, cancel := context.WithCancelCause(context.Background())
	Of[int](ctx, ints...).
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
	fmt.Println()
	Map[int, int](Of[int](ctx, ints...), func(i int) int {
		return i * 2
	}).ForEach(func(i int) {
		t.Logf("forEach: %d", i)
	})
	fmt.Println()
	FlatMap[string, string](Of[string](ctx, strs...), func(s string) Stream[string] {
		var t []string
		for i := range []byte(s) {
			t = append(t, string(s[i]))
		}
		return Of[string](ctx, t...)
	}).ForEach(func(s string) {
		t.Logf("forEach: %s", s)
	})

	fmt.Println()
	ch := make(chan int)
	go func() {
		defer close(ch)
		timer := time.NewTicker(time.Second)
		defer timer.Stop()
		for i := 0; i < 10; i++ {
			select {
			case <-timer.C:
				ch <- 1
			case <-ctx.Done():
				return
			}
		}
	}()
	OfChan[int](ctx, ch).ForEach(func(i int) {
		t.Logf("forEach: %d", i)
	})

	cancel(nil)
}
