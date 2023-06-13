package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5V(str []byte) string {
	//h := md5.New()
	//h.Write(str)
	//return hex.EncodeToString(h.Sum(b))
	sum := md5.Sum(str)
	return hex.EncodeToString(sum[:])
}
