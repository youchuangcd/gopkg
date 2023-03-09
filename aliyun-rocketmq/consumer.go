package aliyunrocketmq

import (
	"context"
	"fmt"
	mqhttpsdk "github.com/aliyunmq/mq-http-go-sdk"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"github.com/gogap/errors"
	"github.com/youchuangcd/gopkg"
	"github.com/youchuangcd/gopkg/common/utils"
	"github.com/youchuangcd/gopkg/mylog"
	"strings"
	"sync"
	"time"
)

// ConsumerMessage
// @Description: 消费消息
// @receiver p
// @param ctx
// @param instanceId
// @param topicName
// @param groupId
// @param tag
// @param callback 业务逻辑处理 返回bool值表示消费成功，否则就不ACK
// @return err
func (p *RocketMQ) ConsumerMessage(ctx context.Context, instanceId, topicName, groupId, tag string, callback func(ctx context.Context, msg *mqhttpsdk.ConsumeMessageEntry) (ok bool)) {
	mqConsumer := p.consumerClient.GetConsumer(instanceId, topicName, groupId, tag)

	for {
		endChan := make(chan int)
		respChan := make(chan mqhttpsdk.ConsumeMessageResponse)
		errChan := make(chan error)
		p.goPool.Submit(func() {
			select {
			case resp := <-respChan:
				{
					var (
						handles []string
						wg      sync.WaitGroup
						lock    sync.Mutex
					)
					// 处理业务逻辑
					wg.Add(len(resp.Messages))
					for _, v := range resp.Messages {
						newCtx := context.WithValue(ctx, gopkg.RequestHeaderTraceIdKey, v.MessageKey)
						newCtx = context.WithValue(newCtx, gopkg.LogMsgIdKey, v.MessageId)
						// 处理成功才ack
						msg := v
						p.goPool.Submit(func() {
							defer wg.Done()
							logMsgBody := utils.CutStrFromLogConfig(msg.MessageBody)
							logWithField := map[string]interface{}{
								"topicName":          topicName,
								"groupId":            groupId,
								"process_result":     false,
								"code":               msg.Code,
								"message_id":         msg.MessageId,
								"request_id":         msg.RequestId,
								"receipt_handle":     msg.ReceiptHandle,
								"message_body_md5":   msg.MessageBodyMD5,
								"message_body":       logMsgBody,
								"publish_time":       msg.PublishTime,
								"next_consume_time":  msg.NextConsumeTime,
								"first_consume_time": msg.FirstConsumeTime,
								"consumed_times":     msg.ConsumedTimes,
								"message_tag":        msg.MessageTag,
							}
							if ok := callback(newCtx, &msg); ok {
								logWithField["process_result"] = ok
								lock.Lock()
								handles = append(handles, msg.ReceiptHandle)
								lock.Unlock()
							}
							mylog.WithInfo(newCtx, gopkg.LogRocketMQ, logWithField, "处理完一条消息")
						})
					}
					wg.Wait()

					if len(handles) > 0 {
						// NextConsumeTime前若不确认消息消费成功，则消息会被重复消费。
						// 消息句柄有时间戳，同一条消息每次消费拿到的都不一样。
						ackerr := mqConsumer.AckMessage(handles)
						if ackerr != nil {
							// 某些消息的句柄可能超时，会导致消息消费状态确认不成功。
							//mylog.Error(ctx, gopkg.RocketMQLog, fmt.Sprintf("%+v", ackerr))
							if errAckItems, ok := ackerr.(errors.ErrCode).Context()["Detail"].([]mqhttpsdk.ErrAckItem); ok {
								for _, errAckItem := range errAckItems {
									mylog.Error(ctx, gopkg.LogRocketMQ, fmt.Sprintf("\tErrorHandle:%s, ErrorCode:%s, ErrorMsg:%s\n",
										errAckItem.ErrorHandle, errAckItem.ErrorCode, errAckItem.ErrorMsg))
								}
							} else {
								mylog.Error(ctx, gopkg.LogRocketMQ, fmt.Sprintf("ack err = %+v", ackerr))
							}
							time.Sleep(time.Duration(3) * time.Second)
						} else {
							mylog.Info(ctx, gopkg.LogRocketMQ, fmt.Sprintf("Ack ---->%s", handles))
						}
					}

					endChan <- 1
				}
			case err := <-errChan:
				{
					// Topic中没有消息可消费。
					if strings.Contains(err.(errors.ErrCode).Error(), "MessageNotExist") {
						//mylog.Info(ctx, gopkg.RocketMQLog, "No new message, continue!")
						//fmt.Println("\nNo new message, continue!")
					} else {
						mylog.Error(ctx, gopkg.LogRocketMQ, fmt.Sprintf("%+v", err))
						//fmt.Println(err)
						time.Sleep(time.Duration(3) * time.Second)
					}
					endChan <- 1
				}
			case <-time.After(35 * time.Second):
				{
					mylog.Info(ctx, gopkg.LogRocketMQ, "Timeout of consumer message ??")
					//fmt.Println("Timeout of consumer message ??")
					endChan <- 1
				}
			}
		})

		// 长轮询消费消息，网络超时时间默认为35s。
		// 长轮询表示如果Topic没有消息，则客户端请求会在服务端挂起3s，3s内如果有消息可以消费则立即返回响应。
		mqConsumer.ConsumeMessage(respChan, errChan,
			p.maxConsumeLimit, // 一次最多消费3条（最多可设置为16条）。
			p.longLoopTime,    // 长轮询时间3s（最多可设置为30s）。
		)
		<-endChan
	}
}

