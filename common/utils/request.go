package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/youchuangcd/gopkg"
	"github.com/youchuangcd/gopkg/mylog"
	"io"
)

func GetParams(c *gin.Context) interface{} {
	var (
		param interface{}
		ok    bool
	)
	// 如果已经获取过了，就不再重新读取了
	if param, ok = c.Value(gopkg.ContextRequestParamKey).(interface{}); !ok {
		paramMap := make(map[string]interface{})
		if c.Request.Method == "GET" {
			for key, _ := range c.Request.URL.Query() {
				paramMap[key] = c.Query(key)
			}
		} else if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "DELETE" {
			contentType := c.ContentType()
			switch contentType {
			case "multipart/form-data":
				c.Request.ParseMultipartForm(200000)
				c.Request.ParseForm()
				for k, v := range c.Request.PostForm {
					paramMap[k] = v[len(v)-1]
				}
				param = paramMap
			case "application/x-www-form-urlencoded":
				err := c.Request.ParseForm()
				if err != nil {
					mylog.WithInfo(c, gopkg.LogCtx, map[string]interface{}{
						"body": c.Request.PostForm,
						"err":  err,
					}, "application/x-www-form-urlencoded读取body原始参数解析到map失败")
				}
				for k, v := range c.Request.PostForm {
					paramMap[k] = v[len(v)-1]
				}
				param = paramMap
			case "application/json":
				data, _ := c.GetRawData()
				c.Request.Body = io.NopCloser(bytes.NewBuffer(data)) // 读取完后再放回去
				if len(data) > 0 {
					var err error
					if data[0] == '[' { // 没有key的数组: [{"A1":"B", "C1":"D"},{"A2":"B", "C2":"D"}]
						err = json.Unmarshal(data, &param)
					} else {
						err = json.Unmarshal(data, &paramMap)
						param = paramMap
					}
					if err != nil {
						mylog.WithError(c, gopkg.LogCtx, map[string]interface{}{
							"body": string(data),
							"err":  err,
						}, "application/json读取body原始参数解析到map失败")
					}
				}
			default:
			}
		}
		// 给recover中间件使用
		c.Set(gopkg.ContextRequestParamKey, param)
	}

	return param
}

// GetUserIdByContext
// @Description: 从上下文中获取用户id
// @param ctx
// @return uint
// @return bool
func GetUserIdByContext(ctx context.Context) (uint, bool) {
	v, ok := ctx.Value(gopkg.ContextRequestUserIdKey).(uint)
	if !ok {
		return 0, false
	}
	return v, true
}

// GetSysUserIdByContext
// @Description: 从上下文中获取系统用户id
// @param ctx
// @return uint
// @return bool
func GetSysUserIdByContext(ctx context.Context) (uint, bool) {
	v, ok := ctx.Value(gopkg.ContextRequestSysUserIdKey).(uint)
	if !ok {
		return 0, false
	}
	return v, true
}

// GetTraceIdByContext
//
//	@Description: 从上下文中获取traceId
//	@param ctx
//	@return string
//	@return bool
func GetTraceIdByContext(ctx context.Context) (string, bool) {
	traceId, ok := ctx.Value(gopkg.RequestHeaderTraceIdKey).(string)
	return traceId, ok
}
