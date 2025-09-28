package _map

// Map 接口，定义字典应具备的基本操作
type Map[K any, V any] interface {

	// 基础操作

	Set(key K, value V)
	Get(key K) (V, bool)
	Remove(key K)
	RemoveAll(keys ...K)
	Pop(key K) (V, bool)
	Contains(key ...K) bool
	ContainsAll(keys ...K) bool
	ContainsAny(keys ...K) bool

	// 集合操作

	Keys() []K
	Values() []V
	Len() int
	IsEmpty() bool
	Merge(other Map[K, V])
	Equal(other Map[K, V]) bool
	Clone() Map[K, V]
	Clear()
	Reset()

	// 线程安全锁相关

	Lock()
	Unlock()
	RLock()
	RUnlock()
}
