package model

import (
	"context"
	"fmt"
	"github.com/youchuangcd/gopkg"
	"github.com/youchuangcd/gopkg/mylog"
	"github.com/youchuangcd/gopkg/redis"
)

func concurrentLockJoinPrefix(ctx context.Context, key string) string {
	res := fmt.Sprintf(gopkg.CacheConcurrentLockPrefix, key)
	return res
}

// CacheConcurrentLock
// @Description: 并发请求锁
// @param ctx
// @param key
// @param expire
// @param args 指定锁的值
// @return locked
// @return err
func CacheConcurrentLock(ctx context.Context, key string, expire int64, args ...interface{}) (locked bool, err error) {
	key = concurrentLockJoinPrefix(ctx, key)
	var value interface{} = 1
	if len(args) == 1 {
		value = args[0]
	}
	locked, _, err = redis.Lock(ctx, key, value, expire)
	if err != nil {
		logContent := map[string]interface{}{
			"err":  err,
			"key":  key,
			"args": args,
		}
		mylog.WithError(ctx, gopkg.LogRedis, logContent, "并发请求加锁失败")
	}
	return
}

// CacheConcurrentUnLock
// @Description: 解除并发请求锁
// @param ctx
// @param key
// @param args 指定锁的值
// @return err
func CacheConcurrentUnLock(ctx context.Context, key string, args ...interface{}) (err error) {
	key = concurrentLockJoinPrefix(ctx, key)
	err = redis.UnLock(ctx, key, args...)
	if err != nil {
		logContent := map[string]interface{}{
			"err":  err,
			"key":  key,
			"args": args,
		}
		mylog.WithError(ctx, gopkg.LogRedis, logContent, "并发请求解锁失败")
	}
	return
}
