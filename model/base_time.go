package model

import (
	"database/sql/driver"
	"fmt"
	"github.com/youchuangcd/gopkg"
	"gorm.io/gorm/schema"
	"strconv"
	"strings"
	"time"
)

// LocalTime json 格式化为 时间戳
type LocalTime struct {
	time.Time
}

// LocalDateTime json 格式化为 2006-01-02 15:04:05格式
type LocalDateTime struct {
	time.Time
}

// LocalDate json 格式化为 2006-01-02格式
type LocalDate struct {
	time.Time
}

type UnixTime2DateTime struct {
	time.Time
}

// LocalDateMsTime json 格式化为 2006-01-02 15:04:05.000格式
type LocalDateMsTime struct {
	time.Time
}

// MarshalJSON
// @Description: json编码时，格式化为时间戳
// @receiver t
// @return []byte
// @return error
func (t LocalTime) MarshalJSON() ([]byte, error) {
	//格式化秒
	seconds := t.Unix()
	return []byte(strconv.FormatInt(seconds, 10)), nil
}

// UnmarshalJSON
// @Description: 把时间戳转化为LocalTime对象
// @receiver t
// @param data
// @return error
func (t *LocalTime) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}
	i, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	// Fractional seconds are handled implicitly by Parse.
	t.Time = time.Unix(i, 0)
	return nil
}

// Value
// @Description:  存储到数据库时取值
// @receiver t
// @return driver.Value
// @return error
func (t LocalTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

//func (t LocalTime) GormDataType() string {
//	return string(schema.Time)
//}

// Scan
// @Description:
// @receiver t
// @param v
// @return error
func (t *LocalTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = LocalTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

// MarshalJSON
// @Description: json编码时，格式化为指定格式字符串
// @receiver t
// @return []byte
// @return error
func (t LocalDateTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%v\"", t.Format(gopkg.DateTimeFormat))
	return []byte(formatted), nil
}

// UnmarshalJSON
// @Description: 把时间字符串转化为LocalTime对象
// @receiver t
// @param data
// @return error
func (t *LocalDateTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	var err error
	//前端接收的时间字符串
	str := string(data)
	//去除接收的str收尾多余的"
	timeStr := strings.Trim(str, "\"")
	t1, err := time.Parse(gopkg.DateTimeFormat, timeStr)
	*t = LocalDateTime{Time: t1}
	return err
}

// Value
// @Description:  存储到数据库时取值
// @receiver t
// @return driver.Value
// @return error
func (t LocalDateTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t LocalDateTime) GormDataType() string {
	return string(schema.Time)
}

// Scan
// @Description:
// @receiver t
// @param v
// @return error
func (t *LocalDateTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = LocalDateTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to datetime", v)
}

// MarshalJSON
// @Description: json编码时，格式化为指定格式字符串
// @receiver t
// @return []byte
// @return error
func (t LocalDate) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%v\"", t.Format(gopkg.DateFormat))), nil
}

// UnmarshalJSON
// @Description: 把时间字符串转化为LocalTime对象
// @receiver t
// @param data
// @return error
func (t *LocalDate) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	var err error
	//前端接收的时间字符串
	str := string(data)
	//去除接收的str收尾多余的"
	timeStr := strings.Trim(str, "\"")
	t1, err := time.Parse(gopkg.DateFormat, timeStr)
	*t = LocalDate{Time: t1}
	return err
}

// Value
// @Description:  存储到数据库时取值
// @receiver t
// @return driver.Value
// @return error
func (t LocalDate) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t LocalDate) GormDataType() string {
	return string(schema.Time)
}

// Scan
// @Description:
// @receiver t
// @param v
// @return error
func (t *LocalDate) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = LocalDate{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to date", v)
}

// MarshalJSON
// @Description: json编码时，格式化为指定格式字符串
// @receiver t
// @return []byte
// @return error
func (t UnixTime2DateTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%v\"", t.Format(gopkg.DateTimeFormat))
	return []byte(formatted), nil
}

// UnmarshalJSON
// @Description: 把时间字符串转化为LocalTime对象
// @receiver t
// @param data
// @return error
func (t *UnixTime2DateTime) UnmarshalJSON(data []byte) error {
	unixTime, err := strconv.ParseInt(string(data), 10, 0)
	if err == nil {
		*t = UnixTime2DateTime{Time: time.Unix(unixTime, 0)}
	}
	return err
}

