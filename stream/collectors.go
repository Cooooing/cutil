package stream

import (
	"cutil"
)

// Reduce 将流中的元素通过累积函数聚合为单一结果。
// identity: 归约的初始值，即使流为空也返回此值。
// accumulator: 累积函数，定义如何将流元素与中间结果合并。
// combiner: 合并函数，定义如何合并并行流中的子结果。
func Reduce[T any, R any](stream Stream[T], identity R, accumulator cutil.BiFunction[T, R], combiner cutil.BinaryOperator[R]) (R, error) {
	var err error
	result := identity
	if stream.IsParallel() {
		// Todo 并行流逻辑
	} else {
		// 顺序流逻辑
		err = stream.ForEach(func(item T) {
			result = accumulator(item, result)
		})
	}
	return result, err
}

func GroupBy[T any, K comparable](stream Stream[T], classifier cutil.Function[T, K]) (map[K][]T, error) {
	var err error
	result := make(map[K][]T)
	err = stream.ForEach(func(item T) {
		key := classifier(item)
		result[key] = append(result[key], item)
	})
	return result, err
}
