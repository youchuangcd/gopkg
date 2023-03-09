package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
)

// AesEncryptCBC AES-CBC 加密
// key 必须是 16(AES-128)、24(AES-192) 或 32(AES-256) 字节的 AES 密钥；
// 初始化向量 iv 为随机的 16 位字符串 (必须是16位)，
// 解密需要用到这个相同的 iv，因此将它包含在密文的开头。
func AesEncryptCBC(plaintext string, key string) string {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("cbc decrypt err:", err)
		}
	}()

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return ""
	}

	blockSize := len(key)
	padding := blockSize - len(plaintext)%blockSize // 填充字节
	if padding == 0 {
		padding = blockSize
	}

	// 填充 padding 个 byte(padding) 到 plaintext
	plaintext += string(bytes.Repeat([]byte{byte(padding)}, padding))
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err = rand.Read(iv); err != nil { // 将同时写到 ciphertext 的开头
		return ""
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], []byte(plaintext))

	//return base64.StdEncoding.EncodeToString(ciphertext)
	return hex.EncodeToString(ciphertext)
}

// AesDecryptCBC AES-CBC 解密
func AesDecryptCBC(ciphertext string, key string) string {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("cbc decrypt err:", err)
		}
	}()

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return ""
	}

	//ciphercode, err := base64.StdEncoding.DecodeString(ciphertext)
	ciphercode, err := hex.DecodeString(ciphertext)
	if err != nil {
		return ""
	}

	iv := ciphercode[:aes.BlockSize]        // 密文的前 16 个字节为 iv
	ciphercode = ciphercode[aes.BlockSize:] // 正式密文

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphercode, ciphercode)

	plaintext := string(ciphercode) // ↓ 减去 padding
	return plaintext[:len(plaintext)-int(plaintext[len(plaintext)-1])]
}

func AesEncryptECB(data, key string) (r []byte) {
	src := []byte(data)
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return r
	}
	//判断加密快的大小
	blockSize := c.BlockSize()
	//填充
	src = PKCS5Padding(src, blockSize)
	//初始化加密数据接收切片
	encryptData := make([]byte, len(src))
	dst := encryptData

	for len(src) > 0 {
		c.Encrypt(dst, src[:blockSize])
		src = src[blockSize:]
		dst = dst[blockSize:]
	}
	return encryptData
}
func AesDecryptECB(data, key string) (r []byte) {
	src := []byte(data)
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return r
	}
	originData := make([]byte, len(src))
	dst := originData
	blockSize := c.BlockSize()

	for len(src) > 0 {
		c.Decrypt(dst, src[:blockSize])
		src = src[blockSize:]
		dst = dst[blockSize:]
	}
	originData = PKCS5UnPadding(originData)
	return originData
}

// AesEncryptECBHex
// @Description:
// @param data
// @param key
// @return s
func AesEncryptECBHex(data, key string) (s string) {
	return strings.ToUpper(hex.EncodeToString(AesEncryptECB(data, key)))
}

// AesDecryptECBHex
// @Description:
// @param data
// @param key
// @return s
func AesDecryptECBHex(data, key string) (s string) {
	d, _e := hex.DecodeString(data)
	if _e != nil {
		return ""
	}
	return string(AesDecryptECB(string(d), key))
}

// 明文补码算法
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// 明文减码算法
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	if length == 0 {
		return []byte{}
	}
	unpadding := int(origData[length-1])
	if length-unpadding < 0 {
		return []byte{}
	}
	return origData[:(length - unpadding)]
}
