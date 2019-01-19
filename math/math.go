package math

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
