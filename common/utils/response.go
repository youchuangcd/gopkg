package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/youchuangcd/gopkg"
	"net/http"
	"reflect"
	"time"
)

// DefaultMapValue 声明一个空的data返回体
var defaultMapValue = make(map[string]interface{})

func ToJson(ctx *gin.Context, e *gopkg.Error, data interface{}, args ...interface{}) {
	if _, ok := data.(map[string]interface{}); !ok && IsNil(data) {
		data = defaultMapValue
	}
	responseData := gin.H{
		"code": e.GetCode(),
		"msg":  e.GetMsg(),
		"data": data,
	}
	statusCode := http.StatusOK
	if len(args) > 0 {
		statusCode = args[0].(int)
	}
	ctx.JSON(statusCode, responseData)
	// 走最后的中间件来发送结果，这样可以做到中间件改写响应结果
	ctx.Set(gopkg.ContextResponseDataKey, responseData)
	return
}

// MerchantResJson
//
//	@Description: 响应给商户的结构
//	@param ctx
//	@param e
//	@param data
//	@param args
func MerchantResJson(ctx *gin.Context, e *gopkg.Error, data interface{}, args ...interface{}) {
	if _, ok := data.(map[string]interface{}); !ok && IsNil(data) {
		data = defaultMapValue
	}
	responseData := gin.H{
		"code": e.GetCode(),
		"msg":  e.GetMsg(),
		"data": data,
	}
	// 响应给商户的结构
	code := 0
	msg := "SUCCESS"
	successBool := true
	traceId, _ := ctx.Value(gopkg.RequestHeaderTraceIdKey).(string)
	if e.GetCode() != gopkg.Success.GetCode() {
		code = 1
		successBool = false
		msg = e.GetMsg()
	}
	merchantResponseData := gin.H{
		"code":      code,
		"success":   successBool,
		"message":   msg,
		"serial_no": traceId,
		"data":      data,
		"timestamp": time.Now().UnixMilli(),
	}
	//merchantResponseData := gin.H{
	//	"nResult":    nResult,
	//	"vcResult":   vcResult,
	//	"vcSerialNo": traceId,
	//	"data":       data,
	//}
	statusCode := http.StatusOK
	if len(args) > 0 {
		statusCode = args[0].(int)
	}
	ctx.JSON(statusCode, merchantResponseData)
	// 走最后的中间件来发送结果，这样可以做到中间件改写响应结果
	ctx.Set(gopkg.ContextResponseDataKey, responseData)
	return
}

func ToRaw(ctx *gin.Context, data interface{}) {
	ctx.String(http.StatusOK, "%v", data)
	// 走最后的中间件来发送结果，这样可以做到中间件改写响应结果
	e := gopkg.Success
	responseData := gin.H{
		"code": e.GetCode(),
		"msg":  e.GetMsg(),
		"data": data,
	}
	ctx.Set(gopkg.ContextResponseDataKey, responseData)
	return
}

// IsNil
// @Description: 判断interface的值是否为nil, 只判断指针、切片、map、chan、func
// @param i
// @return bool
func IsNil(i interface{}) bool {
	ret := i == nil
	if !ret {
		vi := reflect.ValueOf(i)
		kind := vi.Kind()
		if kind == reflect.Slice ||
			kind == reflect.Map ||
			kind == reflect.Chan ||
			kind == reflect.Interface ||
			kind == reflect.Func ||
			kind == reflect.Ptr {
			return vi.IsNil()
		}
	}
	return ret
}

func RetJson(ctx *gin.Context, e *gopkg.Error, data interface{}, args ...interface{}) {
	responseData := gin.H{
		"code":         e.GetCode(),
		"success":      true,
		"message":      e.GetMsg(),
		"service_time": time.Now().Unix(),
		"serial_no":    "",
		"data":         data,
	}
	statusCode := http.StatusOK
	if len(args) > 0 {
		statusCode = args[0].(int)
	}
	ctx.JSON(statusCode, responseData)
	// 走最后的中间件来发送结果，这样可以做到中间件改写响应结果
	ctx.Set(gopkg.ContextResponseDataKey, responseData)

	return
}
