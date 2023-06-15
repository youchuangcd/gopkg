package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/panjf2000/ants/v2"
	"github.com/youchuangcd/gopkg/common"
	"strings"
)

var (
	// 使用kafka包，需要设置日志配置
	logConf LogConfig
	// 上下文中的traceId key
	ctxTraceIdKey string
	// 上下文中的msgId key
	ctxMsgIdKey string
	ctxWithMap  map[string]struct{}
	// 生成唯一id的方法
	genUniqIdFunc func() string
	kafkaClientId string
)

// Init
//
//	@Description:
//	@param clientId kafka客户端名称
//	@param conf
//	@param traceIdKey 上下文中traceId 的key
//	@param msgIdKey
//	@param genUniqIdHandler
func Init(clientId string, conf LogConfig, traceIdKey, msgIdKey string, genUniqIdHandler func() string) {
	kafkaClientId = clientId
	logConf = conf
	ctxTraceIdKey = traceIdKey
	ctxMsgIdKey = msgIdKey
	ctxWithMap = make(map[string]struct{}, 2)
	ctxWithMap[ctxTraceIdKey] = struct{}{}
	ctxWithMap[ctxMsgIdKey] = struct{}{}
	genUniqIdFunc = genUniqIdHandler
}

// SetLogConfig
//
//	@Description: 设置日志配置
//	@param conf
func SetLogConfig(conf LogConfig) {
	logConf = conf
}

// SetRequestHeaderTraceIdKey
//
//	@Description: 设置上下文中的traceId key
//	@param k
func SetRequestHeaderTraceIdKey(k string) {
	ctxTraceIdKey = k
	ctxWithMap[ctxTraceIdKey] = struct{}{}
}

// SetRequestHeaderMsgIdKey
//
//	@Description: 设置上下文中的msgId key
//	@param k
func SetRequestHeaderMsgIdKey(k string) {
	ctxMsgIdKey = k
	ctxWithMap[ctxMsgIdKey] = struct{}{}
}

// SetLogProducer
//
//	@Description: 设置是否记录生产者日志
//	@param v
func SetLogProducer(v bool) {
	logConf.Producer = v
}

// SetLogConsumer
//
//	@Description: 设置是否记录消费者日志
//	@param v
func SetLogConsumer(v bool) {
	logConf.Consumer = v
}

// 公共的日志配置
type LogConfig struct {
	Category              string
	Logger                logger
	Limit                 uint
	ReplaceStr            string
	PanicHandler          func(r interface{})
	GoroutinePanicHandler func(r interface{})
	Producer              bool // 推送日志开启
	Consumer              bool // 消费日志开启
}

type logger interface {
	LogDebug(ctx context.Context, logCategory string, logContent map[string]interface{}, msg string)
	LogInfo(ctx context.Context, logCategory string, logContent map[string]interface{}, msg string)
	LogWarn(ctx context.Context, logCategory string, logContent map[string]interface{}, msg string)
	LogError(ctx context.Context, logCategory string, logContent map[string]interface{}, msg string)
}
type Kafka struct {
	callback             func(ctx context.Context, s *sarama.ConsumerMessage) error
	callbackBatchProcess func(ctx context.Context, msgs []*sarama.ConsumerMessage) error
	aggregator           *common.Aggregator
	group                string
	consumerAddrs        []string
	producerAddrs        []string
	syncProducer         sarama.SyncProducer
	goPool               *ants.Pool // 协程池
	consumerOffsets      int64      // 消费者偏移量类型设置 OffsetNewest or OffsetOldest
}

type Config struct {
	Group           string
	ConsumerHost    string
	ProducerHost    string
	consumerOffsets int64 // 消费者偏移量类型设置 OffsetNewest or OffsetOldest
}

type ConsumerConfig struct {
	ConsumerOffsets int64 // 消费者偏移量类型设置 OffsetNewest or OffsetOldest
}

func New(conf Config) Kafka {
	k := Kafka{
		group:           conf.Group,
		producerAddrs:   strings.Split(conf.ProducerHost, ","),
		consumerAddrs:   strings.Split(conf.ConsumerHost, ","),
		consumerOffsets: conf.consumerOffsets,
	}
	kConf := k.getConfig()
	addrs := k.getProducerAddr()
	syncProducer, err := sarama.NewSyncProducer(addrs, kConf)
	if err != nil {
		panic("NewSyncProducer failed: " + err.Error())
	}
	// Wrap instrumentation
	//syncProducer = otelsarama.WrapSyncProducer(kConf, syncProducer)
	k.syncProducer = syncProducer
	return k
}

func (k Kafka) getProducerAddr() []string {
	return k.producerAddrs
}
func (k Kafka) getConsumerAddr() []string {
	return k.consumerAddrs
}

func (k Kafka) getConfig() *sarama.Config {
	conf := sarama.NewConfig()
	conf.ClientID = kafkaClientId
	conf.Producer.Return.Successes = true
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Partitioner = sarama.NewRandomPartitioner
	conf.Consumer.Return.Errors = true
	conf.Consumer.Offsets.Initial = sarama.OffsetOldest
	if k.consumerOffsets != 0 {
		conf.Consumer.Offsets.Initial = k.consumerOffsets
	}
	//conf.Consumer.Offsets.AutoCommit.Enable = false //手动提交偏移量
	conf.Version = sarama.V0_11_0_1 //kafka server的版本号
	sarama.PanicHandler = logConf.PanicHandler
	return conf
}

func (k Kafka) cutStrFromLogConfig(s string) string {
	return cutStr(s, logConf.Limit, logConf.ReplaceStr)
}

func cutStr(s string, limit uint, rs string) string {
	runeStr := []rune(s)
	sl := len(runeStr)
	if sl > int(limit) {
		halfLen := limit / 2
		var buff strings.Builder
		buff.WriteString(string(runeStr[0:halfLen]))
		buff.WriteString(rs)
		buff.WriteString(string(runeStr[sl-int(halfLen):]))
		return buff.String()
	}
	return s
}

func (k Kafka) Close() {
	k.syncProducer.Close()
}
