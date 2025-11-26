package bitmap

import "github.com/Cooooing/cutil/collections"

// BitMap 定义一个位图应具备的基本功能
// 位图（BitMap）使用位(bit)来高效地表示布尔集合
type BitMap interface {

	//  基础操作

	// Set 将第 i 位设置为 1如果 i 超出当前容量，内部应自动扩容
	Set(i int)
	// Clear 将第 i 位清零（设置为 0）
	Clear(i int)
	// Flip 反转第 i 位（1→0，0→1）
	Flip(i int)
	// Test 判断第 i 位是否为 1
	Test(i int) bool
	// SetRange 设置 [start, end] 范围内的所有位为 1
	SetRange(start, end int)
	// ClearRange 清除 [start, end] 范围内的所有位
	ClearRange(start, end int)

	//  集合逻辑操作

	// Union 并集，执行按位或运算： b = b ∪ other
	Union(other BitMap) BitMap
	// Intersect 交集，执行按位与运算： b = b ∩ other
	Intersect(other BitMap) BitMap
	// Difference 差集，执行差集运算： b = b - other
	Difference(other BitMap) BitMap

	// Clone 返回当前 BitMap 的副本
	Clone() BitMap
	// Equal 判断两个 BitMap 是否完全相同
	Equal(other BitMap) bool
	// Len 返回当前位图的长度（能存储的最大 bit 数）
	Len() int

	//  工具方法

	// ClearAll 清除所有位
	ClearAll()
	// String 返回位图的字符串形式（如 "00101101"）
	String() string

	collections.Lockable
}
