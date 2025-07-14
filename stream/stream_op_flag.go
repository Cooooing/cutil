package stream

// FlagType 标志类型
type FlagType int

const (
	FlagSpliterator        FlagType = iota // 该标志与spliterator特征相关联
	FlagStream                             // 该标志与流标志相关联
	FlagOp                                 // 该标志与中间操作标志相关联
	FlagTerminalOp                         // 该标志与终端操作标志相关联
	FlagUpstreamTerminalOp                 // 该标志与终端操作标志相关联，这些标志在最后一个有状态操作边界上向上游传播
)

// StreamOpFlag 流操作标志
type StreamOpFlag struct {
	BitPosition int
	Set         int
	Clear       int
	Preserve    int
	MaskTable   map[FlagType]int
}

// 标志位模式常数
const (
	setBits      = 0b01 // 设置标志
	clearBits    = 0b10 // 清除标志
	preserveBits = 0b11 // 保留标志
)

// 预定义标志
var (
	DISTINCT       = newStreamOpFlag(0, set(FlagSpliterator).set(FlagStream).setAndClear(FlagOp))
	SORTED         = newStreamOpFlag(1, set(FlagSpliterator).set(FlagStream).setAndClear(FlagOp))
	ORDERED        = newStreamOpFlag(2, set(FlagSpliterator).set(FlagStream).setAndClear(FlagOp).clear(FlagTerminalOp).clear(FlagUpstreamTerminalOp))
	SIZED          = newStreamOpFlag(3, set(FlagSpliterator).set(FlagStream).clear(FlagOp))
	SHORT_CIRCUIT  = newStreamOpFlag(12, set(FlagOp).set(FlagTerminalOp))
	SIZE_ADJUSTING = newStreamOpFlag(13, set(FlagOp))

	// 不同标志类型的掩码，用于过滤特定类型的标志
	SpliteratorCharacteristicsMask = createMask(FlagSpliterator)
	StreamMask                     = createMask(FlagStream)
	OpMask                         = createMask(FlagOp)
	TerminalOpMask                 = createMask(FlagTerminalOp)
	UpstreamTerminalOpMask         = createMask(FlagUpstreamTerminalOp)
	FlagMask                       = createFlagMask()
	FlagMaskIs                     = StreamMask
	FlagMaskNot                    = StreamMask << 1
	InitialOpsValue                = FlagMaskIs | FlagMaskNot

	// 特定标志的Bit值
	IsDistinct      = DISTINCT.Set
	NotDistinct     = DISTINCT.Clear
	IsSorted        = SORTED.Set
	NotSorted       = SORTED.Clear
	IsOrdered       = ORDERED.Set
	NotOrdered      = ORDERED.Clear
	IsSized         = SIZED.Set
	NotSized        = SIZED.Clear
	IsShortCircuit  = SHORT_CIRCUIT.Set
	IsSizeAdjusting = SIZE_ADJUSTING.Set
)

// maskBuilder 用于构建标志位掩码
type maskBuilder struct {
	m map[FlagType]int
}

func newMaskBuilder() *maskBuilder {
	return &maskBuilder{m: make(map[FlagType]int)}
}

func (b *maskBuilder) mask(t FlagType, i int) *maskBuilder {
	b.m[t] = i
	return b
}

func (b *maskBuilder) set(t FlagType) *maskBuilder {
	return b.mask(t, setBits)
}

func (b *maskBuilder) clear(t FlagType) *maskBuilder {
	return b.mask(t, clearBits)
}

func (b *maskBuilder) setAndClear(t FlagType) *maskBuilder {
	return b.mask(t, preserveBits)
}

func (b *maskBuilder) build() map[FlagType]int {
	for _, t := range []FlagType{FlagSpliterator, FlagStream, FlagOp, FlagTerminalOp, FlagUpstreamTerminalOp} {
		if _, exists := b.m[t]; !exists {
			b.m[t] = 0b00
		}
	}
	return b.m
}

func set(t FlagType) *maskBuilder {
	return newMaskBuilder().set(t)
}

func newStreamOpFlag(position int, builder *maskBuilder) StreamOpFlag {
	maskTable := builder.build()
	position *= 2 // 每个标志两个bit位
	return StreamOpFlag{
		BitPosition: position,
		Set:         setBits << position,
		Clear:       clearBits << position,
		Preserve:    preserveBits << position,
		MaskTable:   maskTable,
	}
}

// IsStreamFlag 检查该标志是否是基于流的标志
func (f StreamOpFlag) IsStreamFlag() bool {
	return f.MaskTable[FlagStream] > 0
}

// IsKnown 检查给定标志中是否设置了该标志
func (f StreamOpFlag) IsKnown(flags int) bool {
	return (flags & f.Preserve) == f.Set
}

// IsCleared 检查给定标志中的标志是否已清除
func (f StreamOpFlag) IsCleared(flags int) bool {
	return (flags & f.Preserve) == f.Clear
}

// IsPreserved 检查给定标志中是否保留了该标志
func (f StreamOpFlag) IsPreserved(flags int) bool {
	return (flags & f.Preserve) == f.Preserve
}

// CanSet 检查是否可以为给定类型设置标记
func (f StreamOpFlag) CanSet(t FlagType) bool {
	return (f.MaskTable[t] & setBits) > 0
}

func createMask(t FlagType) int {
	var mask int
	for _, flag := range []StreamOpFlag{DISTINCT, SORTED, ORDERED, SIZED, SHORT_CIRCUIT, SIZE_ADJUSTING} {
		mask |= flag.MaskTable[t] << flag.BitPosition
	}
	return mask
}

func createFlagMask() int {
	var mask int
	for _, flag := range []StreamOpFlag{DISTINCT, SORTED, ORDERED, SIZED, SHORT_CIRCUIT, SIZE_ADJUSTING} {
		mask |= flag.Preserve
	}
	return mask
}

func getMask(flags int) int {
	if flags == 0 {
		return FlagMask
	}
	return ^(flags | ((FlagMaskIs & flags) << 1) | ((FlagMaskNot & flags) >> 1))
}

// CombineOpFlags 将新流或操作标志与之前的组合标志组合在一起
func CombineOpFlags(newStreamOrOpFlags, prevCombOpFlags int) int {
	return (prevCombOpFlags & getMask(newStreamOrOpFlags)) | newStreamOrOpFlags
}

// ToStreamFlags 将组合的流和操作标志转换为流标志
func ToStreamFlags(combOpFlags int) int {
	return ((^combOpFlags) >> 1) & FlagMaskIs & combOpFlags
}

// ToCharacteristics 将流标志转换为FlagSplitterar特性，用于与底层数据源交互
func ToCharacteristics(streamFlags int) int {
	return streamFlags & SpliteratorCharacteristicsMask
}

// FromCharacteristics 将FlagSplitterar特性转换为流标志，特别处理 SORTED（非自然排序时不转换）。
func FromCharacteristics(characteristics int, hasCustomSort bool) int {
	if (characteristics&SpliteratorSorted) != 0 && hasCustomSort {
		// Do not propagate SORTED if it uses a custom comparator
		return characteristics & SpliteratorCharacteristicsMask & ^SpliteratorSorted
	}
	return characteristics & SpliteratorCharacteristicsMask
}

// FlagSpliterator 特征常数
const (
	SpliteratorDistinct = 0x00000001
	SpliteratorSorted   = 0x00000004
	SpliteratorOrdered  = 0x00000010
	SpliteratorSized    = 0x00000040
)
