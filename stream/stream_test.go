package stream

import (
	"context"
	"reflect"
	"testing"
)

func TestMap(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{name: "empty stream", input: []int{}, expected: []int{}},
		{name: "single element", input: []int{1}, expected: []int{1, 1}},
		{name: "multiple elements", input: []int{1, 2, 3}, expected: []int{1, 2, 3, 1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			stream1 := OfBlock(ctx, tt.input...)
			stream2 := OfNoBlock(ctx, tt.input...)

			stream3 := ConcatNoBlock(ctx, stream1, stream2)
			result, err := stream3.ToArray()
			if err != nil {
				t.Errorf("ToArray() error = %v", err)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("OfBlock() = %v, want %v", result, tt.expected)
			}
		})
	}
}
