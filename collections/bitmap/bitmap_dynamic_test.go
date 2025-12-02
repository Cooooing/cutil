package bitmap

import (
	"testing"
)

func TestSetAndTest(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		indexes  []int
		expected []int
	}{
		{"empty bitmap", []int{}, []int{}},
		{"single set", []int{5}, []int{5}},
		{"multiple set", []int{0, 3, 63, 64, 127}, []int{0, 3, 63, 64, 127}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bm := NewDynamicBitMap(16)
			for _, i := range tt.indexes {
				bm.Set(i)
			}
			for _, i := range tt.expected {
				if !bm.Test(i) {
					t.Errorf("bit %d should be set but was not", i)
				}
			}
		})
	}
}

func TestClear(t *testing.T) {
	t.Parallel()
	bm := NewDynamicBitMap(16)
	bm.Set(1)
	bm.Set(5)
	bm.Set(63)
	bm.Clear(5)

	if bm.Test(5) {
		t.Errorf("bit 5 should be cleared")
	}
	if !bm.Test(1) || !bm.Test(63) {
		t.Errorf("other bits should remain set")
	}
}

func TestFlip(t *testing.T) {
	t.Parallel()
	bm := NewDynamicBitMap(16)
	bm.Set(10)
	bm.Flip(10) // 1 -> 0
	if bm.Test(10) {
		t.Errorf("bit 10 should be flipped to 0")
	}
	bm.Flip(10) // 0 -> 1
	if !bm.Test(10) {
		t.Errorf("bit 10 should be flipped to 1")
	}
}

func TestSetRangeAndClearRange(t *testing.T) {
	t.Parallel()
	bm := NewDynamicBitMap(16)
	bm.SetRange(2, 5)
	for i := 2; i <= 5; i++ {
		if !bm.Test(i) {
			t.Errorf("bit %d should be set", i)
		}
	}
	bm.ClearRange(3, 4)
	if bm.Test(3) || bm.Test(4) {
		t.Errorf("bits 3 and 4 should be cleared")
	}
}

func TestUnionAndIntersect(t *testing.T) {
	t.Parallel()
	a := NewDynamicBitMap(16)
	b := NewDynamicBitMap(16)
	a.Set(1)
	a.Set(3)
	b.Set(3)
	b.Set(4)

	union := a.Clone().Union(b)
	if !union.Test(1) || !union.Test(3) || !union.Test(4) {
		t.Errorf("Union failed: expected bits 1,3,4 set")
	}

	intersect := a.Clone().Intersect(b)
	if !intersect.Test(3) || intersect.Test(1) || intersect.Test(4) {
		t.Errorf("Intersect failed: expected only bit 3 set")
	}
}

func TestDifference(t *testing.T) {
	t.Parallel()
	a := NewDynamicBitMap(16)
	b := NewDynamicBitMap(16)
	a.SetRange(0, 5)
	b.Set(1)
	b.Set(2)
	diff := a.Clone().Difference(b)

	if diff.Test(1) || diff.Test(2) {
		t.Errorf("bits 1 and 2 should be cleared in Difference")
	}
	if !diff.Test(0) || !diff.Test(3) || !diff.Test(4) || !diff.Test(5) {
		t.Errorf("other bits should remain set in Difference")
	}
}

func TestCloneAndEqual(t *testing.T) {
	t.Parallel()
	a := NewDynamicBitMap(16)
	a.Set(10)
	b := a.Clone()

	if !a.Equal(b) {
		t.Errorf("Clone should produce an equal bitmap")
	}

	b.Set(20)
	if a.Equal(b) {
		t.Errorf("bitmaps should differ after modifying clone")
	}
}

func TestClearAll(t *testing.T) {
	t.Parallel()
	bm := NewDynamicBitMap(16)
	bm.SetRange(0, 10)
	bm.ClearAll()
	for i := 0; i < 10; i++ {
		if bm.Test(i) {
			t.Errorf("bit %d should be cleared", i)
		}
	}
}
