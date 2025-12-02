package base

import (
	"context"
)

func SetContextValue[K any, V any](ctx context.Context, key K, value V) context.Context {
	return context.WithValue(ctx, key, value)
}

func GetContextValue[K any, V any](ctx context.Context, key K) (V, bool) {
	value := ctx.Value(key)
	if IsNil(value) {
		var t V
		return t, false
	}
	return value.(V), true
}
