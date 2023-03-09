package redis

import (
	"context"
	"github.com/gomodule/redigo/redis"
)

// 向列表头插入元素
func LPush(ctx context.Context, key string, value interface{}) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "lpush", key, value))
}

// 当列表存在则将元素插入表头
func LPushx(ctx context.Context, key string, value interface{}) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "lpushx", key, value))
}

// 将指定元素插入列表末尾
func RPush(ctx context.Context, key string, value interface{}) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "rpush", key, value))
}

// 当列表存在则将元素插入表尾
func RPushx(ctx context.Context, key string, value interface{}) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "rpushx", key, value))
}

// 将元素插入指定位置position:BEFORE|AFTER,当 pivot 不存在于列表 key 时，不执行任何操作。当 key 不存在时， key 被视为空列表，不执行任何操作。
func LInsert(ctx context.Context, key, position, pivot, value string) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "linsert", key, position, pivot, value))
}

// 返回列表头元素
func LPop(ctx context.Context, key string) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "lpop", key))
}

// 阻塞并弹出头元素
func BLpop(ctx context.Context, key, timeout interface{}) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "blpop", redis.Args{}.AddFlat(key).Add(timeout)...))
}

// 返回列表尾元素
func RPop(ctx context.Context, key string) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "rpop", key))
}

// 阻塞并弹出末尾元素
func BRpop(ctx context.Context, key, timeout interface{}) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "brpop", key, timeout))
}

// 返回指定位置的元素
func LIndex(ctx context.Context, key string, index interface{}) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "lindex", key, index))
}

// 获取指定区间的元素
func LRange(ctx context.Context, key string, start, stop interface{}) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "lrange", key, start, stop))
}

// 设置指定位元素
func LSet(ctx context.Context, key string, index, value interface{}) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "lset", key, index, value))
}

// 弹出source尾元素并返回，将弹出元素插入destination列表的开头
func RPoplpush(ctx context.Context, key, source, destination string) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "rpoplpush ", key, source, destination))
}

// 阻塞并弹出尾元素，将弹出元素插入另一列表的开头
func BRpoplpush(ctx context.Context, key, source, destination string, timeout interface{}) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "brpoplpush ", key, source, destination, timeout))
}

// 移除元素,count = 0 : 移除表中所有与 value 相等的值,count!=0,移除与 value 相等的元素，数量为 count的绝对值
func LRem(ctx context.Context, key string, count, value interface{}) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "lrem", key, count, value))
}

// 列表裁剪，让列表只保留指定区间内的元素，不在指定区间之内的元素都将被删除。-1 表示尾部
func LTrim(ctx context.Context, key string, start, stop interface{}) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "ltrim", key, start, stop))
}

func LLen(ctx context.Context, key string) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "llen", key))
}
