package redis

import (
	"context"
	"github.com/gomodule/redigo/redis"
	"time"
)

// Lock
// @Description: 加锁
// @param ctx context.Context
// @param key
// @param value
// @param expire
// @param localTimeout
// @return locked
// @return local
// @return err
func Lock(ctx context.Context, key string, value interface{}, expire int64) (locked bool, local bool, err error) {
	return LockLocalTimeout(ctx, key, value, expire, time.Duration(0))
}

// LockLocalTimeout
// @Description: 加锁，如果加锁失败会一直尝试，直到本地超时才返回
// @param ctx context.Context
// @param key
// @param value
// @param expire 秒级
// @param localTimeout 毫秒级
// @param extraArgs 指定睡眠时间间隔
// @return locked
// @return local
// @return err
func LockLocalTimeout(ctx context.Context, key string, value interface{}, expire int64, localTimeout time.Duration, extraArgs ...time.Duration) (locked bool, local bool, err error) {
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()

	// 默认睡眠100ms重新尝试
	var sleepInterval = 100 * time.Millisecond
	if len(extraArgs) == 1 {
		sleepInterval = extraArgs[0]
	}
	// 本地等待结束时间
	var localWaitEndTime time.Time

	// 本地超时标记，外部可以用来区分
	local = false
	if localTimeout > 0 {
		localWaitEndTime = time.Now().Add(localTimeout)
	}
	args := []interface{}{
		key,
		value,
	}
	if expire > 0 {
		args = append(args, "EX", expire)
	}
	args = append(args, "NX")
	for {
		_, err = redis.String(redis.DoContext(c, ctx, "SET", args...))
		locked = true
		if err == redis.ErrNil {
			// 加锁失败，锁已存在
			locked, err = false, nil
		} else if err != nil { // redis操作报错
			locked = false
		}
		// 不需要本地等待超时
		if localTimeout == 0 {
			break
		} else if localTimeout > 0 { // 需要循环尝试加锁直到超过本地超时时间
			// 检查是否超过本地等待锁定的时间, 超过的话就返回true。防止redis不可用时，一直死循环等锁
			if locked || time.Now().After(localWaitEndTime) {
				if !locked {
					locked, local, err = true, true, nil
				}
				break
			}

		} else if locked { // 一直等锁
			break
		}
		// 睡眠100ms重新尝试
		time.Sleep(sleepInterval)
	}
	return locked, local, err
}

// UnLock
// @Description: 解锁
// @param ctx context.Context
// @param key
// @param args [可选]指定锁的内容，比较一致才解锁
// @return err
func UnLock(ctx context.Context, key string, args ...interface{}) (err error) {
	// 如果是指定锁的内容，需要比对一致才可以解锁
	if len(args) == 1 {
		_, err = EvalScript(ctx, ScriptKeyValueEqualsUnlock, key, args[0]).Int()
		return err
	}
	c, _ := getPoolInstance(ctx).GetContext(ctx)
	defer c.Close()
	_, err = redis.Int(redis.DoContext(c, ctx, "DEL", key))
	return
}
