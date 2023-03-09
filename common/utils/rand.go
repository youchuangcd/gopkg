package utils

import (
	"math/rand"
	"time"
)

// 获取随机字符串
func RandSeq(n int) string {
	var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// 获取随机数字
func RandNumber(n int) string {
	var letters = []rune("0123456789")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// RandInt
// @Description: 获取指定范围内的随机数
// @param min 可为负数
// @param max 不可为负数; 包含max
// @return r int
func RandInt(min, max int) int {
	if min >= max || max <= 0 {
		return max
	}
	rand.Seed(time.Now().UnixNano())
	// 因为随机数是左开右闭，所以要包含max需要+1
	return rand.Intn(max-min+1) + min
}
