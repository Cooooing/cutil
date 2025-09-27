package set

import (
	"fmt"
	"testing"
)

func TestSet(t *testing.T) {
	set := NewComparableSet[int](0)
	// 断言成具体类型
	cs, ok := set.(Set[int])
	if ok {
		fmt.Println("成功断言为 *ComparableSet[int]", cs.Len())
	} else {
		fmt.Println("断言失败")
	}
}
