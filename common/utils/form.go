package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/youchuangcd/gopkg"
	"github.com/youchuangcd/gopkg/mylog"
	"strings"
)

// BindValid
// @Description: 绑定参数并翻译验证错误
// @param c
// @param form
// @return error
func BindValid(c *gin.Context, form interface{}) *gopkg.Error {
	var err error
	if c.Request.Method == "GET" {
		err = c.ShouldBindQuery(form)
	} else {
		err = c.ShouldBind(form)
	}
	if err != nil {
		logContent := map[string]interface{}{
			"form":         form,
			"err":          err,
			"content-type": c.ContentType(),
		}
		invalidParamErr := gopkg.InvalidParam
		invalidParamErr.Set("绑定参数校验失败")
		// 获取validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			if gopkg.Trans != nil {
				errMap := removeTopStruct(errs.Translate(gopkg.Trans))
				logContent["err"] = errMap
				var errBuff strings.Builder
				for _, v := range errMap {
					errBuff.WriteString(v)
					errBuff.WriteString(",")
				}
				invalidParamErr.Set(strings.TrimRight(errBuff.String(), ","))
			}
		}
		mylog.WithWarn(c, gopkg.LogCtx, logContent, "绑定参数校验失败")
		return invalidParamErr
	}

	return gopkg.Success
}

func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}
