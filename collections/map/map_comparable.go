package _map

import (
	"context"
	"reflect"

	"github.com/Cooooing/cutil/common"
	"github.com/Cooooing/cutil/stream"
)

type MapComparable[K comparable, V any] map[K]V

func NewMapComparable[K comparable, V any](size int) Map[K, V] {
	m := make(MapComparable[K, V], size)
	return &m
}

func (m *MapComparable[K, V]) Set(key K, value V) {
	(*m)[key] = value
}

func (m *MapComparable[K, V]) Get(key K) (V, bool) {
	v, ok := (*m)[key]
	return v, ok
}

func (m *MapComparable[K, V]) Remove(key K) {
	delete(*m, key)
}

func (m *MapComparable[K, V]) RemoveAll(keys ...K) {
	for _, key := range keys {
		delete(*m, key)
	}
}

func (m *MapComparable[K, V]) Pop(key K) (V, bool) {
	v, ok := (*m)[key]
	if ok {
		delete(*m, key)
	}
	return v, ok
}

func (m *MapComparable[K, V]) Contains(key K) bool {
	_, ok := (*m)[key]
	return ok
}

func (m *MapComparable[K, V]) ContainsAll(keys ...K) bool {
	for _, key := range keys {
		if _, ok := (*m)[key]; !ok {
			return false
		}
	}
	return true
}

func (m *MapComparable[K, V]) ContainsAny(keys ...K) bool {
	for _, key := range keys {
		if _, ok := (*m)[key]; ok {
			return true
		}
	}
	return false
}

func (m *MapComparable[K, V]) Keys() []K {
	keys := make([]K, 0, m.Len())
	for k, _ := range *m {
		keys = append(keys, k)
	}
	return keys
}

func (m *MapComparable[K, V]) Values() []V {
	values := make([]V, 0, m.Len())
	for _, v := range *m {
		values = append(values, v)
	}
	return values
}

func (m *MapComparable[K, V]) Entries() []*Entry[K, V] {
	entries := make([]*Entry[K, V], 0, m.Len())
	for k, v := range *m {
		entries = append(entries, &Entry[K, V]{
			Key:   k,
			Value: v,
		})
	}
	return entries
}

func (m *MapComparable[K, V]) Len() int {
	return len(*m)
}

func (m *MapComparable[K, V]) IsEmpty() bool {
	return m.Len() == 0
}

func (m *MapComparable[K, V]) Merge(other Map[K, V]) {
	other.Foreach(func(e *Entry[K, V]) bool {
		m.Set(e.Key, e.Value)
		return true
	})
}

func (m *MapComparable[K, V]) Equal(other Map[K, V]) bool {
	if m.Len() != other.Len() {
		return false
	}
	for k, v := range *m {
		if ov, ok := other.Get(k); !ok || !reflect.DeepEqual(ov, v) {
			return false
		}
	}
	return true
}

func (m *MapComparable[K, V]) EqualFunc(other Map[K, V], fn common.Equator[V]) bool {
	if m.Len() != other.Len() {
		return false
	}
	for k, v := range *m {
		if ov, ok := other.Get(k); !ok || !fn(ov, v) {
			return false
		}
	}
	return true
}

func (m *MapComparable[K, V]) Clone() Map[K, V] {
	clone := NewMapComparable[K, V](m.Len())
	for k, v := range *m {
		clone.Set(k, v)
	}
	return clone
}

func (m *MapComparable[K, V]) Clear() {
	for k, _ := range *m {
		delete(*m, k)
	}
}

func (m *MapComparable[K, V]) Reset() {
	*m = make(MapComparable[K, V])
}

func (m *MapComparable[K, V]) Foreach(action common.Predicate[*Entry[K, V]]) {
	for item := range *m {
		action(&Entry[K, V]{
			Key:   item,
			Value: (*m)[item],
		})
	}
}

func (m *MapComparable[K, V]) Stream(ctx context.Context) stream.Stream[*Entry[K, V]] {
	return stream.OfBlock[*Entry[K, V]](ctx, m.Entries()...)
}

func (m *MapComparable[K, V]) Lock() {
	return
}

func (m *MapComparable[K, V]) Unlock() {
	return
}

func (m *MapComparable[K, V]) RLock() {
	return
}

func (m *MapComparable[K, V]) RUnlock() {
	return
}
