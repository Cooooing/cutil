package stream

import (
	"context"
	"testing"
	"time"
)

func TestGroupBy(t *testing.T) {
	stream := Of[string](context.Background(), "a", "b", "cc")
	by, err := GroupBy[string, int](stream, func(i string) int {
		return len(i)
	})
	t.Logf("%v, %v", by, err)
}

func TestName(t *testing.T) {
	ch := make(chan int, 1024)
	go func() {
		for i := 1; ; i++ {
			time.Sleep(time.Second)
			for j := 0; j < 100; j++ {
				ch <- j * i
			}
		}
	}()
	err := OfChan[int](context.Background(), ch).
		Parallel(50).
		Map(func(i int) int {
			return i * i
		}).
		ForEach(func(i int) {
			time.Sleep(time.Millisecond)
			t.Logf("%v", i)
		})
	if err != nil {
		t.Errorf("%v", err)
	}
}
