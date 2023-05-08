package utils

import (
	"context"
	"github.com/youchuangcd/gopkg"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type MyUUID struct {
	Value string
}

func (p MyUUID) String() string {
	return p.Value
}

// GenToken
// @Description: 生成token
// @param userId
// @return string
func GenToken(userId uint) string {
	return MD5V([]byte(RandSeq(30) + strconv.FormatUint(uint64(userId), 10)))
}

// GenPassword
// @Description: 生成密码
// @param password
// @param salt
// @return string
func GenPassword(password, salt string) string {
	return MD5V([]byte(MD5V([]byte(password)) + salt))
}

// FormatMobileStar
// @Description: 手机号中间4位替换为*号
// @param mobile
// @return string
func FormatMobileStar(mobile string) string {
	if len(mobile) <= 10 {
		return mobile
	}
	return mobile[:3] + "****" + mobile[7:]
}

// 字符串切片去重
func RemoveRepeatedStr(s []string) []string {
	var result []string
	m := make(map[string]struct{}) //map的值不重要
	for _, v := range s {
		if _, ok := m[v]; !ok {
			result = append(result, v)
			m[v] = struct{}{}
		}
	}
	return result
}

// UUID
// @Description: 获取UUID 封装是为了方便替换
// @return uid
// @return err
var (
	macAddr, _ = GetMac()
	uuidCount  = uuidCountParam{
		initStr: macAddr + strconv.FormatInt(time.Now().UnixNano(), 10) + strconv.FormatInt(int64(os.Getpid()), 10), // 降低多节点同时获取id可能出现重复的概率
	}
)

type uuidCountParam struct {
	count   atomic.Uint64
	initStr string
}

func UUID() (uid MyUUID, err error) {
	// 计数器值 + 19位时间戳 纳秒级
	uid = MyUUID{
		Value: MD5V([]byte(uuidCount.initStr + strconv.FormatUint(uuidCount.count.Add(1), 10) + strconv.FormatInt(time.Now().UnixNano(), 10))),
	}
	return uid, nil
}

func GenUniqueId() string {
	uid, _ := UUID()
	return uid.String()
}

// IntStr2Uint
// @Description: 整形字符串转Uint
// @param s
// @return res
func IntStr2Uint(s string) (res uint) {
	if s != "" {
		r := []rune(s)
		length := len(r) - 1
		intStr := r[:length]
		lastStr := r[length:] // 最后一个字符
		var (
			fv  float64
			uv  uint64
			err error
		)
		switch string(lastStr) {
		case "w":
			fv, err = strconv.ParseFloat(string(intStr), 64)
			if err == nil {
				fv *= 10000
			}
			res = uint(fv)
		case "亿":
			fv, err = strconv.ParseFloat(string(intStr), 64)
			if err == nil {
				fv *= 100000000
			}
			res = uint(fv)
		default:
			uv, err = strconv.ParseUint(s, 10, 0)
			res = uint(uv)
		}
	}
	return
}

// AgeStr2Uint
// @Description: 年龄字符串转整形
// @param s
// @return res
func AgeStr2Uint(s string) (res uint) {
	if s != "" {
		r := []rune(s)
		length := len(r) - 1
		intStr := r[:length]
		lastStr := r[length:] // 最后一个字符
		var (
			uv uint64
		)
		switch string(lastStr) {
		case "岁":
			uv, _ = strconv.ParseUint(string(intStr), 10, 0)
			res = uint(uv)
		default:
			uv, _ = strconv.ParseUint(s, 10, 0)
			res = uint(uv)
		}
	}
	return
}

// JoinEnvStr
// @Description: 拼接环境变量，默认以下划线加环境变量拼接在后面; eg: str_dev
// @param s
// @param extArgs 指定拼接符 JoinEnvStr("xxx", "环境变量"); JoinEnvStr("xxx", "环境变量", "拼接符")
// @return string
func JoinEnvStr(s string, extArgs ...string) string {
	spliceSymbol := "_"
	env := gopkg.Env
	if len(extArgs) == 1 {
		env = extArgs[0]
	}
	if len(extArgs) == 2 {
		env = extArgs[0]
		spliceSymbol = extArgs[1]
	}
	return s + spliceSymbol + env
}

// TrimEnvStr
// @Description: 修剪字符串后缀，移除环境变量字符串； str_dev => str
// @param s
// @param extArgs
// @return string
func TrimEnvStr(s string, extArgs ...string) string {
	spliceSymbol := "_"
	env := gopkg.Env
	if len(extArgs) == 1 {
		env = extArgs[0]
	}
	if len(extArgs) == 2 {
		env = extArgs[0]
		spliceSymbol = extArgs[1]
	}
	return strings.TrimSuffix(s, spliceSymbol+env)
}

// GenTraceId
// @Description: 生成TraceId
// @param ctx
// @return string
func GenTraceId(ctx context.Context) string {
	//return strings.Replace(time.Now().Format("20060102150405.000"), ".", "", 1) + RandSeq(23)
	// 解决并发traceId冲突的问题
	return strings.Replace(time.Now().Format("20060102150405.000"), ".", "", 1) + GenUniqueId()
}

// AmountToFloat 金额字符串转浮点
func AmountToFloat(amount string) float64 {
	strAmount := strings.Replace(amount, "￥", "", 1)
	floatAmount, err := strconv.ParseFloat(strAmount, 64)
	if err == nil {
		return floatAmount
	}
	return 0.0
}

// CutStrFromLogConfig
// @Description: 根据日志配置截取字符长度
// @param s
// @return string
func CutStrFromLogConfig(s string) string {
	return CutStr(s, gopkg.LogLimitContentLength, gopkg.LogLimitContentReplaceWord)
}

// CutStr
// @Description: 截取字符串
// @param s
// @param limit 不支持负数
// @param rs 中间替换的内容，如...
// @return string
func CutStr(s string, limit uint, rs string) string {
	runeStr := []rune(s)
	sl := len(runeStr)
	if sl > int(limit) {
		halfLen := limit / 2
		var buff strings.Builder
		buff.WriteString(string(runeStr[0:halfLen]))
		buff.WriteString(rs)
		buff.WriteString(string(runeStr[sl-int(halfLen):]))
		return buff.String()
	}
	return s
}

// 字符串在指定的列数处换行文本
// s string 字符串文本数据
// limit int 限制单词个数
func WordWrap(s string, limit int) string {
	if strings.TrimSpace(s) == "" {
		return s
	}
	//将字符串转换为切片
	strSlice := strings.Fields(s)
	var result string = ""
	for len(strSlice) >= 1 {
		// 满足最后几个单词
		if len(strSlice) < limit {
			limit = len(strSlice)
		}
		// 将切片/数组转换回字符串
		// 但在指定的限制处插入 \r\n
		result = result + strings.Join(strSlice[:limit], " ") + "\r\n"
		//丢弃复制到结果的元素
		strSlice = strSlice[limit:]
	}

	return result
}

// 生成随机字符串
func RandStr(n int) string {
	// 保证每次生成的随机数不一样
	rand.Seed(time.Now().UnixNano())
	// 长度40
	bytes := []byte("~!@$%^&*()-=_+ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	result := make([]byte, n)
	for i := 0; i < n; i++ {
		result[i] = bytes[rand.Int31()%40]
	}
	return string(result)
}

// 判断字符串是不是纯数字
func IsNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func FormatStar(s string, args ...int) string {
	starNum := 4
	if len(args) > 0 && args[0] > 1 {
		starNum = args[0]
	}
	ns := []rune(s)
	l := len(ns)
	if l < starNum+2 {
		return s
	}
	half := l / 2
	starHalf := starNum / 2
	var sb strings.Builder
	sb.WriteString(string(ns[:half-starHalf]))
	sb.WriteString(strings.Repeat("*", starNum))
	sb.WriteString(string(ns[half+starHalf:]))
	//return string(ns[:half-starHalf]) + strings.Repeat("*", starNum) + string(ns[half + starHalf:])
	return sb.String()
}
