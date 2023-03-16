package gopkg

import "gorm.io/gorm"

/**
 * 公共包依赖外部的所有配置都在此处，按需修改
 */
var (
	Success                    = NewError(0, "")
	Failure                    = NewError(500, "操作失败，请稍后再试!")
	ErrorRedisInvalidScriptKey = NewError(800, "无效的redis Lua脚本key")
	ErrorInternalServer        = NewError(999, "系统繁忙，请稍后再试!")

	InvalidParam      = NewError(10005, "无效的参数")
	ErrorInvalidToken = NewError(10001, "token无效或者授权已过期，请重新申请授权")
	//Success                    *Error
	//Failure                    *Error
	//ErrorRedisInvalidScriptKey *Error
	//InvalidParam               *Error
	//ErrorInvalidToken          *Error
)

var (
	// Env 设置项目环境
	Env string
	// 环境变量
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvTest  = "test"
	EnvGray  = "gray"
	EnvProd  = "prod"
)

// 日志分类名称
var (
	LogCtx         = "" // 根据上下文来区分到底是哪个模块的日志
	LogDB          = "db"
	LogRedis       = "redis"
	LogKafka       = "kafka"
	LogHttp        = "http"
	LogRocketMQ    = "rocketmq"
	LogRocketMQTCP = "rocketmq-tcp"
	LogNacos       = "nacos"
	LogPanic       = "panic"
	LogSys         = "system"
)

var (
	RequestLogTypeFlag      = "request"
	RequestHeaderTraceIdKey = "X-Request-Id"
	RequestHeaderSpanIdKey  = "X-Request-Span-Id"
	LogTraceIdKey           = "traceId"
	LogSpanIdKey            = "spanId"
	LogUserIdKey            = "userId"
	LogTaskIdKey            = "taskId"
	LogMsgIdKey             = "msgId"
	LogParentMsgIdKey       = "parentMsgId"

	// 请求来源
	ContextRequestSourceKey = "X-Request-Source"
	ContextLogCategoryKey   = "logCategory"

	ContextRequestUserIdKey    = "userId"
	ContextRequestSysUserIdKey = "sysUserId"
	ContextRequestParamKey     = "requestParam"
	ContextResponseDataKey     = "responseData"
	ContextRequestStartTimeKey = "requestStartTime"
	// ContextResponseBodyWriterKey 替换gin c.Writer
	ContextResponseBodyWriterKey = "responseBodyWriter"
)

var (
	DateFormat                      = "2006-01-02"
	DateTimeFormat                  = "2006-01-02 15:04:05"
	DateTimeNoSecondFormat          = "2006-01-02 15:04"
	DateTimeSlashFormat             = "2006/01/02 15:04"
	DateMsTimeFormat                = "2006-01-02 15:04:05.999"
	DateNoDelimitersFormat          = "20060102"
	DateTimeNoDelimitersFormat      = "20060102150405"
	TimeFormat                      = "15:04:05"
	MysqlTimeTypeDefaultValue       = "00:00:00"                  // mysql time类型默认值
	DateTimeUtcFormat               = "2006-01-02T15:04:05+08:00" //utc时间格式
	DateNoDelimitersYearMonthFormat = "200601"
)

// DB相关变量
var (
	GormDBMap           map[string]*gorm.DB
	GormDBMapKeyDefault = "default"
)

// redis相关变量
var (
	// 上下文里记录要切换的redis实例的key
	ContextRedisMapKey = "redisMapKey"
	// Redis 默认名称: 默认使用哪个redis实例
	RedisMapKeyDefault        = "default"
	CacheConcurrentLockPrefix = "conc_lock:%s" // 并发锁请求前缀 s1= 自定义key
)

var (
	// LogLimitContentLength 日志内容长度限制, 超过限制就截取用 关键词替换
	LogLimitContentLength uint = 1000
	// LogLimitContentReplaceWord 日志内容长度超出限制就替换为指定关键词
	LogLimitContentReplaceWord = "...内容截取..."
)

var (
	RobotDataRelationUrl = "" // 机器人数据关联地址
)
