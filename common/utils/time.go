package utils

import (
	"github.com/youchuangcd/gopkg"
	"time"
)

// IsNowInTimeRange
// @Description: 当前时间是否在指定范围内
// @param startTimeStr 参数为时间字符串，格式为"时:分:秒"
// @param endTimeStr 参数为时间字符串，格式为"时:分:秒"
// @return bool
func IsNowInTimeRange(now time.Time, startTimeStr, endTimeStr string) bool {
	//当前时间转换为"年-月-日"的格式
	format := now.Format(gopkg.DateFormat)
	//转换为time类型需要的格式
	layout := gopkg.DateTimeFormat
	//将开始时间拼接“年-月-日 ”转换为time类型
	timeStart, err := time.ParseInLocation(layout, format+" "+startTimeStr, time.Local)
	if err != nil {
		return false
	}
	//将结束时间拼接“年-月-日 ”转换为time类型
	timeEnd, err := time.ParseInLocation(layout, format+" "+endTimeStr, time.Local)
	if err != nil {
		return false
	}
	// 如果结束时间是0点0分0秒，就需要加一天
	//if endTimeStr == "00:00:00" {
	//	timeEnd = timeEnd.AddDate(0,0,1)
	//}
	//使用time的Before和After方法，判断当前时间是否在参数的时间范围
	return now.Before(timeEnd) && now.After(timeStart)
}

func ConvertTimeStrToTime(now time.Time, timeStr string) (t time.Time, err error) {
	//当前时间转换为"年-月-日"的格式
	format := now.Format(gopkg.DateFormat)
	//转换为time类型需要的格式
	layout := gopkg.DateTimeFormat
	//将开始时间拼接“年-月-日 ”转换为time类型
	t, err = time.ParseInLocation(layout, format+" "+timeStr, time.Local)
	return
}

// GetPeriodValueByTime
// @Description: 指定多少天为一个周期，计算今天在周期中的位置
// @param now
// @param periodDayNum
// @return int64
func GetPeriodValueByTime(now time.Time, periodDayNum int64) int64 {
	if periodDayNum <= 0 {
		return 0
	}
	// 取当天零点的日期时间戳
	zeroTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	day := zeroTime.Unix() / 86400
	return day % periodDayNum
}

// GetDiffDays
// @Description: 获取两个时间相差的天数，0表同一天，正数表t1>t2，负数表t1<t2
// @param t1
// @param t2
// @return int
func GetDiffDays(t1, t2 time.Time) int {
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, time.Local)
	t2 = time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, time.Local)

	return int(t1.Sub(t2).Hours() / 24)
}
