package utils

import (
	"fmt"
)

/*
时间常量
*/
const (
	//定义每分钟的秒数
	SecondsPerMinute = 60
	//定义每小时的秒数
	SecondsPerHour = SecondsPerMinute * 60
	//定义每天的秒数
	SecondsPerDay = SecondsPerHour * 24
)

/*
*时间转换函数秒数转换为 分钟:秒
 */
func ResolveTimeMinuteSeconds(s int) (minute, seconds int) {
	//
	seconds = s % SecondsPerMinute
	//每分钟秒数
	minute = s / SecondsPerMinute
	return
}

func ResolveTimeMinuteSecondsStr(s int) string {
	minute, seconds := ResolveTimeMinuteSeconds(s)
	args := make([]interface{}, 0)
	// 默认显示分:秒 不足两位补0
	var format string = "%02d:%02d"
	args = append(args, minute, seconds)
	return fmt.Sprintf(format, args...)
}
