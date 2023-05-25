package gopkg

import (
	ut "github.com/go-playground/universal-translator"
)

var (
	// Validate 验证器 validator.Validate是线程安全的，其变量内会缓存已经验证过结构体的特征，因此用户用一个变量更有利于提高效率
	Trans ut.Translator
)
