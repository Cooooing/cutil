package bitmap

import (
	"strings"
	"sync"
)

type Enum[T ~uint64] struct {
	names map[T]string
	once  sync.Once
	all   []T
}

func NewEnum[T ~uint64](names map[T]string) *Enum[T] {
	return &Enum[T]{names: names}
}

func (e *Enum[T]) String(v T) string {
	if name, ok := e.names[v]; ok {
		return name
	}

	// 处理多选组合
	var parts []string
	for val, name := range e.names {
		if v&val != 0 {
			parts = append(parts, name)
		}
	}
	if len(parts) == 0 {
		return "None"
	}
	return strings.Join(parts, "|")
}

func (e *Enum[T]) Parse(name string) (T, bool) {
	name = strings.TrimSpace(strings.ToLower(name))
	for val, n := range e.names {
		if strings.ToLower(n) == name {
			return val, true
		}
	}
	var zero T
	return zero, false
}

func (e *Enum[T]) Values() []T {
	e.once.Do(func() {
		for v := range e.names {
			e.all = append(e.all, v)
		}
	})
	return e.all
}

// ---

type EnumSet[T ~uint64] struct {
	value T
	enum  *Enum[T]
}

func NewEnumSet[T ~uint64](e *Enum[T], initial ...T) *EnumSet[T] {
	var v T
	for _, i := range initial {
		v |= i
	}
	return &EnumSet[T]{value: v, enum: e}
}

func (s *EnumSet[T]) Add(flag T)         { s.value |= flag }
func (s *EnumSet[T]) Remove(flag T)      { s.value &^= flag }
func (s *EnumSet[T]) Has(flag T) bool    { return s.value&flag != 0 }
func (s *EnumSet[T]) Clear()             { var zero T; s.value = zero }
func (s *EnumSet[T]) Value() T           { return s.value }
func (s *EnumSet[T]) String() string     { return s.enum.String(s.value) }
func (s *EnumSet[T]) Clone() *EnumSet[T] { return &EnumSet[T]{value: s.value, enum: s.enum} }
