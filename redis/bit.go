package redis

import (
	"context"
	"github.com/gomodule/redigo/redis"
)

func SetBit(ctx context.Context, key string, offset, value int64) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "setbit", key, offset, value))
}

func GetBit(ctx context.Context, key string, offset int64) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "getbit", key, offset))
}

func BitCount(ctx context.Context, key string, interval ...int64) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	if len(interval) == 2 {
		return getReply(c.Do("bitcount", key, interval[0], interval[1]))
	}
	return getReply(redis.DoContext(c, ctx, "bitcount", key))
}

// opt 包含 and、or、xor、not
func BitTop(ctx context.Context, opt, destKey string, keys ...string) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "bitop", opt, redis.Args{}.Add(destKey).AddFlat(keys)))
}
