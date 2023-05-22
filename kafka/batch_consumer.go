package kafka

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/youchuangcd/gopkg/common"
	"time"
)

type BatchConsumerConfig struct {
	Topics            []string
	ConsumerGroupName string
	Callback          func(ctx context.Context, msgs []*sarama.ConsumerMessage) error
	BatchSize         int // 达到多少条处理一次
	GoPoolSize        int
	LingerTime        int64 // 多久处理一次 单位毫秒
	ConsumerConfig    ConsumerConfig
}

// 消息批量处理handler，核心的消费者业务实现
type batchConsumerGroupHandler struct {
	Kafka
}

func (h batchConsumerGroupHandler) Setup(s sarama.ConsumerGroupSession) error {
	// 当连接完毕的时候会通知这个，start
	return nil
}
func (h batchConsumerGroupHandler) Cleanup(s sarama.ConsumerGroupSession) error {
	s.Commit()
	// end，当这一次消费完毕，会通知，这里最好commit
	return nil
}
func (h batchConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error { // consume
	ctx := sess.Context()
	// 当第一个ConsumeClaim消费完成，会话就会被关闭
	//ctx := context.WithValue(context.Background(), "logCategory", logConf.Category)
	for msg := range claim.Messages() {
		select {
		case <-sess.Context().Done():
			break
		default:
		}
		// 丢进内存队列中
		h.Kafka.aggregator.TryEnqueue(msg)
		if logConf.Producer {
			newCtx := ctx
			// 从消息头部中取traceId 和msgId 写到上下文中
			for _, v := range msg.Headers {
				headerKey := string(v.Key)
				if _, ok := ctxWithMap[headerKey]; ok {
					newCtx = context.WithValue(newCtx, headerKey, string(v.Value))
				}
			}
			highWaterMarkOffset := claim.HighWaterMarkOffset()
			logMap := map[string]interface{}{
				"topic":        msg.Topic,
				"group":        h.Kafka.group,
				"partition":    msg.Partition,
				"offset":       msg.Offset,
				"maxOffsetSub": highWaterMarkOffset - msg.Offset,
				"key":          string(msg.Key),
				"value":        h.cutStrFromLogConfig(string(msg.Value)),
			}
			logConf.Logger.LogInfo(newCtx, logConf.Category, logMap, "[BatchConsumer] Message Success")
		}
		sess.MarkMessage(msg, "") // 必须设置这个，不然你的偏移量无法提交。
	}
	return nil
}

// BatchConsumer
//
//	@Description: 消费批量处理
//	@receiver k
//	@param ctx
//	@param batchConf
func (k Kafka) BatchConsumer(ctx context.Context, batchConf BatchConsumerConfig) {
	if batchConf.ConsumerConfig.ConsumerOffsets != 0 {
		k.consumerOffsets = batchConf.ConsumerConfig.ConsumerOffsets
	}
	conf := k.getConfig()
	// 没有额外设置地址，取配置地址
	addrs := k.getConsumerAddr()
	var err error
	k.group = batchConf.ConsumerGroupName
	// k.batchProcess里面用到
	k.callbackBatchProcess = batchConf.Callback
	// 创建批处理聚合对象
	k.aggregator = common.NewAggregator(ctx, k.batchProcess, func(option common.AggregatorOption) common.AggregatorOption {
		option.BatchSize = batchConf.BatchSize
		option.Workers = batchConf.GoPoolSize
		option.LingerTime = time.Duration(batchConf.LingerTime) * time.Millisecond // 多久提交一次 单位毫秒
		option.Logger = logConf.Logger
		return option
	})

	client, err := sarama.NewConsumerGroup(addrs, batchConf.ConsumerGroupName, conf)
	if err != nil {
		logConf.Logger.LogError(ctx, logConf.Category, map[string]interface{}{
			"err":     err,
			"topics":  batchConf.Topics,
			"address": k.consumerAddrs,
		}, "Consumer failed")
		panic(fmt.Sprintf("创建消费者分组失败, topics: %v, err: %s", batchConf.Topics, err.Error()))
	}
	// 开启监听通道协程
	k.aggregator.Start()
	defer func() {
		k.aggregator.SafeStop()
		client.Close()
	}()
	handler := batchConsumerGroupHandler{Kafka: k} // 必须传递一个handler
	for {                                          // for循环的目的是因为存在重平衡，他会重新启动
		select {
		case <-ctx.Done():
			break
		default:
		}
		err = client.Consume(ctx, batchConf.Topics, handler) // consume 操作，死循环。exampleConsumerGroupHandler的ConsumeClaim不允许退出，也就是操作到完毕。
		if err != nil {
			logConf.Logger.LogError(ctx, logConf.Category, map[string]interface{}{
				"topics":  batchConf.Topics,
				"err":     err,
				"address": k.consumerAddrs,
			}, "msg consumer failed")
		}
	}
}

// batchProcess
//
//	@Description: 批处理，把any转成ConsumerMessage
//	@receiver k
//	@param items
//	@return error
func (k Kafka) batchProcess(ctx context.Context, items []any) error {
	msgs := make([]*sarama.ConsumerMessage, 0, len(items))
	for _, item := range items {
		if msg, ok := item.(*sarama.ConsumerMessage); ok {
			msgs = append(msgs, msg)
		}
	}
	return k.callbackBatchProcess(ctx, msgs)
}
