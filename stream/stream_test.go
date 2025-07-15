package stream

import (
	"context"
	"testing"
	"time"
)

func TestOf(t *testing.T) {
	t.Run("empty stream", func(t *testing.T) {
		ctx := context.Background()
		s := Of[int](ctx)
		count, err := s.Count()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if count != 0 {
			t.Errorf("expected count 0, got %d", count)
		}
	})

	t.Run("stream with values", func(t *testing.T) {
		ctx := context.Background()
		s := Of(ctx, 1, 2, 3)
		arr, err := s.ToArray()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(arr) != 3 {
			t.Errorf("expected 3 elements, got %d", len(arr))
		}
	})
}

func TestOfChan(t *testing.T) {
	t.Run("channel stream", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		ch := make(chan int)
		go func() {
			defer close(ch)
			ch <- 1
			ch <- 2
			ch <- 3
		}()

		s := OfChan(ctx, ch)
		count, err := s.Count()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if count != 3 {
			t.Errorf("expected count 3, got %d", count)
		}
	})
}

func TestMap(t *testing.T) {
	t.Run("basic mapping", func(t *testing.T) {
		ctx := context.Background()
		s := Of(ctx, 1, 2, 3)
		mapped := Map(s, func(x int) int { return x * 2 })
		arr, err := mapped.ToArray()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := []int{2, 4, 6}
		if len(arr) != len(expected) {
			t.Fatalf("expected %v, got %v", expected, arr)
		}
		for i := range expected {
			if arr[i] != expected[i] {
				t.Errorf("at index %d: expected %d, got %d", i, expected[i], arr[i])
			}
		}
	})
}

func TestFilter(t *testing.T) {
	t.Run("basic filtering", func(t *testing.T) {
		ctx := context.Background()
		s := Of(ctx, 1, 2, 3, 4, 5)
		filtered := s.Filter(func(x int) bool { return x%2 == 0 })
		arr, err := filtered.ToArray()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := []int{2, 4}
		if len(arr) != len(expected) {
			t.Fatalf("expected %v, got %v", expected, arr)
		}
		for i := range expected {
			if arr[i] != expected[i] {
				t.Errorf("at index %d: expected %d, got %d", i, expected[i], arr[i])
			}
		}
	})
}

func TestSkipAndLimit(t *testing.T) {
	t.Run("skip and limit", func(t *testing.T) {
		ctx := context.Background()
		s := Of(ctx, 1, 2, 3, 4, 5)
		limited := s.Skip(1).Limit(2)
		arr, err := limited.ToArray()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := []int{2, 3}
		if len(arr) != len(expected) {
			t.Fatalf("expected %v, got %v", expected, arr)
		}
		for i := range expected {
			if arr[i] != expected[i] {
				t.Errorf("at index %d: expected %d, got %d", i, expected[i], arr[i])
			}
		}
	})
}

func TestForEach(t *testing.T) {
	t.Run("basic forEach", func(t *testing.T) {
		ctx := context.Background()
		s := Of(ctx, 1, 2, 3)
		var sum int
		err := s.ForEach(func(x int) { sum += x })
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sum != 6 {
			t.Errorf("expected sum 6, got %d", sum)
		}
	})
}

func TestForEachOrdered(t *testing.T) {
	t.Run("ordered iteration", func(t *testing.T) {
		ctx := context.Background()
		s := Of(ctx, 3, 1, 2)
		var result []int
		err := s.ForEachOrdered(func(a, b int) int { return a - b }, func(x int) {
			result = append(result, x)
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := []int{1, 2, 3}
		for i := range expected {
			if result[i] != expected[i] {
				t.Errorf("at index %d: expected %d, got %d", i, expected[i], result[i])
			}
		}
	})
}

func TestMatchOperations(t *testing.T) {
	t.Run("anyMatch", func(t *testing.T) {
		ctx := context.Background()
		s := Of(ctx, 1, 2, 3)
		matched, err := s.AnyMatch(func(x int) bool { return x > 2 })
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !matched {
			t.Error("expected anyMatch to find element > 2")
		}
	})

	t.Run("allMatch", func(t *testing.T) {
		ctx := context.Background()
		s := Of(ctx, 1, 2, 3)
		matched, err := s.AllMatch(func(x int) bool { return x > 0 })
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !matched {
			t.Error("expected allMatch to find all elements > 0")
		}
	})

	t.Run("noneMatch", func(t *testing.T) {
		ctx := context.Background()
		s := Of(ctx, 1, 2, 3)
		matched, err := s.NoneMatch(func(x int) bool { return x > 3 })
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !matched {
			t.Error("expected noneMatch to find no elements > 3")
		}
	})
}

func TestContextCancellation(t *testing.T) {
	t.Run("cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // immediately cancel

		s := Of(ctx, 1, 2, 3)
		err := s.ForEach(func(x int) {
			t.Errorf("unexpected element: %d", x)
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestFlatMap(t *testing.T) {
	t.Run("basic flatMap", func(t *testing.T) {
		ctx := context.Background()
		s := Of(ctx, []int{1, 2}, []int{3, 4})
		flattened := FlatMap(s, func(x []int) Stream[int] {
			return Of(ctx, x...)
		})
		arr, err := flattened.ToArray()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := []int{1, 2, 3, 4}
		if len(arr) != len(expected) {
			t.Fatalf("expected %v, got %v", expected, arr)
		}
		for i := range expected {
			if arr[i] != expected[i] {
				t.Errorf("at index %d: expected %d, got %d", i, expected[i], arr[i])
			}
		}
	})
}

func TestConcat(t *testing.T) {
	t.Run("concatenate streams", func(t *testing.T) {
		ctx := context.Background()
		s1 := Of(ctx, 1, 2)
		s2 := Of(ctx, 3, 4)
		concatenated := Concat(ctx, s1, s2)
		arr, err := concatenated.ToArray()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := []int{1, 2, 3, 4}
		if len(arr) != len(expected) {
			t.Fatalf("expected %v, got %v", expected, arr)
		}
		for i := range expected {
			if arr[i] != expected[i] {
				t.Errorf("at index %d: expected %d, got %d", i, expected[i], arr[i])
			}
		}
	})
}

func TestGenerate(t *testing.T) {
	t.Run("infinite stream with limit", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		counter := 0
		s := Generate(ctx, func() int {
			counter++
			return counter
		})
		limited := s.Limit(3)
		arr, err := limited.ToArray()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := []int{1, 2, 3}
		if len(arr) != len(expected) {
			t.Fatalf("expected %v, got %v", expected, arr)
		}
		for i := range expected {
			if arr[i] != expected[i] {
				t.Errorf("at index %d: expected %d, got %d", i, expected[i], arr[i])
			}
		}
	})
}
