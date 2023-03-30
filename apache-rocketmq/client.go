package apacherocketmq

import (
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/panjf2000/ants/v2"
	"github.com/youchuangcd/gopkg/common/utils"
)

type Config struct {
	EndPoint   string
	RetryTimes int // 推送消息重试次数
	GoPoolSize int // 协程池数量
	GoPool     *ants.Pool
}

type RocketMQ struct {
	producerClient rocketmq.Producer
	consumerClient rocketmq.PushConsumer
	endPoint       string
	goPool         *ants.Pool
}

// New
//
//	@Description:
//	@param conf
//	@return *RocketMQ
func New(conf Config) *RocketMQ {
	mq := &RocketMQ{
		endPoint: conf.EndPoint,
	}
	var err error
	mq.producerClient, err = rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{mq.endPoint})),
		producer.WithRetry(conf.RetryTimes),
	)
	if err != nil {

	}
	mq.consumerClient, err = rocketmq.NewPushConsumer(
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{mq.endPoint})),
		consumer.WithRetry(conf.RetryTimes),
	)

	mq.goPool = conf.GoPool
	if mq.goPool == nil {
		var goPoolSize = 1000
		if conf.GoPoolSize > 0 {
			goPoolSize = conf.GoPoolSize
		}
		goPool, err := utils.NewPool(goPoolSize)
		if err != nil {
			panic("初始化协程池失败")
		}
		mq.goPool = goPool
	}
	return mq
}
