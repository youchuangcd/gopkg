package apacherocketmq

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

var (
	logConf LogConfig
)

// 公共的日志配置
type LogConfig struct {
	Category   string
	Logger     logger
	Limit      uint
	ReplaceStr string
	Producer   bool // 推送日志开启
}

// Init
//
//	@Description:
//	@param conf
func Init(conf LogConfig) {
	logConf = conf
}

type logger interface {
	LogDebug(ctx context.Context, logCategory string, logContent map[string]interface{}, msg string)
	LogInfo(ctx context.Context, logCategory string, logContent map[string]interface{}, msg string)
	LogWarn(ctx context.Context, logCategory string, logContent map[string]interface{}, msg string)
	LogError(ctx context.Context, logCategory string, logContent map[string]interface{}, msg string)
}

type Config struct {
	EndPoint   string
	RetryTimes int // 推送消息重试次数
	Namespace  string
}

type RocketMQ struct {
	rocketmq.Producer
	endPoint   string
	retryTimes int // 推送消息重试次数
	namespace  string
}

// New
//
//	@Description:
//	@param conf
//	@return *RocketMQ
func New(conf Config) *RocketMQ {
	mq := &RocketMQ{
		endPoint:   conf.EndPoint,
		retryTimes: conf.RetryTimes,
		namespace:  conf.Namespace,
	}
	var err error
	mq.Producer, err = rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{mq.endPoint})),
		producer.WithRetry(mq.retryTimes),
		producer.WithNamespace(mq.namespace),
	)
	if err != nil {
		panic("rocketmq producer new failed：" + err.Error())
	}
	err = mq.Producer.Start()
	if err != nil {
		panic("rocketmq producer start failed：" + err.Error())
	}
	return mq
}

func (p *RocketMQ) SendSync(ctx context.Context, msgs ...*primitive.Message) (res *primitive.SendResult, err error) {
	res, err = p.Producer.SendSync(ctx, msgs...)
	if err != nil {
		logConf.Logger.LogError(ctx, logConf.Category, map[string]interface{}{
			"err":  err,
			"msgs": msgs,
		}, "send msg failed")
	} else if logConf.Producer {
		logConf.Logger.LogInfo(ctx, logConf.Category, map[string]interface{}{
			"producerMsgId": res.MsgID,
			"msgs":          msgs,
			"res":           res,
		}, "send msg success")
	}
	return
}

func (p *RocketMQ) SendAsync(ctx context.Context, mq func(ctx context.Context, result *primitive.SendResult, err error), msgs ...*primitive.Message) (err error) {
	err = p.Producer.SendAsync(ctx, func(ctx context.Context, res *primitive.SendResult, err error) {
		mq(ctx, res, err)
		if err == nil && logConf.Producer {
			logConf.Logger.LogInfo(ctx, logConf.Category, map[string]interface{}{
				"producerMsgId": res.MsgID,
				"msgs":          msgs,
				"res":           res,
			}, "send msg success")
		}
	}, msgs...)
	if err != nil {
		logConf.Logger.LogError(ctx, logConf.Category, map[string]interface{}{
			"err":  err,
			"msgs": msgs,
		}, "send async msg failed")
	}
	return
}

//
//func (p *RocketMQ) SendOneWay(ctx context.Context, mq ...*primitive.Message) error {
//	return p.Producer.SendOneWay(ctx, mq...)
//}
//func (p *RocketMQ) Request(ctx context.Context, ttl time.Duration, msg *primitive.Message) (*primitive.Message, error) {
//	return p.Producer.Request(ctx, ttl, msg)
//}
//func (p *RocketMQ) RequestAsync(ctx context.Context, ttl time.Duration, callback func(ctx context.Context, msg *primitive.Message, err error), msg *primitive.Message) error {
//	return p.Producer.RequestAsync(ctx, ttl, callback, msg)
//}
