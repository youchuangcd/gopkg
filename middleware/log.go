package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/youchuangcd/gopkg"
	"github.com/youchuangcd/gopkg/common/utils"
	"github.com/youchuangcd/gopkg/mylog"
	"time"
)

func LoggerHandle(logCategory string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 标记来源，业务逻辑用来做兼容处理
		rs := logCategory
		c.Set(gopkg.ContextRequestSourceKey, rs)
		// 提供给mylog.loggerEntry方法用
		c.Set(gopkg.ContextLogCategoryKey, logCategory)

		// 开始时间
		startTime := time.Now()
		c.Set(gopkg.ContextRequestStartTimeKey, startTime)

		// 请求内容
		//reqParam := utils.GetParams(c)
		reqParamByte, _ := json.Marshal(utils.GetParams(c))
		// 请求IP
		clientIP := c.ClientIP()
		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		reqUri := c.Request.Host + c.Request.RequestURI
		// 请求内容
		//reqHeader := c.Request.Header
		reqHeaderByte, _ := json.Marshal(c.Request.Header)

		// 日志json
		logContent := mylog.Fields{
			"client_ip":      clientIP,
			"request_method": reqMethod,
			"request_domain": c.Request.Host,
			"request_uri":    c.Request.URL.Path,
			"request_url":    reqUri,
			//"request_param":  reqParam,
			//"request_header": reqHeader,
			"request_param":  string(reqParamByte),
			"request_header": string(reqHeaderByte),
			"request_source": rs,
		}

		// 处理请求
		c.Next()

		// 执行时间
		latencyTime := float64(time.Since(startTime).Nanoseconds()) / 1e6

		// 状态码
		statusCode := c.Writer.Status()

		// 程序响应状态码
		responseData, _ := c.Value(gopkg.ContextResponseDataKey).(gin.H)
		responseDataByte, _ := json.Marshal(responseData["data"])

		// 日志json
		logContent["code"] = statusCode
		logContent["latency_time"] = latencyTime
		logContent["response_code"] = responseData["code"]
		logContent["response_msg"] = responseData["msg"]
		logContent["response_data"] = string(responseDataByte)
		logContent["log_type"] = gopkg.RequestLogTypeFlag
		mylog.WithInfo(c, logCategory, logContent, "request-end")
	}
}
