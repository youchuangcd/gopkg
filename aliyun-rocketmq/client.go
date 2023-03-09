package aliyunrocketmq

import (
	mqhttpsdk "github.com/aliyunmq/mq-http-go-sdk"
	"github.com/panjf2000/ants/v2"
	"github.com/youchuangcd/gopkg/common/utils"
	"sync"
	"time"
)

type Config struct {
	AccessKeyId     string
	AccessKeySecret string
	EndPoint        string
	// Producer
	ProducerTimeout int64 // 生产者请求超时设置
	MaxConnsPerHost int   // 最大连接数
	// Consumer
	MaxConsumeLimit int32 // 一次最多消费3条（最多可设置为16条）。
	LongLoopTime    int64 // 长轮询时间3s（最多可设置为30s）。
	TcpClient       struct {
		EndPoint string
	}
	GoPoolSize int // 协程池数量
	GoPool     *ants.Pool
}

// New
//
//	@Description:
//	@param conf
//	@return RocketMQ
func New(conf Config) RocketMQ {
	mq := RocketMQ{
		accessKeyId:     conf.AccessKeyId,
		accessKeySecret: conf.AccessKeySecret,
		endPoint:        conf.EndPoint,
		maxConsumeLimit: conf.MaxConsumeLimit,
		longLoopTime:    conf.LongLoopTime,
		//clients: &sync.Pool{
		//	New: func() interface{} {
		//		return mq.NewClient()
		//	},
		//},
		producerTimeout: conf.ProducerTimeout, // 生产者超时时间; 注意：不宜设置过长，因为会导致连接一直等到超时后才释放。造成没有可用的连接错误
		maxConnsPerHost: conf.MaxConnsPerHost, // 最大连接数
		tcpEndPoint:     conf.TcpClient.EndPoint,
	}
	mq.producerClient = mq.newProducerClient()
	mq.consumerClient = mq.newConsumerClient()

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

type RocketMQ struct {
	producerClient  mqhttpsdk.MQClient
	consumerClient  mqhttpsdk.MQClient
	clients         *sync.Pool
	accessKeyId     string
	accessKeySecret string
	endPoint        string
	tcpEndPoint     string
	securityToken   string

	// Consumer
	maxConsumeLimit int32 // 一次最多消费3条（最多可设置为16条）。
	longLoopTime    int64 // 长轮询时间3s（最多可设置为30s）。

	producerTimeout int64 // 生产者超时 单位 秒
	maxConnsPerHost int   // 最大连接数
	goPool          *ants.Pool
}

func (p *RocketMQ) newProducerClient() mqhttpsdk.MQClient {
	timeout := p.producerTimeout
	if timeout == 0 {
		timeout = 5
	}
	return mqhttpsdk.NewAliyunMQClientWithTimeout(p.endPoint, p.accessKeyId, p.accessKeySecret, p.securityToken, time.Second*time.Duration(timeout), p.maxConnsPerHost)
}

func (p *RocketMQ) newConsumerClient() mqhttpsdk.MQClient {
	return mqhttpsdk.NewAliyunMQClientWithTimeout(p.endPoint, p.accessKeyId, p.accessKeySecret, p.securityToken, time.Second*time.Duration(mqhttpsdk.DefaultTimeout), p.maxConnsPerHost)
}
