package comm

import "math"

// Ptr 获取一个值的指针
func Ptr[T any](val T) *T {
	return &val
}

// SortMap 排序，数字越小越靠前，没有靠最后
type SortMap map[string]int

// SortVal 根据KEY获取排序值，不存在返回INT最大值
func (m SortMap) SortVal(key string) int {
	value, ok := m[key]
	if !ok {
		value = math.MaxInt
	}
	return value
}
