package utils

import (
	"crypto/md5"
)

const hextable = "0123456789abcdef"

func MD5V(str []byte) string {
	//sum := md5.Sum(str)
	//return hex.EncodeToString(sum[:])

	src := md5.Sum(str)
	var dst = make([]byte, 32)
	j := 0
	for _, v := range src {
		dst[j] = hextable[v>>4]
		dst[j+1] = hextable[v&0x0f]
		j += 2
	}
	return string(dst)
}
