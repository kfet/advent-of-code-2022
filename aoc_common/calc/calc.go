package calc

func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func LCD(a, b int) int {
	return a * b / GCD(a, b)
}

func IsPresent[K comparable, V any](m map[K]V, key K) bool {
	if _, ok := m[key]; ok {
		return true
	}
	return false
}

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
