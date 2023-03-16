package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/youchuangcd/gopkg"
	"github.com/youchuangcd/gopkg/common/utils"
)

func TraceId(c *gin.Context) {
	traceId := c.Request.Header.Get(gopkg.RequestHeaderTraceIdKey)
	if len(traceId) == 0 || len(traceId) > 100 {
		traceId = utils.GenTraceId(c)
	}
	c.Set(gopkg.RequestHeaderTraceIdKey, traceId)
	// Response head
	c.Header(gopkg.RequestHeaderTraceIdKey, traceId)
	c.Next()
}
