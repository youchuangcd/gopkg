package redis

import (
	"context"
	"github.com/gomodule/redigo/redis"
)

type Z struct {
	Score  float64
	Member interface{}
}

// 添加元素
func ZAdd(ctx context.Context, key string, members ...Z) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	args := redis.Args{}.Add(key)
	for _, v := range members {
		args = args.Add(v.Score).Add(v.Member)
	}
	return getReply(redis.DoContext(c, ctx, "zadd", args...))
}

// 增加元素权重
func ZIncrBy(ctx context.Context, key string, increment, member interface{}) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "zincrby", key, increment, member))
}

// 增加元素权重
func ZCard(ctx context.Context, key string) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "zcard", key))
}

// 返回指定元素的排名
func ZRank(ctx context.Context, key string, member interface{}) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "zrank", key, member))
}

// 返回指定元素的权重
func ZScore(ctx context.Context, key string, member interface{}) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "zscore", key, member))
}

// 返回集合两个权重间的元素数
func ZCount(ctx context.Context, key string, min, max interface{}) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "zcount", key, min, max))
}

// 返回指定区间内的元素
func ZRange(ctx context.Context, key string, start, stop interface{}, withScore ...bool) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	if len(withScore) > 0 && withScore[0] {
		return getReply(redis.DoContext(c, ctx, "zrange", key, start, stop, "WITHSCORES"))
	}
	return getReply(redis.DoContext(c, ctx, "zrange", key, start, stop))
}

// 倒序返回指定区间内的元素
func ZRevRange(ctx context.Context, key string, start, stop interface{}, withScore ...bool) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	if len(withScore) > 0 && withScore[0] {
		return getReply(redis.DoContext(c, ctx, "zrevrange", key, start, stop, "WITHSCORES"))
	}
	return getReply(redis.DoContext(c, ctx, "zrevrange", key, start, stop))
}
