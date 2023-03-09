package redis

import (
	"context"
	"github.com/gomodule/redigo/redis"
)

// 设置值
func Set(ctx context.Context, key string, value interface{}, expire ...int64) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	if len(expire) == 0 {
		return getReply(redis.DoContext(c, ctx, "set", key, value))
	}
	return getReply(redis.DoContext(c, ctx, "set", key, value, "ex", expire[0]))
}

// 获取值
func Get(ctx context.Context, key string) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "get", key))
}

// key不存在是在设置值
func SetNX(ctx context.Context, key string, value interface{}) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "setnx", key, value))
}

// 设置并返回旧值
func GetSet(ctx context.Context, key string, value interface{}) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "getset", key, value))
}

// 设置key并指定生存时间
func SetEX(ctx context.Context, key string, value interface{}, seconds int64) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "setex", key, seconds, value))
}

// 设置key值并指定生存时间(毫秒)
func PSetEX(ctx context.Context, key string, value interface{}, milliseconds int64) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "psetex", key, milliseconds, value))
}

// 设置子字符串
func SetRange(ctx context.Context, key string, value interface{}, offset int64) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "setrange", key, offset, value))
}

// 获取子字符串
func GetRange(ctx context.Context, key string, start, end int64) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "getrange", key, start, end))
}

// 设置多个值
func MSet(ctx context.Context, kv map[string]interface{}) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "mset", redis.Args{}.AddFlat(kv)))
}

// key不存在时设置多个值
func MSetNx(ctx context.Context, kv map[string]interface{}) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "msetnx", redis.Args{}.AddFlat(kv)))
}

// 返回多个key的值
func MGet(ctx context.Context, keys []string) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "mget", redis.Args{}.AddFlat(keys)...))
}

// Incr
// @Description: 自增
// @param key
// @param args 过期时间
// @return *Reply
func Incr(ctx context.Context, key string, args ...int) *Reply {
	if len(args) == 1 {
		return EvalScript(ctx, ScriptKeyIncr, key, args[0])
	}
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "incr", key))
}

// IncrBy
// @Description: 增加指定值
// @param key
// @param increment
// @param args 过期时间
// @return *Reply
func IncrBy(ctx context.Context, key string, increment int64, args ...int) *Reply {
	if len(args) == 1 {
		return EvalScript(ctx, ScriptKeyIncrBy, key, increment, args[0])
	}
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "incrby", key, increment))
}

// 增加一个浮点值
func IncrByFloat(ctx context.Context, key string, increment float64) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "incrbyfloat", key, increment))
}

// 自减
func Decr(ctx context.Context, key string) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "decr", key))
}

// 自减指定值
func DecrBy(ctx context.Context, key string, increment int64) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "decrby", key, increment))
}

// IncrReset
// @Description: 增量计数器，并在增量值超过最大值时，重置为0;
// @param key
// @param max
// @param expire
// @return res
// @return err
func IncrReset(ctx context.Context, key string, max int, expire int) (res int, err error) {
	return EvalScript(ctx, ScriptKeyIncrReset, key, max, expire).Int()
}

// IncrByOutMaxReset
// @Description: 增量计数器，并在超出最大值时，重置为0；或主动重置为0;
// @param key
// @param increment 本次要加的数量
// @param max 最大值
// @param isResetBool 本次是否要重置
// @param resetVal 重置的值
// @return resIsReset
// @return currentValue
// @return err
func IncrByOutMaxReset(ctx context.Context, key string, increment int, max int, isResetBool bool, resetVal int) (resIsReset bool, currentValue int, err error) {
	isReset := 0
	if isResetBool {
		isReset = 1
	}
	res, err := EvalScript(ctx, ScriptKeyIncrByOutMaxReset, key, increment, max, isReset, resetVal).Values()
	if err != nil {
		return
	}
	if v, ok := res[0].(int64); ok && v == 1 {
		resIsReset = true
	}
	if v, ok := res[1].(int64); ok {
		currentValue = int(v)
	}
	return
}

// IncrMax
// @Description: 增量计数器，如果当前值没有大于限定值，才可以加一并返回累加后的值，否则返回当前值
// @param key
// @param max
// @param expire
// @return outMax 是否达到最大值
// @return currentValue 当前值
// @return err
func IncrMax(ctx context.Context, key string, max int, expire int) (outMax bool, currentValue int, err error) {
	res, err := EvalScript(ctx, ScriptKeyIncrMax, key, max, expire).Values()
	if err != nil {
		return
	}
	if v, ok := res[0].(int64); ok && v == 0 {
		outMax = true
	}
	if v, ok := res[1].(int64); ok {
		currentValue = int(v)
	}
	return
}
