package gopkg

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"log"
	"reflect"
	"strings"
	"time"
)

var (
	timeType = reflect.TypeOf(time.Time{})
	// Validate 验证器 validator.Validate是线程安全的，其变量内会缓存已经验证过结构体的特征，因此用户用一个变量更有利于提高效率
	Trans ut.Translator
)

// InitTrans
// @Description: 初始化验证器翻译
// @param locale
// @return err
func InitTrans(locale string) {
	var err error
	// 修改gin框架中的Validator引擎属性，实现自定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册一个获取json tag的自定义方法
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
		//注册翻译器
		zhT := zh.New() // 中文翻译器
		enT := en.New() // 英文翻译器

		// 第一个参数是备用（fallback）的语言环境
		// 后面的参数是应该支持的语言环境（支持多个）
		// uni := ut.New(zhT, zhT) 也是可以的
		uni := ut.New(enT, zhT, enT)

		// locale 通常取决于 http 请求头的 'Accept-Language'
		var ok bool
		// 也可以使用 uni.FindTranslator(...) 传入多个locale进行查找
		Trans, ok = uni.GetTranslator(locale)
		if !ok {
			log.Panic(fmt.Errorf("uni.GetTranslator(%s) failed", locale))
		} else {
			// 注册翻译器
			switch locale {
			case "en":
				err = enTranslations.RegisterDefaultTranslations(v, Trans)
			case "zh":
				err = zhTranslations.RegisterDefaultTranslations(v, Trans)
			default:
				err = enTranslations.RegisterDefaultTranslations(v, Trans)
			}
			if err != nil {
				log.Panic(err)
			}
		}
	}
}

// ValidatorFunc
// @Description: 初始化自定义验证规则
// @return err
func ValidatorFunc() {
	// 修改gin框架中的Validator引擎属性，实现自定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 时间字符串 "10:20:30" 时分秒
		_ = v.RegisterValidation("time", validatorFuncTime)
		// 大于其他时间字段 时分秒的时间字符串 "10:20:30"
		_ = v.RegisterValidation("gttimefield", isGtTimeField)
		// 大于等于其他时间字段 时分秒的时间字符串 "10:20:30"
		_ = v.RegisterValidation("gtetimefield", isGteTimeField)
	}
}

func validatorFuncTime(fl validator.FieldLevel) bool {
	if timeStr, ok := fl.Field().Interface().(string); ok {
		//当前时间
		now := time.Now()
		//当前时间转换为"年-月-日"的格式
		format := now.Format(DateFormat)
		//转换为time类型需要的格式
		layout := DateTimeFormat
		//将开始时间拼接“年-月-日 ”转换为time类型
		_, err := time.ParseInLocation(layout, format+" "+timeStr, time.Local)
		if err != nil {
			return false
		}
		return true
	}
	return false
}

func gtTimeField(fl validator.FieldLevel, isGteType bool) bool {
	field := fl.Field()
	kind := field.Kind()
	// currentField = 被比较的那个字段; 如 gttimefield=StartPeriod, currentField = StartPeriod
	currentField, currentKind, _, ok := fl.GetStructFieldOK2()
	if !ok || currentKind != kind {
		return false
	}

	switch kind {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

		return field.Int() > currentField.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:

		return field.Uint() > currentField.Uint()

	case reflect.Float32, reflect.Float64:

		return field.Float() > currentField.Float()

	case reflect.Struct:

		fieldType := field.Type()

		// Not Same underlying type i.e. struct and time
		if fieldType != currentField.Type() {
			return false
		}

		if fieldType == timeType {

			t := currentField.Interface().(time.Time)
			fieldTime := field.Interface().(time.Time)

			return fieldTime.After(t)
		}

	case reflect.String: // 字符串时间 10:20:10

		if field.Type() != currentField.Type() {
			return false
		}
		currentTimeStr := currentField.Interface().(string)
		fieldTimeStr := field.Interface().(string)

		//当前时间
		now := time.Now()
		//当前时间转换为"年-月-日"的格式
		format := now.Format(DateFormat)
		//转换为time类型需要的格式
		layout := DateTimeFormat
		//将开始时间拼接“年-月-日 ”转换为time类型
		timeStart, err := time.ParseInLocation(layout, format+" "+currentTimeStr, time.Local)
		if err != nil {
			return false
		}
		//将结束时间拼接“年-月-日 ”转换为time类型
		timeEnd, err := time.ParseInLocation(layout, format+" "+fieldTimeStr, time.Local)
		if err != nil {
			return false
		}
		//使用time的After方法，判断结束时间是否在参数的时间之后
		if isGteType {
			return timeEnd.After(timeStart) || timeEnd.Equal(timeStart)
		}
		return timeEnd.After(timeStart)
	}

	// default reflect.String
	return false
}

// isGtTimeField
// @Description: 是否大于其他时间字段 10:20:30 时:分:秒
// @param fl
// @return bool
func isGtTimeField(fl validator.FieldLevel) bool {
	return gtTimeField(fl, false)
}

// isGteTimeField
// @Description: 是否大于等于
// @param fl
// @return bool
func isGteTimeField(fl validator.FieldLevel) bool {
	return gtTimeField(fl, true)
}
