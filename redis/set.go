package redis

import (
	"context"
	"github.com/gomodule/redigo/redis"
)

// 添加元素
func SAdd(ctx context.Context, key string, member interface{}, members ...interface{}) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	args := redis.Args{}.Add(key).AddFlat(member)
	if len(members) > 0 {
		args = args.AddFlat(members)
	}
	return getReply(redis.DoContext(c, ctx, "sadd", args...))
}

// 集合元素个数
func SCard(ctx context.Context, key string) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "scard", key))
}

// 返回集合中成员
func SMembers(ctx context.Context, key string) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "smembers", key))
}

// 判断元素是否是集合成员
func SisMember(ctx context.Context, key string, member interface{}) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "sismember", key, member))
}

// 随机返回并移除一个元素
func SPop(ctx context.Context, key string) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "spop", key))
}

// 随机返回一个或多个元素
func SRandMember(ctx context.Context, key string, count ...int64) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	if len(count) > 0 {
		return getReply(redis.DoContext(c, ctx, "srandmember", key, count[0]))
	}
	return getReply(redis.DoContext(c, ctx, "srandmember", key))
}

// 移除指定的元素
func SRem(ctx context.Context, key string, member interface{}, members ...interface{}) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	args := redis.Args{}.Add(key).AddFlat(member)
	if len(members) > 0 {
		args = args.AddFlat(members)
	}
	return getReply(redis.DoContext(c, ctx, "srem", args...))
}

// 将元素从集合移至另一个集合
func SMove(ctx context.Context, sourceKey, destinationKey string, member interface{}) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "smove", sourceKey, destinationKey, member))
}

// 返回一或多个集合的差集
func SDiff(ctx context.Context, keys []string) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "sdiff", redis.Args{}.AddFlat(keys)...))
}

// 将一或多个集合的差集保存至另一集合(destinationKey)
func SDiffStore(ctx context.Context, destinationKey string, keys []string) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "sdiffstore", redis.Args{}.Add(destinationKey).AddFlat(keys)...))
}

// 将keys的集合的并集 写入到 destinationKey中
func SInterStore(ctx context.Context, destinationKey string, keys []string) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "sinterstore", redis.Args{}.Add(destinationKey).AddFlat(keys)...))
}

// 一个或多个集合的交集
func SInter(ctx context.Context, keys []string) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "sinter", redis.Args{}.AddFlat(keys)...))
}

// 返回集合的并集
func SUnion(ctx context.Context, keys []string) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "sunion", redis.Args{}.AddFlat(keys)...))
}

// 将 keys 的集合的并集 写入到 destinationKey 中
func SUnionStore(ctx context.Context, destinationKey string, keys []string) *Reply {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	return getReply(redis.DoContext(c, ctx, "sunionstore", redis.Args{}.Add(destinationKey).AddFlat(keys)...))
}
