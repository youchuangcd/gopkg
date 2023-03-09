package redis

import (
	"context"
	"github.com/gomodule/redigo/redis"
)

// exist 为true 表示字段不存则设置其值
func HSet(ctx context.Context, key string, filed, value interface{}, exist ...bool) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	if len(exist) > 0 && exist[0] {
		return getReply(redis.DoContext(c, ctx, "hsetex", key, filed, value))
	}
	return getReply(redis.DoContext(c, ctx, "hset", key, filed, value))
}

// 获取指定字段值
func HGet(ctx context.Context, key string, filed interface{}) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "hget", key, filed))
}

// 获取所有字段及值
func HGetAll(ctx context.Context, key string) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "hgetall", key))
}

// 设置多个字段及值 [map]
func HMSetFromMap(ctx context.Context, key string, mp map[interface{}]interface{}) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "hmset", redis.Args{}.Add(key).AddFlat(mp)...))
}

// 设置多个字段及值 [struct]
func HMSetFromStruct(ctx context.Context, key string, obj interface{}) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "hmset", redis.Args{}.Add(key).AddFlat(obj)...))
}

// 返回多个字段值
func HMGet(ctx context.Context, key string, fields interface{}) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "hmget", redis.Args{}.Add(key).AddFlat(fields)...))
}

// 字段删除
func HDel(ctx context.Context, key string, fields interface{}) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "hdel", redis.Args{}.Add(key).AddFlat(fields)...))
}

// 判断字段是否存在
func HExists(ctx context.Context, key string, field interface{}) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "hexists", key, field))
}

// 返回所有字段
func HKeys(ctx context.Context, key string) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "hkeys", key))
}

// 返回字段数量
func HLen(ctx context.Context, key string) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "hlen", key))
}

// 返回所有字段值
func HVals(ctx context.Context, key string) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "hvals", key))
}

// 为指定字段值增加
func HIncrBy(ctx context.Context, key string, field interface{}, increment interface{}) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "hincrby", key, field, increment))
}

// 为指定字段值增加浮点数
func HIncrByFloat(ctx context.Context, key string, field interface{}, increment float64) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "hincrbyfloat", key, field, increment))
}

// HMIncrBy
// @Description: hash结构批量增加计数器(原子性)，返回每个域增加后的值
// @param key
// @param elems {"域1":要增加的值, "域2":要增加的值}
// @return map[string]int {"域1": 增加后的值, "域2": 增加后的值}
// @return error
func HMIncrBy(ctx context.Context, key string, elems map[string]int) (map[string]int, error) {
	var (
		args   = make([]interface{}, 0, len(elems)*2+1)
		fields = make([]string, 0, len(elems)) // 用来按传递参数的顺序记录域，拿到结果后按相同的顺序读取结果
	)
	args = append(args, key)
	for k, v := range elems {
		args = append(args, k, v)
		fields = append(fields, k)
	}
	res, err := EvalScript(ctx, ScriptKeyHMIncrBy, args...).Values()
	if err != nil {
		return nil, err
	}
	for k, field := range fields {
		elems[field], _ = redis.Int(res[k], nil)
	}
	return elems, nil
}
