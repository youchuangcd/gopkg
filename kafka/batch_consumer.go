package kafka

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/gogap/errors"
	"github.com/youchuangcd/gopkg/common"
	"time"
)

type BatchConsumerConfig struct {
	Topics            []string
	ConsumerGroupName string
	Callback          func(ctx context.Context, msgs []*sarama.ConsumerMessage) error
	BatchSize         int // 达到多少条处理一次
	ChannelBufferSize int // 缓冲通道大小
	GoPoolSize        int
	LingerTime        int64 // 多久处理一次 单位毫秒
	ConsumerConfig    ConsumerConfig
}

type batchConsumerMessageExt struct {
	ctx  context.Context
	sess sarama.ConsumerGroupSession
	msg  *sarama.ConsumerMessage
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
		newMsg := msg
		msgExt := batchConsumerMessageExt{
			ctx:  ctx,
			sess: sess,
			msg:  newMsg,
		}
		for {
			// 丢进内存队列中
			res := h.Kafka.aggregator.TryEnqueue(msgExt)
			if res {
				// 标记消息偏移量
				sess.MarkMessage(msg, "")
				break
			}
		}
	}
	return nil
}

// batchProcess
//
//	@Description: 批处理消息，把any转成ConsumerMessage
//	@receiver k
//	@param items
//	@return error
func (k Kafka) batchProcess(ctx context.Context, items []any) (err error) {
	msgs := make([]*sarama.ConsumerMessage, 0, len(items))
	for _, item := range items {
		if msgExt, ok := item.(batchConsumerMessageExt); ok {
			select {
			case <-ctx.Done(): // 程序退出
				return nil
			case <-msgExt.sess.Context().Done(): // kafka消费者会话退出
				return nil
			default:
			}
			msgs = append(msgs, msgExt.msg)
		}
	}
	if len(msgs) == 0 {
		return errors.New("无效的消息类型")
	}
	err = k.callbackBatchProcess(ctx, msgs)
	if err == nil {
		if logConf.Producer {
			for _, item := range items {
				if msgExt, ok := item.(batchConsumerMessageExt); ok {
					msg := msgExt.msg
					newCtx := msgExt.ctx
					// 从消息头部中取traceId 和msgId 写到上下文中
					for _, v := range msg.Headers {
						headerKey := string(v.Key)
						if _, ok2 := ctxWithMap[headerKey]; ok2 {
							newCtx = context.WithValue(newCtx, headerKey, string(v.Value))
						}
					}
					logMap := map[string]interface{}{
						"topic":     msg.Topic,
						"group":     k.group,
						"partition": msg.Partition,
						"offset":    msg.Offset,
						"key":       string(msg.Key),
						"value":     k.cutStrFromLogConfig(string(msg.Value)),
					}
					logConf.Logger.LogInfo(newCtx, logConf.Category, logMap, "[BatchConsumer] Message Success")
				}
			}
		}
		lastMsgExt := items[len(items)-1]
		msgExt, _ := lastMsgExt.(batchConsumerMessageExt)
		msgExt.sess.Commit()
	}
	return
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
	// 手动提交消费偏移量
	conf.Consumer.Offsets.AutoCommit.Enable = false
	// 没有额外设置地址，取配置地址
	addrs := k.getConsumerAddr()
	var err error
	k.group = batchConf.ConsumerGroupName
	// k.batchProcess里面用到
	k.callbackBatchProcess = batchConf.Callback
	// 创建批处理聚合对象
	k.aggregator = common.NewAggregator(ctx, k.batchProcess, func(option common.AggregatorOption) common.AggregatorOption {
		option.BatchSize = batchConf.BatchSize
		//option.Workers = batchConf.GoPoolSize
		// 因为偏移量的关系，暂时只支持单协程消费一批消息
		option.Workers = 1
		option.LingerTime = time.Duration(batchConf.LingerTime) * time.Millisecond // 多久提交一次 单位毫秒
		option.Logger = logConf.Logger
		option.ChannelBufferSize = batchConf.ChannelBufferSize // 设置缓冲通道大小
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
