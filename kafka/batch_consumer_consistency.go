package kafka

import (
	"context"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/youchuangcd/gopkg"
	"github.com/youchuangcd/gopkg/common"
	"time"
)

// BatchConsumerConsistency
//
//	@Description: 强一致性批处理：只有1个协程消费处理，成功才提交偏移量
//	@receiver k
//	@param ctx
//	@param batchConf
func (k Kafka) BatchConsumerConsistency(ctx context.Context, batchConf BatchConsumerConfig) {
	if batchConf.ConsumerConfig.ConsumerOffsets != 0 {
		k.consumerOffsets = batchConf.ConsumerConfig.ConsumerOffsets
	}
	conf := k.getConfig()
	// 手动提交消费偏移量
	conf.Consumer.Offsets.AutoCommit.Enable = false
	// 没有额外设置地址，取配置地址
	addrs := k.getConsumerAddr()
	if common.EnvLocal() || common.EnvDev() { // 开发环境会追加环境变量，与其他环境隔开
		batchConf.ConsumerGroupName += "_" + gopkg.EnvDev
	}
	var err error
	k.group = batchConf.ConsumerGroupName
	// k.batchProcess里面用到
	k.callbackBatchProcess = batchConf.Callback
	// 创建批处理聚合对象
	k.aggregator = common.NewAggregator(ctx, k.batchProcessConsistency, func(option common.AggregatorOption) common.AggregatorOption {
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
	//consumerHandler := batchConsumerGroupHandler{Kafka: k}
	//handler := otelsarama.WrapConsumerGroupHandler(&consumerHandler)
Loop:
	for { // for循环的目的是因为存在重平衡，他会重新启动
		select {
		case <-ctx.Done():
			break Loop
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
//	@Description: 批处理消息，把any转成ConsumerMessage
//	@receiver k
//	@param items
//	@return error
func (k Kafka) batchProcessConsistency(ctx context.Context, items []any) (err error) {
	msgs := make([]*sarama.ConsumerMessage, 0, len(items))
	//spanSlice := make([]trace.Span, 0, len(items))
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
			// Extract tracing info from message
			//newCtx := otel.GetTextMapPropagator().Extract(ctx, otelsarama.NewConsumerMessageCarrier(msgExt.msg))
			//
			//_, span := otel.Tracer("consumer").Start(newCtx, "consume batch message", trace.WithAttributes(
			//	semconv.MessagingOperationProcess,
			//))
			//spanSlice = append(spanSlice, span)
		}
	}
	if len(msgs) == 0 {
		return errors.New("无效的消息类型")
	}

	//defer func() {
	//	for _, span := range spanSlice {
	//		span.End()
	//	}
	//}()
	err = k.callbackBatchProcess(ctx, msgs)
	if err != nil || logConf.Consumer {
		logMsg := "[BatchConsumer] Message Success"
		msg := msgs[len(msgs)-1]
		logFunc := logConf.Logger.LogInfo
		// 从消息头部中取traceId 和msgId 写到上下文中
		for _, v := range msg.Headers {
			headerKey := string(v.Key)
			if _, ok2 := ctxWithMap[headerKey]; ok2 {
				ctx = context.WithValue(ctx, headerKey, string(v.Value))
			}
		}
		logMap := map[string]any{
			"topic":     msg.Topic,
			"group":     k.group,
			"partition": msg.Partition,
			"offset":    msg.Offset,
			"key":       string(msg.Key),
			"value":     k.cutStrFromLogConfig(string(msg.Value)),
		}
		if err != nil {
			logMsg = "[BatchConsumer] Message Failed"
			logMap["err"] = err
			logMap["address"] = k.consumerAddrs
			logFunc = logConf.Logger.LogError
		}
		logFunc(ctx, logConf.Category, logMap, logMsg)
	}
	if err == nil {
		lastMsgExt := items[len(items)-1]
		msgExt, _ := lastMsgExt.(batchConsumerMessageExt)
		msgExt.sess.Commit()
	}
	return
}
