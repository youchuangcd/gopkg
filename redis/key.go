package redis

import (
	"context"
	"github.com/gomodule/redigo/redis"
)

// 查找键 [*模糊查找]
func Keys(ctx context.Context, key string) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "keys", key))
}

// 判断key是否存在
func Exists(ctx context.Context, key string) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "exists", key))
}

// 随机返回一个key
func RandomKey(ctx context.Context) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "randomkey"))
}

// 返回值类型
func Type(ctx context.Context, key string) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "type", key))
}

// 删除key
func Del(ctx context.Context, keys ...string) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "del", redis.Args{}.AddFlat(keys)...))
}

// 重命名
func Rename(ctx context.Context, key, newKey string) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "rename", key, newKey))
}

// 仅当newkey不存在时重命名
func RenameNX(ctx context.Context, key, newKey string) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "renamenx", key, newKey))
}

// 序列化key
func Dump(ctx context.Context, key string) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "dump", key))
}

// 反序列化
func Restore(ctx context.Context, key string, ttl, serializedValue interface{}) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "restore", key, ttl, serializedValue))
}

// 秒
func Expire(ctx context.Context, key string, seconds int64) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "expire", key, seconds))
}

// 秒
func ExpireAt(ctx context.Context, key string, timestamp int64) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "expireat", key, timestamp))
}

// 毫秒
func Persist(ctx context.Context, key string) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "persist", key))
}

// 毫秒
func PersistAt(ctx context.Context, key string, milliSeconds int64) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "persistat", key, milliSeconds))
}

// 秒
func TTL(ctx context.Context, key string) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "ttl", key))
}

// 毫秒
func PTTL(ctx context.Context, key string) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "pttl", key))
}

// 同实例不同库间的键移动
func Move(ctx context.Context, key string, db int64) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "move", key, db))
}
