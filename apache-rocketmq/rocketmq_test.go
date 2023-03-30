package apacherocketmq

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"os"
	"strconv"
	"sync"
	"testing"
)

var address = []string{"123.60.185.100:8200", "123.60.76.68:8200"}

//var address = []string{"121.36.207.59:8200"}

func TestMq(t *testing.T) {
	topicName := "paas-data-center-misc-order"
	p, _ := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver(address)),
		producer.WithRetry(2),
		producer.WithNamespace("dev"),
	)
	err := p.Start()
	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}

	for i := 0; i < 100; i++ {
		msg := &primitive.Message{
			Topic: topicName,
			Body:  []byte("Hello RocketMQ Go Client! " + strconv.Itoa(i)),
		}
		orderId := strconv.Itoa(i % 10)
		// 相同orderId的消息需要保证顺序，不同orderId的消息不需要保证顺序，所以将orderId作为选择队列的sharding key。
		msg.WithShardingKey(orderId)
		res, err := p.SendSync(context.Background(), msg)

		if err != nil {
			fmt.Printf("send message error: %s\n", err)
		} else {
			fmt.Printf("send message success: result=%s\n", res.String())
		}
	}
	defer func() {
		err = p.Shutdown()
		if err != nil {
			fmt.Printf("shutdown producer error: %s", err.Error())
		}
	}()
}

func TestConsumer(t *testing.T) {
	topicName := "paas-data-center-misc-order"
	groupId := "gid-paas-data-center-misc-order"
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithGroupName(groupId),
		consumer.WithNsResolver(primitive.NewPassthroughResolver(address)),
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromFirstOffset),
		consumer.WithNamespace("dev"),
	)
	tmp := sync.Map{}
	err := c.Subscribe(topicName, consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		//orderlyCtx, _ := primitive.GetOrderlyCtx(ctx)
		for _, msg := range msgs {
			content := string(msg.Body)
			if _, ok := tmp.Load(content); ok {
				fmt.Println("重复的消息：", msg.Body)
			}
			tmp.Store(content, "存在")
		}
		//fmt.Printf("orderly context: %v\n", orderlyCtx)
		fmt.Printf("subscribe orderly callback: %v \n", msgs)
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	// Note: start after subscribe
	err = c.Start()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	defer func() {
		err = c.Shutdown()
		if err != nil {
			fmt.Printf("Shutdown Consumer error: %s", err.Error())
		}
	}()
	<-(chan interface{})(nil)
}
