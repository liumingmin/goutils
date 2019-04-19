package utils

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func MinU(x, y uint32) uint32 {
	if x < y {
		return x
	}
	return y
}

func MaxU(x, y uint32) uint32 {
	if x > y {
		return x
	}
	return y
}

func Min64(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func Max64(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

func Abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
