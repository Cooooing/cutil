package stream

import (
	"context"
	"cutil"
	"reflect"
	"testing"
	"time"
)

// TestOf 测试 Of 方法，创建有限流
func TestOf(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{name: "empty stream", input: []int{}, expected: []int{}},
		{name: "single element", input: []int{1}, expected: []int{1}},
		{name: "multiple elements", input: []int{1, 2, 3}, expected: []int{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			stream := Of(ctx, tt.input...)
			result, err := stream.ToArray()
			if err != nil {
				t.Errorf("Of() error = %v, want nil", err)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Of() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestMap 测试 Map 方法
func TestMap(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    []int
		mapper   cutil.UnaryOperator[int]
		expected []int
	}{
		{name: "double elements", input: []int{1, 2, 3}, mapper: func(x int) int { return x * 2 }, expected: []int{2, 4, 6}},
		{name: "add one", input: []int{1, 2, 3}, mapper: func(x int) int { return x + 1 }, expected: []int{2, 3, 4}},
		{name: "empty stream", input: []int{}, mapper: func(x int) int { return x * 2 }, expected: []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			stream := Of(ctx, tt.input...).Map(tt.mapper)
			result, err := stream.ToArray()
			if err != nil {
				t.Errorf("Map() error = %v, want nil", err)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Map() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestFilter 测试 Filter 方法
func TestFilter(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		input     []int
		predicate cutil.Predicate[int]
		expected  []int
	}{
		{name: "even numbers", input: []int{1, 2, 3, 4}, predicate: func(x int) bool { return x%2 == 0 }, expected: []int{2, 4}},
		{name: "greater than 2", input: []int{1, 2, 3, 4}, predicate: func(x int) bool { return x > 2 }, expected: []int{3, 4}},
		{name: "empty stream", input: []int{}, predicate: func(x int) bool { return x%2 == 0 }, expected: []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			stream := Of(ctx, tt.input...).Filter(tt.predicate)
			result, err := stream.ToArray()
			if err != nil {
				t.Errorf("Filter() error = %v, want nil", err)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Filter() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestSkip 测试 Skip 方法
func TestSkip(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    []int
		n        int
		expected []int
	}{
		{name: "skip 2", input: []int{1, 2, 3, 4}, n: 2, expected: []int{3, 4}},
		{name: "skip all", input: []int{1, 2}, n: 3, expected: []int{}},
		{name: "skip zero", input: []int{1, 2, 3}, n: 0, expected: []int{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			stream := Of(ctx, tt.input...).Skip(tt.n)
			result, err := stream.ToArray()
			if err != nil {
				t.Errorf("Skip() error = %v, want nil", err)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Skip() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestLimit 测试 Limit 方法
func TestLimit(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    []int
		maxSize  int
		expected []int
	}{
		{name: "limit 2", input: []int{1, 2, 3, 4}, maxSize: 2, expected: []int{1, 2}},
		{name: "limit exceed", input: []int{1, 2}, maxSize: 3, expected: []int{1, 2}},
		{name: "limit zero", input: []int{1, 2, 3}, maxSize: 0, expected: []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			stream := Of(ctx, tt.input...).Limit(tt.maxSize)
			result, err := stream.ToArray()
			if err != nil {
				t.Errorf("Limit() error = %v, want nil", err)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Limit() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestDistinct 测试 Distinct 方法
func TestDistinct(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{name: "remove duplicates", input: []int{1, 2, 2, 3, 1}, expected: []int{1, 2, 3}},
		{name: "no duplicates", input: []int{1, 2, 3}, expected: []int{1, 2, 3}},
		{name: "empty stream", input: []int{}, expected: []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			stream := Of(ctx, tt.input...).Distinct()
			result, err := stream.ToArray()
			if err != nil {
				t.Errorf("Distinct() error = %v, want nil", err)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Distinct() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestSorted 测试 Sorted 方法
func TestSorted(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		input      []int
		comparator cutil.Comparator[int]
		expected   []int
	}{
		{name: "ascending order", input: []int{3, 1, 2}, comparator: func(a, b int) int { return a - b }, expected: []int{1, 2, 3}},
		{name: "descending order", input: []int{3, 1, 2}, comparator: func(a, b int) int { return b - a }, expected: []int{3, 2, 1}},
		{name: "empty stream", input: []int{}, comparator: func(a, b int) int { return a - b }, expected: []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			stream := Of(ctx, tt.input...).Sorted(tt.comparator)
			result, err := stream.ToArray()
			if err != nil {
				t.Errorf("Sorted() error = %v, want nil", err)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Sorted() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestForEach 测试 ForEach 方法
func TestForEach(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    []int
		consumer cutil.Consumer[int]
		expected []int
	}{
		{name: "collect elements", input: []int{1, 2, 3}, consumer: func(x int) {}, expected: []int{1, 2, 3}},
		{name: "empty stream", input: []int{}, consumer: func(x int) {}, expected: []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			collected := []int{}
			consumer := func(x int) { collected = append(collected, x) }
			stream := Of(ctx, tt.input...)
			err := stream.ForEach(consumer)
			if err != nil {
				t.Errorf("ForEach() error = %v, want nil", err)
			}
			if !reflect.DeepEqual(collected, tt.expected) {
				t.Errorf("ForEach() collected = %v, want %v", collected, tt.expected)
			}
		})
	}
}

// TestCount 测试 Count 方法
func TestCount(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    []int
		expected int
	}{
		{name: "three elements", input: []int{1, 2, 3}, expected: 3},
		{name: "empty stream", input: []int{}, expected: 0},
		{name: "single element", input: []int{1}, expected: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			stream := Of(ctx, tt.input...)
			count, err := stream.Count()
			if err != nil {
				t.Errorf("Count() error = %v, want nil", err)
			}
			if count != tt.expected {
				t.Errorf("Count() = %d, want %d", count, tt.expected)
			}
		})
	}
}

// TestParallel 测试并行流
func TestParallel(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	input := []int{1, 2, 3, 4, 5}
	stream := Of(ctx, input...).
		Parallel(2).
		Map(func(x int) int {
			time.Sleep(100 * time.Millisecond) // 模拟耗时操作
			return x * 2
		}).
		Sorted(func(a, b int) int {
			return a - b
		})
	result, err := stream.ToArray()
	if err != nil {
		t.Errorf("Parallel() error = %v, want nil", err)
	}
	expected := []int{2, 4, 6, 8, 10}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Parallel() = %v, want %v", result, expected)
	}
}

// TestContextCancel 测试上下文取消
func TestContextCancel(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream := Of(ctx, 1, 2, 3, 4, 5).Map(func(x int) int {
		time.Sleep(100 * time.Millisecond) // 模拟耗时操作
		return x
	})

	cancel() // 立即取消上下文
	_, err := stream.ToArray()
	if err == nil {
		t.Error("ToArray() expected context canceled error, got nil")
	}
	if err != context.Canceled {
		t.Errorf("ToArray() error = %v, want context.Canceled", err)
	}
}
