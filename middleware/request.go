package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/youchuangcd/gopkg"
	"github.com/youchuangcd/gopkg/common/utils"
)

// 中间件(验证参数签名)
func CheckSign(c *gin.Context) {
	traceId := c.Request.Header.Get(gopkg.RequestHeaderTraceIdKey)
	if len(traceId) == 0 || len(traceId) > 100 {
		traceId = utils.GenTraceId(c)
	}
	c.Set(gopkg.RequestHeaderTraceIdKey, traceId)
	// Response head
	c.Header(gopkg.RequestHeaderTraceIdKey, traceId)

	//请求参数
	// 验证通过，会继续访问下一个中间件

	c.Next()
}
