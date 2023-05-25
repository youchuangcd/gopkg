package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	"go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama/otelsarama"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

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
	// Create root span
	tr := otel.Tracer("producer")
	ctx, span := tr.Start(ctx, "produce message")
	defer span.End()

	msg := &sarama.ProducerMessage{
		Topic:   topic,
		Value:   sarama.StringEncoder(content),
		Headers: headers,
	}
	// Inject tracing info into message
	otel.GetTextMapPropagator().Inject(ctx, otelsarama.NewProducerMessageCarrier(msg))
	// 发送消息
	partition, offset, err = k.syncProducer.SendMessage(msg)
	//content = k.cutStrFromLogConfig(content)
	if err != nil {
		// 标记span状态
		span.SetStatus(codes.Error, err.Error())
		//content = k.cutStrFromLogConfig(content)
		logConf.Logger.LogError(ctx, logConf.Category, map[string]interface{}{
			"err":     err,
			"topic":   topic,
			"msg":     msg,
			"address": k.producerAddrs,
		}, "send msg failed")
		return
	}
	// 是否开启生产者日志
	if logConf.Producer {
		logConf.Logger.LogInfo(ctx, logConf.Category, map[string]interface{}{
			"partition": partition,
			"offset":    offset,
			"msg":       msg,
		}, "send msg success")
	}
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
	traceId, _ := ctx.Value(ctxTraceIdKey).(string)
	// Create root span
	tr := otel.Tracer("producer")
	ctx, span := tr.Start(ctx, "produce batch message")
	defer span.End()
	for _, content := range contents {
		headers := make([]sarama.RecordHeader, 0, 2)
		if traceId != "" {
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
		// Inject tracing info into message
		otel.GetTextMapPropagator().Inject(ctx, otelsarama.NewProducerMessageCarrier(msg))
		msgs = append(msgs, msg)
	}

	// 发送消息
	err = k.syncProducer.SendMessages(msgs)
	if err != nil {
		logConf.Logger.LogError(ctx, logConf.Category, map[string]interface{}{
			"msgs":    msgs,
			"err":     err,
			"topic":   topic,
			"address": k.producerAddrs,
		}, "send msg failed")
		return
	}
	// 是否开启生产者日志
	if logConf.Producer {
		logConf.Logger.LogInfo(ctx, logConf.Category, map[string]interface{}{
			"msgs":  msgs,
			"topic": topic,
		}, "send msg success")
	}
	return
}
