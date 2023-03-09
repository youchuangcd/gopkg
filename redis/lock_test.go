package redis

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	//准备工作
	fmt.Println("start prepare")
	configs := make([]Config, 0, 1)
	configs = append(configs, Config{
		Name:        "default",
		Host:        "127.0.0.1",
		Port:        6379,
		Password:    "",
		MaxIdle:     10,
		MaxActive:   500,
		IdleTimeout: 200,
		Database:    0,
	})
	// 初始化Redis
	InitRedis(configs)

	exitCode := m.Run()

	//清理工作
	fmt.Println("prepare to clean")
	os.Exit(exitCode)
}

func TestLockLocalTimeout(t *testing.T) {
	key := strconv.FormatInt(time.Now().UnixNano(), 10)
	val := time.Now().UnixNano()
	localTimeout := time.Duration(0)
	sleepInterval := 100 * time.Millisecond
	// 应该成功
	startTime := time.Now()
	ctx := context.Background()
	locked, local, err := LockLocalTimeout(ctx, key, val, 10, localTimeout)
	endTime := time.Now()
	if err != nil {
		t.Error("[应该成功] redis操作错误1") // 如果不是如预期的那么就报错
	}
	if !locked {
		t.Error("[应该成功] 加锁失败1") // 如果不是如预期的那么就报错
	} else {
		remainDuration, _ := TTL(ctx, key).Int()
		t.Log("加锁成功1, 是否为本地加锁: ", local, ",本地最多等待", localTimeout, "剩余时长：", remainDuration, "s 等待", endTime.Sub(startTime).Milliseconds(), "ms") // 记录一些你期望记录的信息
	}
	localTimeout = 3 * time.Second
	// 假成功，本地认为加锁成功
	startTime = time.Now()
	locked, local, err = LockLocalTimeout(ctx, key, val, 2, localTimeout, sleepInterval)
	endTime = time.Now()
	if err != nil {
		t.Error("redis操作错误2") // 如果不是如预期的那么就报错
	}
	if !locked {
		t.Error("本地超时加锁失败2") // 如果不是如预期的那么就报错
	} else {
		remainDuration, _ := TTL(ctx, key).Int()
		t.Log("加锁成功2, 是否为本地加锁: ", local, ",本地最多等待", localTimeout, "剩余时长：", remainDuration, "s 等待", endTime.Sub(startTime).Milliseconds(), "ms") // 记录一些你期望记录的信息
	}
	// 不等待直接加锁，应该失败
	locked, local, err = LockLocalTimeout(ctx, key, val, 10, 0)
	if err != nil {
		t.Error("redis操作错误222") // 如果不是如预期的那么就报错
	}
	if locked {
		t.Error("本地加锁应该失败222，但是成功了") // 如果不是如预期的那么就报错
	}

	// 真成功，非本地加锁成功
	localTimeout = 10 * time.Second
	sleepInterval = 1 * time.Second
	startTime = time.Now()
	locked, local, err = LockLocalTimeout(ctx, key, val, 20, localTimeout, sleepInterval)
	endTime = time.Now()
	if err != nil {
		t.Error("redis操作错误3") // 如果不是如预期的那么就报错
	}
	if !locked {
		t.Error("本地超时加锁失败3") // 如果不是如预期的那么就报错
	} else {
		remainDuration, _ := TTL(ctx, key).Int()
		t.Log("加锁成功3, 是否为本地加锁: ", local, ",本地最多等待", localTimeout, "剩余时长：", remainDuration, "s 等待", endTime.Sub(startTime).Milliseconds(), "ms") // 记录一些你期望记录的信息
	}
	// 不等待直接加锁，应该失败
	locked, local, err = LockLocalTimeout(ctx, key, val, 10, 0)
	if err != nil {
		t.Error("redis操作错误333") // 如果不是如预期的那么就报错
	}
	if locked {
		t.Error("本地加锁应该失败333，但是成功了") // 如果不是如预期的那么就报错
	}

	// 真成功，非本地加锁成功
	localTimeout = time.Duration(-1)
	startTime = time.Now()
	locked, local, err = LockLocalTimeout(ctx, key, val, 5, localTimeout, sleepInterval)
	endTime = time.Now()
	if err != nil {
		t.Error("redis操作错误4") // 如果不是如预期的那么就报错
	}
	if !locked {
		t.Error("本地超时加锁失败4") // 如果不是如预期的那么就报错
	} else {
		remainDuration, _ := TTL(ctx, key).Int()
		t.Log("阻塞式等待加锁成功4, 是否为本地加锁: ", local, ",本地最多等待", localTimeout, "剩余时长：", remainDuration, "s 等待", endTime.Sub(startTime).Milliseconds(), "ms") // 记录一些你期望记录的信息
	}
}
