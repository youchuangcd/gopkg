package model

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/youchuangcd/gopkg"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
	"time"
)

// gorm.Model 的定义
type GormModel struct {
	ID        uint      `json:"id" gorm:"primaryKey" redis:"id" mapstructure:"id"`
	CreatedAt LocalTime `json:"created_at" gorm:"<-:false" redis:"created_at" mapstructure:"created_at"` // allow read, disable write permission
	UpdatedAt LocalTime `json:"-" gorm:"<-:false" redis:"-"`                                             // allow read, disable write permission
}

type GormModelIntID struct {
	ID        int       `json:"id" gorm:"primaryKey" redis:"id" mapstructure:"id"`
	CreatedAt LocalTime `json:"created_at" gorm:"<-:false" redis:"created_at" mapstructure:"created_at"` // allow read, disable write permission
	UpdatedAt LocalTime `json:"-" gorm:"<-:false" redis:"-"`                                             // allow read, disable write permission
}

// gorm.Model 的定义
type GormJsonHiddenTimeModel struct {
	ID        uint      `json:"id" gorm:"primaryKey" mapstructure:"id" redis:"id"`
	CreatedAt time.Time `json:"-" gorm:"<-:false" redis:"created_at"` // allow read, disable write permission
	UpdatedAt time.Time `json:"-" gorm:"<-:false"`                    // allow read, disable write permission
}

type GormModelTime struct {
	ID          uint  `json:"id" gorm:"primaryKey" mapstructure:"id" redis:"id"`
	CreatedTime int64 `json:"created_time" redis:"created_time" gorm:"autoCreateTime"` // allow read, disable write permission
	UpdatedTime int64 `json:"updated_time" redis:"-" gorm:"autoUpdateTime"`            // allow read, disable write permission
}

// gorm.Model 的定义
type GormModelAt struct {
	ID        uint            `json:"id" gorm:"primaryKey" redis:"id"`
	CreatedAt LocalDateMsTime `json:"created_at" gorm:"autoCreateTime" redis:"created_at"`
	UpdatedAt LocalDateMsTime `json:"updated_at" gorm:"autoUpdateTime" redis:"-"`
}

// gorm.Model dws的定义
type GormModelAtDws struct {
	ID        uint            `json:"id" gorm:"autoIncrement:true" redis:"id"`
	CreatedAt LocalDateMsTime `json:"created_at" gorm:"autoCreateTime" redis:"created_at"`
	UpdatedAt LocalDateMsTime `json:"updated_at" gorm:"autoUpdateTime" redis:"-"`
}

type Interface interface {
	TableName() string
}

type JsonValue struct {
	Val     interface{}
	valByte []byte // 存储原json字节，用来后续反序列化到其他类型上
}

func UpdateAddUpdatedTime(ctx context.Context, tx *gorm.DB) {
	tx.Statement.SetColumn("UpdatedTime", time.Now().Unix())
	if tx.Statement.Selects != nil {
		tx.Statement.Selects = append(tx.Statement.Selects, "updated_time")
	}
}

// GetDB
// @Description: 获取连接句柄；默认使用datalake数据库
// @param ctx
// @param tx
// @param args 指定数据库key名 @see global.GormDBMapKeyDataLake or global.GormDBMapKeyKfpt
// @return db
func GetDB(ctx context.Context, tx *gorm.DB, args ...interface{}) (db *gorm.DB) {
	if tx != nil {
		db = tx
	} else {
		gormDBMapKey := gopkg.GormDBMapKeyDefault
		if len(args) > 0 {
			if key, ok := args[0].(string); ok {
				gormDBMapKey = key
			}
		} else if key, ok := ctx.Value(gopkg.ContextDBMapKey).(string); ok { // 从上下文里切换
			gormDBMapKey = key
		}
		newCtx := ctx
		// http请求的话，要提取request里面的上下文才可以获取到b3请求头
		//if ginCtx, ok := ctx.Value(gin.ContextKey).(*gin.Context); ok {
		//	newCtx = ginCtx.Request.Context()
		//}
		db = gopkg.GormDBMap[gormDBMapKey].WithContext(newCtx)
	}
	return
}

// CounterColumn
// @Description: 自定义数据类型; 计数器的列，用SQL 表达式更新列;
// @link https://gorm.io/zh_CN/docs/data_types.html#gorm_valuer_interface
type CounterColumn struct {
	FieldName string // 默认 updated_num
	Value     int
}

// Scan 方法实现了 sql.Scanner 接口
func (p *CounterColumn) Scan(v interface{}) error {
	// Scan a value into struct from database driver
	if val, ok := v.(int); ok {
		p.Value = val
	}
	return nil
}

// GormDataType
// @Description: 告诉gorm数据库中的字段类型
// @receiver p
// @return string
func (p CounterColumn) GormDataType() string {
	return "int"
}

// GormValue
// @Description:
// @receiver p
// @param ctx
// @param db
// @return clause.Expr
func (p CounterColumn) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	field := "updated_num"
	if p.FieldName != "" {
		field = p.FieldName
	}
	sql := field + " + ?"                         // eg: "updated_num + ?"
	if db.Config.Dialector.Name() == "postgres" { // NOTE: 兼容postgres
		sql = strconv.FormatUint(uint64(time.Now().UnixNano()), 10) + " + ?"
	}
	return clause.Expr{
		SQL:  sql, // eg: "updated_num + ?"
		Vars: []interface{}{p.Value},
	}
}

// MarshalJSON
// @Description: json编码时，格式化为字段值
// @receiver t
// @return []byte
// @return error
func (p CounterColumn) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(p.Value), 10)), nil
}

// UnmarshalJSON
// @Description: json解码时，把值赋值到字段的value上
// @receiver t
// @param data
// @return error
func (p *CounterColumn) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}
	i, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	p.Value = int(i)
	return nil
}

// Scan 方法实现了 sql.Scanner 接口; 把数据库字段值扫描到结构体上
func (p *JsonValue) Scan(v interface{}) error {
	if s, ok := v.([]byte); ok {
		return json.Unmarshal(s, p)
	}
	return fmt.Errorf("can not convert %v to JsonValue", v)
}

// Value
// @Description: 把结构体值json存到数据库字段中
// @receiver p
// @return driver.Value
// @return error
func (p JsonValue) Value() (driver.Value, error) {
	b, err := json.Marshal(p)
	return string(b), err
}

func (p JsonValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Val)
}

// UnmarshalJSON
// @Description:
// @receiver t
// @param data
// @return error
func (p *JsonValue) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	tmp := JsonValue{}
	// 这里可以省略，直接不反序列化，直接把data存储到ValByte中，在后续逻辑中再反序列化到对应的结构体中
	err := json.Unmarshal(data, &tmp.Val)
	if err != nil {
		return err
	}
	tmp.valByte = data
	*p = tmp
	return nil
}

// GetValByte
//
//	@Description: 获取未反序列化原内容
//	@receiver p
//	@return []byte
func (p *JsonValue) GetValByte() []byte {
	return p.valByte
}
