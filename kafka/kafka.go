package kafka

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/panjf2000/ants/v2"
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

// 公共的日志配置
type LogConfig struct {
	Category              string
	Logger                logger
	Limit                 uint
	ReplaceStr            string
	PanicHandler          func(r interface{})
	GoroutinePanicHandler func(r interface{})
}

type logger interface {
	LogDebug(ctx context.Context, logCategory string, logContent map[string]interface{}, msg string)
	LogInfo(ctx context.Context, logCategory string, logContent map[string]interface{}, msg string)
	LogWarn(ctx context.Context, logCategory string, logContent map[string]interface{}, msg string)
	LogError(ctx context.Context, logCategory string, logContent map[string]interface{}, msg string)
}
type Kafka struct {
	callback      func(ctx context.Context, s *sarama.ConsumerMessage) error
	group         string
	consumerAddrs []string
	producerAddrs []string
	logReplaceStr string
	syncProducer  sarama.SyncProducer
	goPool        *ants.Pool // 协程池
}

type Config struct {
	Group        string
	ConsumerHost string
	ProducerHost string
}

// handler，核心的消费者业务实现
type consumerGroupHandler struct {
	Kafka
}

func New(conf Config) Kafka {
	k := Kafka{
		group:         conf.Group,
		producerAddrs: strings.Split(conf.ProducerHost, ","),
		consumerAddrs: strings.Split(conf.ConsumerHost, ","),
	}
	kConf := k.getConfig()
	addrs := k.getProducerAddr()
	syncProducer, err := sarama.NewSyncProducer(addrs, kConf)
	if err != nil {
		panic("NewSyncProducer failed: " + err.Error())
	}
	k.syncProducer = syncProducer
	return k
}

func (k Kafka) getProducerAddr() []string {
	return k.producerAddrs
}
func (k Kafka) getConsumerAddr() []string {
	return k.consumerAddrs
}

