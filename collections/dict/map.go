package dict

import (
	"context"

	"github.com/Cooooing/cutil/base"
	"github.com/Cooooing/cutil/stream"
)

// Map 接口，定义字典应具备的基本操作
type Map[K any, V any] interface {

	// 基础操作

	Set(key K, value V)
	Get(key K) (V, bool)
	Remove(key K)
	RemoveAll(keys ...K)
	Pop(key K) (V, bool)
	Contains(key K) bool
	ContainsAll(keys ...K) bool
	ContainsAny(keys ...K) bool

	// 集合操作

	Keys() []K
	Values() []V
	Entries() []*Entry[K, V]
	Len() int
	IsEmpty() bool
	Merge(other Map[K, V])
	Equal(other Map[K, V]) bool
	EqualFunc(other Map[K, V], fn base.Equator[V]) bool
	Clone() Map[K, V]
	Clear()
	Reset()
	Foreach(action base.Predicate[*Entry[K, V]])
	Stream(ctx context.Context) stream.Stream[*Entry[K, V]]

	// 线程安全锁相关

	Lock()
	Unlock()
	RLock()
	RUnlock()
}

type Entry[T any, R any] struct {
	Key   T
	Value R
}