func (t UnixTime2DateTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan
// @Description:
// @receiver t
// @param v
// @return error
func (t *UnixTime2DateTime) Scan(v interface{}) error {
	value, ok := v.(int64)
	if ok {
		*t = UnixTime2DateTime{Time: time.Unix(value, 0)}
		return nil
	}
	return fmt.Errorf("can not convert %v to UnixTime2DateTime", v)
}

// RedisArg
// @Description: redis扫描对象到hash参数处理
// @receiver t
// @return interface{}
func (t LocalDateTime) RedisArg() interface{} {
	return t.Format(gopkg.DateTimeFormat)
}

// RedisScan
// @Description: 读取hash等字段时，映射成对应的时间结构
// @receiver t
// @param x
// @return error
func (t *LocalDateTime) RedisScan(x interface{}) error {
	bs, ok := x.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte, got %T", x)
	}
	tt, err := time.Parse(gopkg.DateTimeFormat, string(bs))
	if err != nil {
		return err
	}
	*t = LocalDateTime{Time: tt}
	return nil
}

// RedisArg
// @Description: redis扫描对象到hash参数处理
// @receiver t
// @return interface{}
func (t LocalTime) RedisArg() interface{} {
	unix := t.Unix()
	if unix < 0 {
		unix = 0
	}
	return unix
}

// RedisScan
// @Description: 读取hash等字段时，映射成对应的时间结构
// @receiver t
// @param x
// @return error
func (t *LocalTime) RedisScan(x interface{}) error {
	bs, ok := x.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte, got %T", x)
	}
	i, err := strconv.ParseInt(string(bs), 10, 64)
	if err != nil {
		return err
	}
	tt := time.Unix(i, 0)
	*t = LocalTime{Time: tt}
	return nil
}

// RedisArg
// @Description: redis扫描对象到hash参数处理
// @receiver t
// @return interface{}
func (t LocalDate) RedisArg() interface{} {
	return t.Format(gopkg.DateFormat)
}

// RedisScan
// @Description: 读取hash等字段时，映射成对应的时间结构
// @receiver t
// @param x
// @return error
func (t *LocalDate) RedisScan(x interface{}) error {
	bs, ok := x.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte, got %T", x)
	}
	tt, err := time.ParseInLocation(gopkg.DateFormat, string(bs), time.Local)
	if err != nil {
		return err
	}
	*t = LocalDate{Time: tt}
	return nil
}

// MarshalJSON
// @Description: json编码时，格式化为指定格式字符串
// @receiver t
// @return []byte
// @return error
func (t LocalDateMsTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%v\"", t.Format(gopkg.DateMsTimeFormat))
	return []byte(formatted), nil
	//unixMilli := t.Unix()*1e3 + (t.UnixNano()-t.Unix()*1e9)/1e6 // 等同于go 1.17 time.UnixMilli()
	//return []byte(strconv.FormatInt(unixMilli, 10)), nil
}

// UnmarshalJSON
// @Description: 把时间字符串(2006-01-02 15:04:05.999)转化为LocalDateMsTime对象
// @receiver t
// @param data
// @return error
func (t *LocalDateMsTime) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}
	//前端接收的时间字符串
	str := string(data)
	//去除接收的str收尾多余的"
	timeStr := strings.Trim(str, "\"")
	t1, err := time.Parse(gopkg.DateMsTimeFormat, timeStr)
	if err != nil {
		return err
	}
	*t = LocalDateMsTime{Time: t1}
	//msec, err := strconv.ParseInt(string(data), 10, 64)
	//if err != nil {
	//	return err
	//}
	//t.Time = time.Unix(msec/1e3, (msec%1e3)*1e6)
	//t.Time = time.UnixMilli(i) // go 1.17才支持
	return nil

}

// Value
// @Description: 存储到数据库时取值
// @receiver t
// @return driver.Value
// @return error
func (t LocalDateMsTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t LocalDateMsTime) GormDataType() string {
	return string(schema.Time)
}

// Scan
// @Description:
// @receiver t
// @param v
// @return error
func (t *LocalDateMsTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = LocalDateMsTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to datetime", v)
}
