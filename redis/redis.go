package redis

import (
	"context"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/youchuangcd/gopkg"
	"github.com/youchuangcd/gopkg/mylog"
	"sync"
	"time"
)

// var clientPool *redis.Pool
var (
	// redisCollections redis对象集合
	redisCollections map[string]*Pool
	// once 确保全局Redis对象只实例一次
	once sync.Once
)

type Pool struct {
	*redis.Pool
	config Config
}

type Config struct {
	Name           string `mapstructure:"name"` // 实例名称
	Host           string `mapstructure:"host"`
	Port           int    `mapstructure:"port"`
	Password       string `mapstructure:"password"`
	MaxIdle        int    `mapstructure:"maxIdle"`        //最大空闲连接数
	MaxActive      int    `mapstructure:"maxActive"`      //最大连接数
	IdleTimeout    int    `mapstructure:"idleTimeout"`    //空闲连接超时时间 单位：纳秒
	Database       int    `mapstructure:"database"`       // 选择db
	ConnectTimeout int    `mapstructure:"connectTimeout"` //连接超时 单位毫秒
	ReadTimeout    int    `mapstructure:"readTimeout"`    //读取超时 单位毫秒
	WriteTimeout   int    `mapstructure:"writeTimeout"`   //写入超时 单位毫秒
}

func InitRedis(configs []Config) {
	once.Do(func() {
		if redisCollections == nil {
			redisCollections = make(map[string]*Pool, len(configs))
		}
		for _, v := range configs {
			if v.ConnectTimeout == 0 {
				v.ConnectTimeout = 1000
			}
			if v.ReadTimeout == 0 {
				v.ReadTimeout = 1000
			}
			if v.WriteTimeout == 0 {
				v.WriteTimeout = 1000
			}
			nv := v
			// 建立连接池
			redisCollections[v.Name] = &Pool{
				config: nv,
				Pool: &redis.Pool{
					MaxIdle:     v.MaxIdle,                                       //最大空闲连接数
					MaxActive:   v.MaxActive,                                     //最大连接数
					IdleTimeout: time.Duration(v.IdleTimeout) * time.Millisecond, //空闲连接超时时间
					Wait:        true,
					DialContext: func(ctx context.Context) (redis.Conn, error) {
						con, err := redis.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", nv.Host, nv.Port),
							redis.DialPassword(nv.Password),
							redis.DialDatabase(nv.Database),
							redis.DialConnectTimeout(time.Duration(nv.ConnectTimeout)*time.Millisecond),
							redis.DialReadTimeout(time.Duration(nv.ReadTimeout)*time.Millisecond),
							redis.DialWriteTimeout(time.Duration(nv.WriteTimeout)*time.Millisecond))
						if err != nil {
							mylog.Error(ctx, gopkg.LogRedis, "[redis init] "+err.Error())
							return nil, err
						}
						return con, nil
					},
				},
			}
		}
	})
}

// getPoolInstance
//
//	@Description: 获取一个连接池对象
//	@param ctx
//	@return *Pool
func getPoolInstance(ctx context.Context) *Pool {
	mapKey := gopkg.RedisMapKeyDefault
	if key, ok := ctx.Value(gopkg.ContextRedisMapKey).(string); ok { // 从上下文里切换
		mapKey = key
	}
	if p, ok := redisCollections[mapKey]; ok {
		return p
	}
	panic(fmt.Sprintf("无效的redis实例key: %s", mapKey))
}

// GetPoolInstance
//
//	@Description: 获取一个连接池对象
//	@param ctx
//	@return *Pool
func GetPoolInstance(ctx context.Context) *Pool {
	return getPoolInstance(ctx)
}

// GetConfig
//
//	@Description: 获取连接池配置
//	@receiver p
//	@return Config
func (p *Pool) GetConfig() Config {
	return p.config
}

// SwitchRedisByCtx
//
//	@Description: 根据上下文切换redis实例
//	@param ctx
//	@param name
//	@return context.Context
func SwitchRedisByCtx(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, gopkg.ContextRedisMapKey, name)
}

// GetConn
// @Description: 从连接池中获取一个连接，要记得close; 不可以用于subscribed to pubsub channel, transaction started, ...
// @return redis.Conn
func GetConn(ctx context.Context) (redis.Conn, error) {
	return getPoolInstance(ctx).GetContext(ctx)
}

func ScanStruct(v []interface{}, obj interface{}) error {
	return redis.ScanStruct(v, obj)
}
func ScanSlice(v []interface{}, obj interface{}, fieldNames ...string) error {
	return redis.ScanSlice(v, obj, fieldNames...)
}
