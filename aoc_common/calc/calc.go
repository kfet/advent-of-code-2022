package calc

func TaxiCab(x, y, x2, y2 int) int {
	return Abs(x-x2) + Abs(y-y2)
}

func Abs[T int | int64](a T) T {
	if a < 0 {
		return -a
	}
	return a
}

func Min(a, b int) int {
	if b < a {
		return b
	}
	return a
}

func Max(a, b int) int {
	if b > a {
		return b
	}
	return a
}
