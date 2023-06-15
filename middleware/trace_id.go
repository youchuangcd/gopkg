package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/youchuangcd/gopkg"
	"github.com/youchuangcd/gopkg/common/utils"
)

func TraceId(c *gin.Context) {
	traceId := c.Request.Header.Get(gopkg.RequestHeaderTraceIdKey)
	parentSpanId := c.Request.Header.Get(gopkg.RequestHeaderSpanIdKey)
	if len(traceId) == 0 || len(traceId) > 100 {
		traceId = utils.GenTraceId(c)
	}
	c.Set(gopkg.RequestHeaderTraceIdKey, traceId)
	// 记录上游spanId
	c.Set(gopkg.RequestHeaderSpanIdKey, parentSpanId)
	if parentSpanId != "" {
		spanId := utils.GenTraceId(c)
		// 记录自己的spanId
		c.Set(gopkg.LogSpanIdKey, spanId)
		// 响应给客户端本应用的spanId
		c.Header(gopkg.RequestHeaderSpanIdKey, spanId)
	}
	// Response head
	c.Header(gopkg.RequestHeaderTraceIdKey, traceId)
	// 接收istio B3 header
	for _, key := range gopkg.RequestB3Headers {
		v := c.Request.Header.Get(key)
		if v != "" {
			c.Set(key, v)
		}
	}
	c.Next()
}
