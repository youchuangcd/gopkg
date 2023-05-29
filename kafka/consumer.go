package kafka

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/panjf2000/ants/v2"
)

func (k Kafka) Consumer(ctx context.Context, topics []string, consumerGroupName string, cb func(ctx context.Context, s *sarama.ConsumerMessage) error, goPoolSize int, args ...interface{}) {
	if len(args) > 0 {
		if conf, ok := args[0].(ConsumerConfig); ok {
			k.consumerOffsets = conf.ConsumerOffsets
		}
	}
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
	if consumerGroupName == "" && k.group != "" {
		consumerGroupName = k.group
	}
	k.group = consumerGroupName
	client, err := sarama.NewConsumerGroup(addrs, consumerGroupName, conf)
	if err != nil {
		logConf.Logger.LogError(ctx, logConf.Category, map[string]interface{}{
			"err":     err,
			"topics":  topics,
			"address": k.consumerAddrs,
		}, "Consumer failed")
		panic(fmt.Sprintf("创建消费者分组失败, topics: %v, err: %s", topics, err.Error()))
	}
	defer client.Close()
	handler := consumerGroupHandler{Kafka: k} // 必须传递一个handler
	// 使用链路追踪消费者
	//consumerHandler := consumerGroupHandler{Kafka: k}
	//handler := otelsarama.WrapConsumerGroupHandler(&consumerHandler)
	for { // for循环的目的是因为存在重平衡，他会重新启动
		select {
		case <-ctx.Done():
			break
		default:
		}
		err = client.Consume(ctx, topics, handler) // consume 操作，死循环。exampleConsumerGroupHandler的ConsumeClaim不允许退出，也就是操作到完毕。
		if err != nil {
			logConf.Logger.LogError(ctx, logConf.Category, map[string]interface{}{
				"topics":  topics,
				"err":     err,
				"address": k.consumerAddrs,
			}, "msg consumer failed")
		}
	}
}

// handler，核心的消费者业务实现
type consumerGroupHandler struct {
	Kafka
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
	ctx := sess.Context()
	//ctx := context.WithValue(context.Background(), "logCategory", logConf.Category)
	for msg := range claim.Messages() {
		select {
		case <-sess.Context().Done():
			break
		default:
		}
		tmpMsg := msg
		newCtx := ctx
		// 从消息头部中取traceId 和msgId 写到上下文中
		for _, v := range tmpMsg.Headers {
			headerKey := string(v.Key)
			if _, ok := ctxWithMap[headerKey]; ok {
				newCtx = context.WithValue(newCtx, headerKey, string(v.Value))
			}
		}
		highWaterMarkOffset := claim.HighWaterMarkOffset()
		_err := h.goPool.Submit(func() {
			// Extract tracing info from message
			//newCtx = otel.GetTextMapPropagator().Extract(newCtx, otelsarama.NewConsumerMessageCarrier(tmpMsg))
			//
			//_, span := otel.Tracer("consumer").Start(newCtx, "consume message", trace.WithAttributes(
			//	semconv.MessagingOperationProcess,
			//))
			//defer span.End()
			// 业务逻辑处理
			err := h.callback(newCtx, tmpMsg)
			if logConf.Producer || err != nil {
				logMap := map[string]interface{}{
					"topic":        tmpMsg.Topic,
					"group":        h.Kafka.group,
					"partition":    tmpMsg.Partition,
					"offset":       tmpMsg.Offset,
					"maxOffsetSub": highWaterMarkOffset - tmpMsg.Offset,
					"key":          string(tmpMsg.Key),
					"value":        h.cutStrFromLogConfig(string(tmpMsg.Value)),
				}
				if err != nil {
					logMap["err"] = err
					logMap["address"] = h.consumerAddrs
					logConf.Logger.LogError(newCtx, logConf.Category, logMap, "[Consumer] Message Failed")
					// 扔到重试队列或死信队列
				} else if logConf.Producer {
					logConf.Logger.LogInfo(newCtx, logConf.Category, logMap, "[Consumer] Message Success")
				}
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
