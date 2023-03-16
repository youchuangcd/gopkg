package model

import (
	"context"
	"fmt"
	"github.com/youchuangcd/gopkg"
	"github.com/youchuangcd/gopkg/mylog"
	"github.com/youchuangcd/gopkg/redis"
	"strconv"
	"time"
)

// Timestamp
// @Description: 针对redis time.time的支持
type Timestamp struct {
	time.Time
}

// RedisArg
// @Description: redis扫描对象到hash参数处理
// @receiver t
// @return interface{}
func (t Timestamp) RedisArg() interface{} {
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
func (t *Timestamp) RedisScan(x interface{}) error {
	bs, ok := x.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte, got %T", x)
	}
	i, err := strconv.ParseInt(string(bs), 10, 64)
	if err != nil {
		return err
	}
	tt := time.Unix(i, 0)
	*t = Timestamp{Time: tt}
	return nil
}

func DelKey(ctx context.Context, key ...string) (err error) {
	err = redis.Del(ctx, key...).Error()
	if err != nil {
		logContent := map[string]interface{}{
			"err": err,
			"key": key,
		}
		mylog.WithError(ctx, gopkg.LogRedis, logContent, "删除key失败")
	}
	return
}

// HIncrCounter
// @Description: hash int计数器累加
// @param ctx
// @param key
// @param field
// @param args
// @return res
// @return err
func HIncrCounter(ctx context.Context, key string, field string, args ...int) (res int, err error) {
	var increment = 1
	if len(args) > 0 {
		increment = args[0]
	}
	res, err = redis.HIncrBy(ctx, key, field, increment).Int()
	if err != nil && err != redis.ErrNil {
		logContent := map[string]interface{}{
			"err": err,
			"key": key,
		}
		mylog.WithError(ctx, gopkg.LogRedis, logContent, "hash计数器累加失败")
	}
	return
}

// HIncrByFloatCounter
// @Description: hash 浮点数计数器累加
// @param ctx
// @param key
// @param field
// @param args
// @return res
// @return err
func HIncrByFloatCounter(ctx context.Context, key string, field string, args ...float64) (res float64, err error) {
	var increment = 1.0
	if len(args) > 0 {
		increment = args[0]
	}
	res, err = redis.HIncrByFloat(ctx, key, field, increment).Float64()
	if err != nil && err != redis.ErrNil {
		logContent := map[string]interface{}{
			"err": err,
			"key": key,
		}
		mylog.WithError(ctx, gopkg.LogRedis, logContent, "hash float计数器累加失败")
	}
	return
}

// HGetString
// @Description: hash 获取指定域值失败
// @param ctx
// @param key
// @param field
// @return res string
// @return err
func HGetString(ctx context.Context, key string, field string) (res string, err error) {
	res, err = redis.HGet(ctx, key, field).String()
	if err != nil && err != redis.ErrNil {
		logContent := map[string]interface{}{
			"err": err,
			"key": key,
		}
		mylog.WithError(ctx, gopkg.LogRedis, logContent, "hash获取指定域值失败")
	}
	return
}

// HMGetString
// @Description:
// @param ctx
// @param key
// @param fields
// @return res
// @return err
func HMGetString(ctx context.Context, key string, fields []string) (res map[string]string, err error) {
	r, err := redis.HMGet(ctx, key, fields).Strings()
	if err != nil {
		if err != redis.ErrNil {
			logContent := map[string]interface{}{
				"err":   err,
				"key":   key,
				"param": fields,
			}
			mylog.WithError(ctx, gopkg.LogRedis, logContent, "hmget获取参数域失败")
		}
		return nil, err
	}
	res = make(map[string]string, len(fields))
	// 生成返回值field => value map
	for k, v := range fields {
		res[v] = r[k]
	}
	return
}

// HGetInt
// @Description:
// @param ctx
// @param key
// @param field
// @return res
// @return err
func HGetInt(ctx context.Context, key string, field string) (res int, err error) {
	res, err = redis.HGet(ctx, key, field).Int()
	if err != nil && err != redis.ErrNil {
		logContent := map[string]interface{}{
			"err": err,
			"key": key,
		}
		mylog.WithError(ctx, gopkg.LogRedis, logContent, "hash获取指定域值失败")
	}
	return
}

// HMGetInts
// @Description: 批量获取多个域的值转换成整数值
// @param ctx
// @param key
// @param fields
// @return res
// @return err
func HMGetInts(ctx context.Context, key string, fields []string) (res map[string]int, err error) {
	r, err := redis.HMGet(ctx, key, fields).Ints()
	if err != nil {
		if err != redis.ErrNil {
			logContent := map[string]interface{}{
				"err":    err,
				"key":    key,
				"fields": fields,
			}
			mylog.WithError(ctx, gopkg.LogRedis, logContent, "hash获取指定域值失败")
		}
		return nil, err
	}
	res = make(map[string]int, len(fields))
	// 生成返回值field => value map
	for k, v := range fields {
		res[v] = r[k]
	}
	return
}

// GetInt
// @Description: 获取int类型的值
// @param ctx
// @param key
// @return res
// @return err
func GetInt(ctx context.Context, key string) (res int, err error) {
	res, err = redis.Get(ctx, key).Int()
	if err != nil {
		if err != redis.ErrNil {
			logContent := map[string]interface{}{
				"err": err,
				"key": key,
			}
			mylog.WithError(ctx, gopkg.LogRedis, logContent, "获取int类型的key值失败")
		}
	}
	return
}
