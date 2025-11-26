package bitmap

import "strings"

// DynamicBitMap 非线程安全实现：基于 uint64 切片 words 存储位（每个 uint64 存 64 个 bit）
// 位编号从 0 开始（最低位为 0）
type DynamicBitMap struct {
	words []uint64 // 存储实际位的数组（每个元素 64 位）
	size  int      // 当前位图的有效长度（bit 数），bits >= 0
}

func New(n int) *DynamicBitMap {
	return NewDynamicBitMap(n)
}

// NewDynamicBitMap 根据初始位数 n 创建 DynamicBitMap（按 64 位对齐分配内存）。
// 如果 n <= 0，返回空的 DynamicBitMap（len=0，words=nil）。
func NewDynamicBitMap(n int) *DynamicBitMap {
	if n <= 0 {
		return &DynamicBitMap{}
	}
	w := (n + 63) >> 6 // 需要的 uint64 数量（向上取整）
	return &DynamicBitMap{
		words: make([]uint64, w),
		size:  n,
	}
}

// ensureWords 保证内部 words 切片至少有 w 个元素（可扩容，保留已有数据）。
func (b *DynamicBitMap) ensureWords(w int) {
	if len(b.words) < w {
		newWords := make([]uint64, w)
		copy(newWords, b.words)
		b.words = newWords
	}
}

// ensureIndex 确保位 index 可以被访问：当 index 超出当前 size 时自动扩容到 index+1。
func (b *DynamicBitMap) ensureIndex(i int) {
	if i < 0 || i < b.size {
		return
	}
	b.size = i + 1
	need := (b.size + 63) >> 6
	b.ensureWords(need)
}

// Set 将指定 index 的位设置为 1；如果 index 超出当前 size，会扩容到 index+1。
func (b *DynamicBitMap) Set(i int) {
	if i < 0 {
		return
	}
	b.ensureIndex(i)
	word, bit := i>>6, uint(i&63)
	b.words[word] |= 1 << bit
}

// Clear 将指定 index 位设置为 0；如果 index 超出当前 size，则什么也不做。
func (b *DynamicBitMap) Clear(i int) {
	if i < 0 || i >= b.size {
		return
	}
	word, bit := i>>6, uint(i&63)
	b.words[word] &^= 1 << bit // &^= 是 clear 指定位的惯用操作
}

// Flip 翻转指定 index 的位（0->1, 1->0）。若 index 超出则扩容。
func (b *DynamicBitMap) Flip(i int) {
	if i < 0 {
		return
	}
	b.ensureIndex(i)
	word, bit := i/64, i%64
	b.words[word] ^= 1 << bit
}

// Test 判断第 i 位是否为 1。
func (b *DynamicBitMap) Test(i int) bool {
	if i < 0 || i >= b.size {
		return false
	}
	return (b.words[i>>6] & (1 << (i & 63))) != 0
}

// SetRange 设置 [start, end] 范围内的所有位为 1
func (b *DynamicBitMap) SetRange(start, end int) {
	for i := start; i <= end; i++ {
		b.Set(i)
	}
}

// ClearRange 清空 [start, end] 范围内的所有位为 0
func (b *DynamicBitMap) ClearRange(start, end int) {
	for i := start; i <= end; i++ {
		b.Clear(i)
	}
}

func (b *DynamicBitMap) Union(other BitMap) BitMap {
	o, ok := other.(*DynamicBitMap)
	if !ok {
		return b.Clone()
	}
	maxSize := b.size
	if o.size > maxSize {
		maxSize = o.size
	}
	result := NewDynamicBitMap(maxSize)
	for i := 0; i < maxSize; i++ {
		if b.Test(i) || o.Test(i) {
			result.Set(i)
		}
	}
	return result
}

func (b *DynamicBitMap) Intersect(other BitMap) BitMap {
	o, ok := other.(*DynamicBitMap)
	if !ok {
		return NewDynamicBitMap(0)
	}
	minSize := b.size
	if o.size < minSize {
		minSize = o.size
	}
	result := NewDynamicBitMap(minSize)
	for i := 0; i < minSize; i++ {
		if b.Test(i) && o.Test(i) {
			result.Set(i)
		}
	}
	return result
}

func (b *DynamicBitMap) Difference(other BitMap) BitMap {
	o, ok := other.(*DynamicBitMap)
	if !ok {
		return b.Clone()
	}
	result := b.Clone().(*DynamicBitMap)
	for i := 0; i < result.size && i < o.size; i++ {
		if o.Test(i) {
			result.Clear(i)
		}
	}
	return result
}

func (b *DynamicBitMap) Clone() BitMap {
	newWords := make([]uint64, len(b.words))
	copy(newWords, b.words)
	return &DynamicBitMap{
		words: newWords,
		size:  b.size,
	}
}

func (b *DynamicBitMap) Equal(other BitMap) bool {
	o, ok := other.(*DynamicBitMap)
	if !ok {
		return false
	}
	if b.size != o.size {
		return false
	}
	for i := 0; i < len(b.words); i++ {
		if b.words[i] != o.words[i] {
			return false
		}
	}
	return true
}

// Len 返回位图当前长度（有效位数）。
func (b *DynamicBitMap) Len() int {
	return b.size
}

func (b *DynamicBitMap) ClearAll() {
	for i := range b.words {
		b.words[i] = 0
	}
}

func (b *DynamicBitMap) String() string {
	var sb strings.Builder
	for i := 0; i < b.size; i++ {
		if b.Test(i) {
			sb.WriteByte('1')
		} else {
			sb.WriteByte('0')
		}
	}
	return sb.String()
}

func (b *DynamicBitMap) Lock() {}

func (b *DynamicBitMap) Unlock() {}

func (b *DynamicBitMap) RLock() {}

func (b *DynamicBitMap) RUnlock() {}