func (h consumerGroupHandler) Setup(s sarama.ConsumerGroupSession) error {
	// 当连接完毕的时候会通知这个，start
	return nil
}
func (h consumerGroupHandler) Cleanup(s sarama.ConsumerGroupSession) error {
	s.Commit()
	// end，当这一次消费完毕，会通知，这里最好commit
	return nil
}
func (h consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error { // consume
	// 当第一个ConsumeClaim消费完成，会话就会被关闭
	//ctx := sess.Context()
	ctx := context.WithValue(context.Background(), "logCategory", logConf.Category)
	for msg := range claim.Messages() {
		select {
		case <-sess.Context().Done():
			break
		default:
		}
		newCtx := ctx
		// 从消息头部中取traceId 和msgId 写到上下文中
		for _, v := range msg.Headers {
			headerKey := string(v.Key)
			if _, ok := ctxWithMap[headerKey]; ok {
				newCtx = context.WithValue(newCtx, headerKey, string(v.Value))
			}
		}
		highWaterMarkOffset := claim.HighWaterMarkOffset()
		tmpMsg := msg
		_err := h.goPool.Submit(func() {
			// 业务逻辑处理
			err := h.callback(newCtx, tmpMsg)
			valStr := h.cutStrFromLogConfig(string(tmpMsg.Value))
			logMap := map[string]interface{}{
				"topic":        tmpMsg.Topic,
				"partition":    tmpMsg.Partition,
				"offset":       tmpMsg.Offset,
				"maxOffsetSub": highWaterMarkOffset - tmpMsg.Offset,
				"key":          string(tmpMsg.Key),
				"value":        valStr,
			}
			if err != nil {
				logMap["err"] = err
				logConf.Logger.LogError(newCtx, logConf.Category, logMap, "[Consumer] Message Failed")
				// 扔到重试队列或死信队列
			} else {
				logConf.Logger.LogInfo(newCtx, logConf.Category, logMap, "[Consumer] Message Success")
			}
		})
		if _err != nil {
			logConf.Logger.LogError(newCtx, logConf.Category, map[string]interface{}{
				"err": _err,
			}, "kafka消费者提交消息到协程池失败")
		}
		sess.MarkMessage(msg, "") // 必须设置这个，不然你的偏移量无法提交。
	}
	return nil
}

func (k Kafka) getConfig() *sarama.Config {
	conf := sarama.NewConfig()
	conf.ClientID = kafkaClientId
	conf.Producer.Return.Successes = true
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Partitioner = sarama.NewRandomPartitioner
	conf.Consumer.Return.Errors = true
	conf.Consumer.Offsets.Initial = sarama.OffsetOldest
	//conf.Consumer.Offsets.AutoCommit.Enable = false //手动提交偏移量
	conf.Version = sarama.V0_11_0_1 //kafka server的版本号
	sarama.PanicHandler = logConf.PanicHandler
	return conf
}

func (k Kafka) Consumer(ctx context.Context, topics []string, consumerGroupName string, cb func(ctx context.Context, s *sarama.ConsumerMessage) error, goPoolSize int) {
	conf := k.getConfig()
	// 没有额外设置地址，取配置地址
	addrs := k.getConsumerAddr()
	k.callback = cb
	var err error
	// 创建协程池
	k.goPool, err = ants.NewPool(goPoolSize, ants.WithPanicHandler(logConf.GoroutinePanicHandler))
	if err != nil {
		panic("消费者" + topics[0] + "初始化协程池失败")
	}
	if consumerGroupName == "" {
		consumerGroupName = k.group
	}
	client, err := sarama.NewConsumerGroup(addrs, consumerGroupName, conf)
	if err != nil {
		logConf.Logger.LogError(ctx, logConf.Category, map[string]interface{}{
			"err":    err,
			"topics": topics,
		}, "Consumer failed")
		panic(fmt.Sprintf("创建消费者分组失败, topics: %v, err: %s", topics, err.Error()))
	}
	defer client.Close()
	handler := consumerGroupHandler{Kafka: k} // 必须传递一个handler
	for {                                     // for循环的目的是因为存在重平衡，他会重新启动
		err = client.Consume(ctx, topics, handler) // consume 操作，死循环。exampleConsumerGroupHandler的ConsumeClaim不允许退出，也就是操作到完毕。
		if err != nil {
			logConf.Logger.LogError(ctx, logConf.Category, map[string]interface{}{
				"topics": topics,
				"err":    err,
			}, "msg consumer failed")
		}
	}
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

// Producer
//
//	@Description: 推送单条消息到kafka
//	@receiver k
//	@param ctx
//	@param topic
//	@param content
//	@return partition
//	@return offset
//	@return err
func (k Kafka) Producer(ctx context.Context, topic string, content string) (partition int32, offset int64, err error) {
	headers := make([]sarama.RecordHeader, 0, 2)
	if traceId, ok := ctx.Value(ctxTraceIdKey).(string); ok {
		headers = append(headers, sarama.RecordHeader{
			Key:   []byte(ctxTraceIdKey),
			Value: []byte(traceId),
		})
	}
	// 生成每条消息的id
	headers = append(headers, sarama.RecordHeader{
		Key:   []byte(ctxMsgIdKey),
		Value: []byte(genUniqIdFunc()),
	})
	msg := &sarama.ProducerMessage{
		Topic:   topic,
		Value:   sarama.StringEncoder(content),
		Headers: headers,
	}

	// 发送消息
	partition, offset, err = k.syncProducer.SendMessage(msg)
	//content = k.cutStrFromLogConfig(content)
	if err != nil {
		content = k.cutStrFromLogConfig(content)
		logConf.Logger.LogError(ctx, logConf.Category, map[string]interface{}{
			"err":     err,
			"topic":   topic,
			"content": content,
		}, "send msg failed")
		return
	}
	//logConf.Logger.LogInfo(ctx, logConf.Category, map[string]interface{}{
	//	"partition": partition,
	//	"offset":  offset,
	//	"content": content,
	//}, "send msg success")
	return
}

// SyncProducerBatch
//
//	@Description: 推送多条消息到kafka
//	@receiver k
//	@param ctx
//	@param topic
//	@param contents
//	@return err
func (k Kafka) SyncProducerBatch(ctx context.Context, topic string, contents []string) (err error) {
	msgs := make([]*sarama.ProducerMessage, 0, len(contents))
	headers := make([]sarama.RecordHeader, 0, 2)
	if traceId, ok := ctx.Value(ctxTraceIdKey).(string); ok {
		headers = append(headers, sarama.RecordHeader{
			Key:   []byte(ctxTraceIdKey),
			Value: []byte(traceId),
		})
	}
	for _, content := range contents {
		msgHeaders := headers
		// 生成每条消息的id
		msgHeaders = append(msgHeaders, sarama.RecordHeader{
			Key:   []byte(ctxMsgIdKey),
			Value: []byte(genUniqIdFunc()),
		})
		msgs = append(msgs, &sarama.ProducerMessage{
			Topic:   topic,
			Value:   sarama.StringEncoder(content),
			Headers: msgHeaders,
		})
	}

	// 发送消息
	err = k.syncProducer.SendMessages(msgs)
	if err != nil {
		logConf.Logger.LogError(ctx, logConf.Category, map[string]interface{}{
			"err":   err,
			"topic": topic,
		}, "send msg failed")
		return
	}
	//logConf.Logger.LogInfo(ctx, logConf.Category, map[string]interface{}{}, "send msg success")
	return
}

func (k Kafka) Close() {
	k.syncProducer.Close()
}
