package aliyunrocketmq

import (
	"context"
	"github.com/youchuangcd/gopkg"
	"github.com/youchuangcd/gopkg/common/utils"
	mqhttpsdk "github.com/youchuangcd/gopkg/mq-http-go-sdk"
	"github.com/youchuangcd/gopkg/mylog"
	"time"
)

// PushMessage
// @Description:
// @receiver p
// @param ctx
// @param instanceId
// @param topicName
// @param tag
// @param body
// @param startDeliverTime 延时消息，发送时间为10s后。该参数格式为毫秒级别的时间戳。若发送定时消息，设置该参数时需要计算定时时间与当前时间的时间差。
// @param messageKey
// @param properties
// @return err
func (p *RocketMQ) PushMessage(ctx context.Context, instanceId, topicName, tag, body string, startDeliverTime int64, messageKey string, properties map[string]string) (err error) {
	//client := p.clients.Get().(mqhttpsdk.MQClient)
	//defer p.clients.Put(client)
	//mqProducer := client.GetProducer(instanceId, topicName) // topicName拼接环境变量到最后，用来在线上实例区分环境
	mqProducer := p.producerClient.GetProducer(instanceId, topicName) // topicName拼接环境变量到最后，用来在线上实例区分环境
	if messageKey == "" {
		// 用traceId来做MessageKey 可以追踪整个链路
		if traceId, ok := ctx.Value(gopkg.RequestHeaderTraceIdKey).(string); ok {
			messageKey = traceId
		} else {
			messageKey = utils.GenUniqueId()
		}
	}
	msg := mqhttpsdk.PublishMessageRequest{
		MessageBody: body,       //消息内容。
		MessageTag:  tag,        // 消息标签。
		Properties:  properties, // 消息属性。
		MessageKey:  messageKey,
	}
	// 延时消息，发送时间为10s后。该参数格式为毫秒级别的时间戳。
	//msg.StartDeliverTime = time.Now().UTC().Unix() * 1000 + 10 * 1000
	msgType := "普通"
	if startDeliverTime != 0 {
		// 如果设置的值不是毫秒级时间戳，就默认填充
		if startDeliverTime < 1000 {
			startDeliverTime = time.Now().UnixNano()/1e6 + startDeliverTime*1000
		}
		msg.StartDeliverTime = startDeliverTime
		msgType = "延迟"
	}
	resp, err := mqProducer.PublishMessage(msg)
	msg.MessageBody = utils.CutStrFromLogConfig(msg.MessageBody)
	logContent := map[string]interface{}{
		"param":            msg,
		"startDeliverTime": msg.StartDeliverTime,
		"resp":             resp,
	}
	if err != nil {
		logContent["err"] = err
		mylog.WithError(ctx, gopkg.LogRocketMQ, logContent, "推送"+msgType+"消息失败")
	} else {
		//mylog.WithInfo(ctx, gopkg.LogRocketMQ, logContent, "推送"+msgType+"消息成功")
	}
	return
}