func InitTcpLog() {
	rlog.SetLogger(mylog.GetSpecialLogger())
}

func (p *RocketMQ) BroadcastConsumerMessage(ctx context.Context, instanceId, topicName, groupId string, selector consumer.MessageSelector, f func(context.Context, ...*primitive.MessageExt) (consumer.ConsumeResult, error)) {
	credentials := primitive.Credentials{
		AccessKey:     p.accessKeyId,
		SecretKey:     p.accessKeySecret,
		SecurityToken: p.securityToken,
	}
	isFirstRun := true
	c, err := rocketmq.NewPushConsumer(
		consumer.WithInstance(instanceId),
		consumer.WithNamespace(instanceId),
		consumer.WithGroupName(groupId),
		consumer.WithCredentials(credentials),
		consumer.WithNameServer([]string{p.endPoint}),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromLastOffset),
		consumer.WithConsumerModel(consumer.BroadCasting),
		consumer.WithTrace(&primitive.TraceConfig{
			GroupName:    groupId,
			Access:       primitive.Cloud,
			NamesrvAddrs: []string{p.endPoint},
			Credentials:  credentials,
		}),
	)
	if err != nil {
		mylog.Error(ctx, gopkg.LogRocketMQTCP, fmt.Sprintf("TCP客户端连接服务端失败, err: %s", err.Error()))
		panic("无法连接RocketMQ，err: " + err.Error())
	}
	defer func() {
		err = c.Shutdown()
		if err != nil {
			mylog.Error(ctx, gopkg.LogRocketMQTCP, fmt.Sprintf("TCP客户端shutdown失败, err: %s", err.Error()))
			return
		}
	}()
	err = c.Subscribe(topicName, selector, f)
	if err != nil {
		mylog.Error(ctx, gopkg.LogRocketMQTCP, fmt.Sprintf("TCP客户端订阅topic失败, err: %s", err.Error()))
		panic("无法订阅RocketMQ " + topicName + "，err: " + err.Error())
	}
	for {
		// Note: start after subscribe
		err = c.Start()
		if err != nil {
			mylog.Error(ctx, gopkg.LogRocketMQTCP, fmt.Sprintf("TCP客户端start失败, err: %s", err.Error()))
			if isFirstRun {
				panic("TCP客户端start失败，err: " + err.Error())
			} else {
				time.Sleep(10 * time.Second)
				continue
			}
		}
		time.Sleep(30 * time.Second)
		//time.Sleep(time.Hour)
		isFirstRun = false
	}
}
