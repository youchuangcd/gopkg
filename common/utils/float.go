package utils

const MIN = 0.000001

// MIN 为用户自定义的比较精度
func IsEqual(f1, f2 float64) bool {
	if f1 > f2 {
		return f1-f2 < MIN
	} else {
		return f2-f1 < MIN
	}
}
