package common

//这个方案需要实现以下几点：
//1.消息聚合后处理（最大条数为BatchSize），核心：
//（1）带buffer的channel相当于一个FIFO的队列
//（2）多个常驻的goroutine来提高并发
//（3）goroutine之间是并行的，但每个goroutine内是串行的，所以对batch操作是不用加锁的。
//2.延迟处理（延迟时间为LingerTime）
//  注意：为什么使用time.Timer而不是time.After，是因为time.After在for select中使用时，会发生内存泄露。
//3.自定义错误处理
//4.并发处理
import (
	"context"
	"github.com/youchuangcd/gopkg"
	"runtime"
	"sync"
	"time"
)

type logger interface {
	LogDebug(ctx context.Context, logCategory string, logContent map[string]interface{}, msg string)
	LogInfo(ctx context.Context, logCategory string, logContent map[string]interface{}, msg string)
	LogWarn(ctx context.Context, logCategory string, logContent map[string]interface{}, msg string)
	LogError(ctx context.Context, logCategory string, logContent map[string]interface{}, msg string)
}

// Represents the aggregator
type Aggregator struct {
	ctx            context.Context
	option         AggregatorOption
	wg             *sync.WaitGroup
	quit           chan struct{}
	eventQueue     chan any
	batchProcessor BatchProcessFunc
}

// Represents the aggregator option
type AggregatorOption struct {
	BatchSize         int
	Workers           int
	ChannelBufferSize int
	LingerTime        time.Duration
	ErrorHandler      ErrorHandlerFunc
	Logger            logger
}

// the func to batch process items
type BatchProcessFunc func(ctx context.Context, items []any) error

// the func to set option for aggregator
type SetAggregatorOptionFunc func(option AggregatorOption) AggregatorOption

// the func to handle error
type ErrorHandlerFunc func(ctx context.Context, err error, items []any, batchProcessFunc BatchProcessFunc, aggregator *Aggregator)

// Creates a new aggregator
func NewAggregator(ctx context.Context, batchProcessor BatchProcessFunc, optionFuncs ...SetAggregatorOptionFunc) *Aggregator {
	option := AggregatorOption{
		BatchSize:  8,
		Workers:    runtime.NumCPU(),
		LingerTime: 1 * time.Minute,
	}

	for _, optionFunc := range optionFuncs {
		option = optionFunc(option)
	}

	if option.ChannelBufferSize <= option.Workers {
		option.ChannelBufferSize = option.Workers
	}
	if option.Logger == nil {
		panic("请配置Aggregator日志组件")
	}

	return &Aggregator{
		ctx:            ctx,
		eventQueue:     make(chan any, option.ChannelBufferSize),
		option:         option,
		quit:           make(chan struct{}),
		wg:             new(sync.WaitGroup),
		batchProcessor: batchProcessor,
	}
}

// Try enqueue an item, and it is non-blocked
func (agt *Aggregator) TryEnqueue(item any) bool {
	select {
	case agt.eventQueue <- item:
		return true
	default:
		agt.option.Logger.LogWarn(agt.ctx, gopkg.LogCtx, map[string]any{
			"item": item,
		}, "[Aggregator] Event queue is full and try reschedule")

		runtime.Gosched()

		select {
		case agt.eventQueue <- item:
			return true
		default:
			agt.option.Logger.LogWarn(agt.ctx, gopkg.LogCtx, map[string]any{
				"item": item,
			}, "[Aggregator] Event queue is still full and is skipped.")
			return false
		}
	}
}

// Enqueue an item, will be blocked if the queue is full
func (agt *Aggregator) Enqueue(item any) {
	agt.eventQueue <- item
}

// Start the aggregator
func (agt *Aggregator) Start() {
	for i := 0; i < agt.option.Workers; i++ {
		index := i
		go agt.work(index)
	}
}

// Stop the aggregator
func (agt *Aggregator) Stop() {
	close(agt.quit)
	agt.wg.Wait()
}

// Stop the aggregator safely, the difference with Stop is it guarantees no item is missed during stop
func (agt *Aggregator) SafeStop() {
	if len(agt.eventQueue) == 0 {
		close(agt.quit)
	} else {
		ticker := time.NewTicker(50 * time.Millisecond)
		for range ticker.C {
			if len(agt.eventQueue) == 0 {
				close(agt.quit)
				break
			}
		}
		ticker.Stop()
	}
	agt.wg.Wait()
}

func (agt *Aggregator) work(index int) {
	defer func() {
		if r := recover(); r != nil {
			agt.option.Logger.LogError(agt.ctx, gopkg.LogCtx, map[string]any{
				"err": r,
			}, "[Aggregator] recover worker as bad thing happens")
			agt.work(index)
		}
	}()

	agt.wg.Add(1)
	defer agt.wg.Done()

	batch := make([]any, 0, agt.option.BatchSize)
	lingerTimer := time.NewTimer(0)
	if !lingerTimer.Stop() {
		<-lingerTimer.C
	}
	defer lingerTimer.Stop()

loop:
	for {
		select {
		case req := <-agt.eventQueue:
			batch = append(batch, req)

			batchSize := len(batch)
			if batchSize < agt.option.BatchSize {
				if batchSize == 1 {
					lingerTimer.Reset(agt.option.LingerTime)
				}
				break // 跳出select, 继续从队列中取数据，直到达到批处理数量
			}

			agt.batchProcess(batch)

			if !lingerTimer.Stop() {
				<-lingerTimer.C
			}
			batch = make([]any, 0, agt.option.BatchSize)
		case <-lingerTimer.C: // 如果时间到了，也要进行批处理
			if len(batch) == 0 {
				break
			}

			agt.batchProcess(batch)
			batch = make([]any, 0, agt.option.BatchSize)
		case <-agt.quit:
			if len(batch) != 0 {
				agt.batchProcess(batch)
			}

			break loop
		}
	}
}

func (agt *Aggregator) batchProcess(items []any) {
	agt.wg.Add(1)
	defer agt.wg.Done()
	if err := agt.batchProcessor(agt.ctx, items); err != nil {
		agt.option.Logger.LogError(agt.ctx, gopkg.LogCtx, map[string]any{
			"err": err,
		}, "[Aggregator] error happens")
		if agt.option.ErrorHandler != nil {
			go agt.option.ErrorHandler(agt.ctx, err, items, agt.batchProcessor, agt)
		} else {
			agt.option.Logger.LogError(agt.ctx, gopkg.LogCtx, map[string]any{
				"err": err,
			}, "[Aggregator] error happens in batchProcess and is skipped")
		}
	} else {
		agt.option.Logger.LogInfo(agt.ctx, gopkg.LogCtx, map[string]any{
			"itemLen": len(items),
		}, "[Aggregator] items have been sent.")
	}
}
