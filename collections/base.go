package collections

// Keyer 接口，用于自定义集合元素的键值
type Keyer interface {
	// Key 获取键值
	Key() string
}
