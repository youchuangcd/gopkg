package utils

import (
	"fmt"
	"strconv"
)

// Decimal
// @Description: 浮点数保留几位小数
// @param v
// @param args
// @return float64
func Decimal(v float64, args ...int) float64 {
	// 默认保留3位，第二位参数指定保留位数
	format := "%.3f"
	if len(args) == 1 {
		format = fmt.Sprintf("%%.%df", args[0])
	}
	v, _ = strconv.ParseFloat(fmt.Sprintf(format, v), 64)
	return v
}

const (
	MININT64 = -922337203685477580
	MAXINT64 = 9223372036854775807
)

func Max(nums ...int64) int64 {
	var maxNum int64 = MININT64
	for _, num := range nums {
		if num > maxNum {
			maxNum = num
		}
	}
	return maxNum
}

func Min(nums ...int64) int64 {
	var minNum int64 = MAXINT64
	for _, num := range nums {
		if num < minNum {
			minNum = num
		}
	}
	return minNum
}

func Sum(nums ...int64) int64 {
	var sumNum int64 = 0
	for _, num := range nums {
		sumNum += num
	}
	return sumNum
}
