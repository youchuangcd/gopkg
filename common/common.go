package common

import (
	"encoding/base64"
	"fmt"
	"github.com/youchuangcd/gopkg"
	"math"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

func EnvLocal() bool {
	return gopkg.Env == gopkg.EnvLocal
}

func EnvDev() bool {
	return gopkg.Env == gopkg.EnvDev
}

func EnvTest() bool {
	return gopkg.Env == gopkg.EnvTest
}

func EnvGray() bool {
	return gopkg.Env == gopkg.EnvGray
}

func EnvProd() bool {
	return gopkg.Env == gopkg.EnvProd
}

// GetResponseError
// @Description: 获取响应的错误码
// @param err
// @return _err
func GetResponseError(err error, args ...interface{}) (_err *gopkg.Error) {
	_err = gopkg.Success
	if err != nil {
		_err = gopkg.Failure
		if len(args) > 0 {
			if e, ok := args[0].(*gopkg.Error); ok {
				_err = e
			}
		}

		if e, ok := err.(*gopkg.Error); ok {
			_err = e
		}
	}
	return _err
}

// GetCurrentFuncName
// @Description: 获取当前方法名称
// @return string
func GetCurrentFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}

// GetCallerFuncName
// @Description: 获取当前方法的调用者的方法名称
// @return string
func GetCallerFuncName() string {
	pc, _, _, _ := runtime.Caller(2)
	return runtime.FuncForPC(pc).Name()
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

// Clone
// @Description: 只支持导出的字段，不能设置未导出的字段
// @param oldObj
// @return interface{}
func Clone(oldObj interface{}) interface{} {
	newObj := reflect.New(reflect.TypeOf(oldObj).Elem())
	oldVal := reflect.ValueOf(oldObj).Elem()
	newVal := newObj.Elem()
	for i := 0; i < oldVal.NumField(); i++ {
		newValField := newVal.Field(i)
		if newValField.CanSet() {
			newValField.Set(oldVal.Field(i))
		}
	}

	return newObj.Interface()
}

func CreateFolder(s string) error {
	err := os.MkdirAll(s, 0766)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

// SetErrorMsg
// @Description: 设置自定义错误信息
// @param err
// @param msg
// @param args true=追加
// @return gopkg.ErrorInterface
func SetErrorMsg(err *gopkg.Error, msg string, args ...bool) *gopkg.Error {
	if len(args) == 1 && args[0] {
		err.Append(msg)
		//err.ErrMsg += msg
	} else {
		//err.ErrMsg = msg
		err.Set(msg)
	}
	return err
}

// Int64AmountToFloat64
// @Description: 整数型金额转换为float64
// @param price
// @return float64
func Int64AmountToFloat64(amount int64) float64 {
	return float64(amount) / 100
}

// Float64AmountToInt64
// @Description: 浮点金额转为整数型金额
// @param amount
// @return int64
func Float64AmountToInt64(amount float64) int64 {
	return int64(amount * 100)
}

// UnitsPerIntRateMulInt64Amount
// @Description: 万分单位整数比例 乘以 整数金额
// @param rate
// @param amount
// @return int64
func UnitsPerIntRateMulInt64Amount(rate int64, amount int64) int64 {
	return (rate * amount) / 10000
}

// float64向下取整
func FloatFloorToInt64(x float64) int64 {
	return int64(math.Floor(x))
}

// 将float64转成精确的int64
func Wrap(num float64, retain int) int64 {
	return int64(num * math.Pow10(retain))
}

// 将int64恢复成正常的float64
func Unwrap(num int64, retain int) float64 {
	return float64(num) / math.Pow10(retain)
}

// 精准float64
func WrapToFloat64(num float64, retain int) float64 {
	return num * math.Pow10(retain)
}

// 精准int64
func UnwrapToInt64(num int64, retain int) int64 {
	return int64(Unwrap(num, retain))
}

// 价格字符串转int
func StringPriceToInt(price string) int {
	float64Price, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return 0
	}
	intPrice := float64Price * 100
	return int(intPrice)
}

// GetMobileJoinAreaCode
// @Description: 获取拼接了手机区号的完整手机号，去除了+号
// @param m
// @param args
// @return string
func GetMobileJoinAreaCode(m string, args ...string) string {
	areaCode := "86"
	if len(args) == 1 && args[0] != "" {
		areaCode = strings.TrimLeft(args[0], "+")
	}
	return areaCode + m
}

// ValidatePlatform
//
//	@Description: 验证平台类型是否有效
//	@param platform
//	@return valid
func ValidatePlatform(platform uint16, platformSlice []uint16) (valid bool) {
	for _, v := range platformSlice {
		if platform == v {
			valid = true
			break
		}
	}
	return
}

// IsMobile
//
//	@Description: 检测是否是手机号
//	@param mobile
//	@return bool
func IsMobile(mobile string) bool {
	res, _ := regexp.MatchString(`^(?:\+?86)?1\d{10}$`, mobile)
	return res
}

// GetKsAvatarTraitCode
//
//	@Description: 获取快手头像特征码
//	@param avatar
//	@return traitCode
//	@return err
func GetKsAvatarTraitCode(avatar string) (traitCode string, err error) {
	if avatar == "" {
		return
	}
	anchorArr := strings.Split(avatar, "/")
	if len(anchorArr) > 0 {
		tmp := anchorArr[len(anchorArr)-1]
		if strings.Contains(tmp, ":") {
			tmp = strings.Split(tmp, ":")[0]
		} else if strings.Contains(tmp, ".") {
			tmp = strings.Split(tmp, ".")[0]
		}
		tmp = strings.TrimSuffix(tmp, "1y8")
		if len(tmp) > 20 {
			var anchorImgByte []byte
			anchorImgByte, err = base64.RawURLEncoding.DecodeString(tmp)
			if err == nil {
				anchorImgArr := strings.Split(string(anchorImgByte), ".")
				if len(anchorImgArr) > 0 {
					traitCode = strings.TrimLeft(strings.TrimSpace(anchorImgArr[0]), ";-")
				}
			}
		} else {
			traitCode = tmp
		}

	}
	return
}

// GetDyAvatarTraitCode
//
//	@Description: 抖音头像获取特征码
//	@param avatar
//	@return traitCode
//	@return err
func GetDyAvatarTraitCode(avatar string) (traitCode string, err error) {
	if avatar == "" {
		return
	}
	anchorArr := strings.Split(avatar, "/")
	if len(anchorArr) > 0 {
		anchorImgArr := strings.Split(anchorArr[len(anchorArr)-1], ".")
		if len(anchorImgArr) > 0 {
			traitCode = anchorImgArr[0]
		}
	}
	return
}
