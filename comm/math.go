package comm

import "golang.org/x/exp/constraints"

// Abs 获取一个数的绝对值
func Abs[T constraints.Integer](val T) T {
	if val < 0 {
		return -val
	}
	return val
}
