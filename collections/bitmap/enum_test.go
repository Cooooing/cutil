package bitmap

import (
	"fmt"
	"testing"
)

type Color uint64

const (
	Red Color = 1 << iota
	Green
	Blue
	Yellow
)

var Colors = NewEnum(map[Color]string{
	Red:    "Red",
	Green:  "Green",
	Blue:   "Blue",
	Yellow: "Yellow",
})

func TestEnum(t *testing.T) {
	t.Parallel()

	s := NewEnumSet(Colors, Red, Blue)

	fmt.Println("set:", s.String())             // Red|Blue
	fmt.Println("has blue:", s.Has(Blue))       // true
	fmt.Println("all values:", Colors.Values()) // [1 2 4 8]
	fmt.Println("value:", s.Value())
}