//
//func (p *RocketMQ) PushTransMessage(ctx context.Context, instanceId, topicName, tag, body string, transCheckImmunityTime int, messageKey string, properties map[string]string) (err error) {
//	transCheckImmunityTime = int(utils.Max(10, utils.Min(300, int64(transCheckImmunityTime)))) // 10-300
//	mqProducer := p.producerClient.GetProducer(instanceId, topicName) // topicName拼接环境变量到最后，用来在线上实例区分环境
//	if messageKey == "" {
//		// 用traceId来做MessageKey 可以追踪整个链路
//		if traceId, ok := ctx.Value(global.RequestHeaderTraceIdKey).(string); ok {
//			messageKey = traceId
//		} else {
//			messageKey = utils.GenUniqueId()
//		}
//	}
//	msg := mqhttpsdk.PublishMessageRequest{
//		MessageBody: body,       //消息内容。
//		MessageTag:  tag,        // 消息标签。
//		Properties:  properties, // 消息属性。
//		MessageKey:  messageKey,
//		TransCheckImmunityTime: transCheckImmunityTime,
//	}
//	msgType := "事务"
//	resp, err := mqProducer.PublishMessage(msg)
//	msg.MessageBody = utils.CutStrFromLogConfig(msg.MessageBody)
//	logContent := map[string]interface{}{
//		"param":            msg,
//		"transCheckImmunityTime": transCheckImmunityTime,
//		"resp":             resp,
//	}
//	if err != nil {
//		logContent["err"] = err
//		mylog.WithError(ctx, global.RocketMQLog, logContent, "推送"+msgType+"消息失败")
//	} else {
//		mylog.WithInfo(ctx, global.RocketMQLog, logContent, "推送"+msgType+"消息成功")
//	}
//	return
//}
//
//
//func (p *RocketMQ) ProcessError(ctx context.Context, err error) {
//	// 如果Commit或Rollback时超过了TransCheckImmunityTime（针对发送事务消息的句柄）或者超过10s（针对consumeHalfMessage的句柄），则Commit或Rollback失败。
//	if err == nil {
//		return
//	}
//	for _, errAckItem := range err.(errors.ErrCode).Context()["Detail"].([]mqhttpsdk.ErrAckItem) {
//		mylog.WithError(ctx, global.RocketMQLog, map[string]interface{}{
//			"err": errAckItem,
//		}, "事务消息错误信息")
//	}
//}
//
//
//func (p *RocketMQ) ConsumeHalfMsg(ctx context.Context, mqTransProducer *mqhttpsdk.MQTransProducer) {
//	loopCount := 0
//	for {
//		if loopCount >= 10 {
//			return
//		}
//		loopCount++
//		endChan := make(chan int)
//		respChan := make(chan mqhttpsdk.ConsumeMessageResponse)
//		errChan := make(chan error)
//		go func() {
//			select {
//			case resp := <-respChan:
//				{
//					// 处理业务逻辑。
//					var handles []string
//					for _, v := range resp.Messages {
//						handles = append(handles, v.ReceiptHandle)
//
//
//						a, _ := strconv.Atoi(v.Properties["id"])
//						var comRollErr error
//						if a == 1 {
//							// 确认提交事务消息。
//							comRollErr = (*mqTransProducer).Commit(v.ReceiptHandle)
//							fmt.Println("Commit---------->")
//						} else if a == 2 && v.ConsumedTimes > 1 {
//							// 确认提交事务消息。
//							comRollErr = (*mqTransProducer).Commit(v.ReceiptHandle)
//							fmt.Println("Commit---------->")
//						} else if a == 3 {
//							// 确认回滚事务消息。
//							comRollErr = (*mqTransProducer).Rollback(v.ReceiptHandle)
//							fmt.Println("Rollback---------->")
//						} else {
//							// 什么都不做，下次再检查。
//							fmt.Println("Unknown---------->")
//						}
//						p.ProcessError(ctx, comRollErr)
//					}
//					endChan <- 1
//				}
//			case err := <-errChan:
//				{
//					// Topic中没有消息可消费。
//					if strings.Contains(err.(errors.ErrCode).Error(), "MessageNotExist") {
//						fmt.Println("\nNo new message, continue!")
//					} else {
//						fmt.Println(err)
//						time.Sleep(time.Duration(3) * time.Second)
//					}
//					endChan <- 1
//				}
//			case <-time.After(35 * time.Second):
//				{
//					fmt.Println("Timeout of consumer message ??")
//					return
//				}
//			}
//		}()
//
//		// 长轮询检查半事务消息。
//		// 长轮询表示如果Topic没有消息则请求会在服务端挂起3s，3s内如果有消息可以消费则立即返回响应。
//		(*mqTransProducer).ConsumeHalfMessage(respChan, errChan,
//			p.MaxConsumeLimit, // 一次最多消费3条（最多可设置为16条）。
//			p.LongLoopTime,    // 长轮询时间3s（最多可设置为30s）。
//		)
//		<-endChan
//	}
//}
