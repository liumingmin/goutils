package utils

type T interface {
	comparable
}

func Min[T int | int32 | uint32 | int64 | uint64 | float32 | float64](x, y T) T {
	if x < y {
		return x
	}
	return y
}

func Max[T int | int32 | uint32 | int64 | uint64 | float32 | float64](x, y T) T {
	if x > y {
		return x
	}
	return y
}

func Abs[T int | int32 | int64 | float32 | float64](x T) T {
	if x < 0 {
		return -x
	}
	return x
}
