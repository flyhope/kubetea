package comm

// Ptr 获取一个值的指针
func Ptr[T any](val T) *T {
	return &val
}
