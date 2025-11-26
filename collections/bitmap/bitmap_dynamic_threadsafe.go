package bitmap

import (
	"sync"
)

// ThreadSafeDynamicBitMap 是线程安全的位图实现。
type ThreadSafeDynamicBitMap struct {
	mu sync.RWMutex
	b  *DynamicBitMap
}

// NewThreadSafeDynamicBitMap 创建一个线程安全的动态位图实例。
func NewThreadSafeDynamicBitMap(n int) *ThreadSafeDynamicBitMap {
	return &ThreadSafeDynamicBitMap{b: NewDynamicBitMap(n)}
}

// Set 将第 i 位设置为 1。
func (s *ThreadSafeDynamicBitMap) Set(i int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.b.Set(i)
}

// Clear 清除第 i 位。
func (s *ThreadSafeDynamicBitMap) Clear(i int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.b.Clear(i)
}

// Flip 反转第 i 位。
func (s *ThreadSafeDynamicBitMap) Flip(i int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.b.Flip(i)
}

// Test 判断第 i 位是否为 1。
func (s *ThreadSafeDynamicBitMap) Test(i int) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.b.Test(i)
}

// SetRange 设置 [start, end] 区间的所有位。
func (s *ThreadSafeDynamicBitMap) SetRange(start, end int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.b.SetRange(start, end)
}

// ClearRange 清除 [start, end] 区间的所有位。
func (s *ThreadSafeDynamicBitMap) ClearRange(start, end int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.b.ClearRange(start, end)
}

// ====================== 集合操作 ======================

// Union 并集。
func (s *ThreadSafeDynamicBitMap) Union(other BitMap) BitMap {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.b = s.b.Union(other).(*DynamicBitMap)
	return s
}

// Intersect 交集。
func (s *ThreadSafeDynamicBitMap) Intersect(other BitMap) BitMap {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.b = s.b.Intersect(other).(*DynamicBitMap)
	return s
}

// Difference 差集。
func (s *ThreadSafeDynamicBitMap) Difference(other BitMap) BitMap {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.b = s.b.Difference(other).(*DynamicBitMap)
	return s
}

// Clone 返回副本。
func (s *ThreadSafeDynamicBitMap) Clone() BitMap {
	s.mu.RLock()
	defer s.mu.RUnlock()
	clone := &ThreadSafeDynamicBitMap{b: s.b.Clone().(*DynamicBitMap)}
	return clone
}

// Equal 判断是否相等。
func (s *ThreadSafeDynamicBitMap) Equal(other BitMap) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.b.Equal(other)
}

// Len 返回长度。
func (s *ThreadSafeDynamicBitMap) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.b.Len()
}

// ClearAll 清除所有位。
func (s *ThreadSafeDynamicBitMap) ClearAll() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.b.ClearAll()
}

// String 返回字符串形式。
func (s *ThreadSafeDynamicBitMap) String() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.b.String()
}

func (s *ThreadSafeDynamicBitMap) Lock()    { s.mu.Lock() }
func (s *ThreadSafeDynamicBitMap) Unlock()  { s.mu.Unlock() }
func (s *ThreadSafeDynamicBitMap) RLock()   { s.mu.RLock() }
func (s *ThreadSafeDynamicBitMap) RUnlock() { s.mu.RUnlock() }
