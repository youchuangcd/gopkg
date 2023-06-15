package common

import (
	"context"
	"errors"
	"github.com/youchuangcd/gopkg/mylog"
	"sync"
	"testing"
	"time"
)

func TestAggregator_Basic(t *testing.T) {
	batchProcess := func(ctx context.Context, items []interface{}) error {
		t.Logf("handler %d items", len(items))
		return nil
	}

	aggregator := NewAggregator(context.Background(), batchProcess, func(option AggregatorOption) AggregatorOption {
		option.Logger = mylog.GetLogger()
		return option
	})

	aggregator.Start()

	for i := 0; i < 1000; i++ {
		aggregator.TryEnqueue(i)
	}

	aggregator.SafeStop()
}

func TestAggregator_Complex(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(100)

	batchProcess := func(ctx context.Context, items []interface{}) error {
		defer wg.Add(-len(items))
		time.Sleep(20 * time.Millisecond)
		if len(items) != 4 {
			return errors.New("len(items) != 4")
		}
		return nil
	}

	errorHandler := func(ctx context.Context, err error, items []interface{}, batchProcessFunc BatchProcessFunc, aggregator *Aggregator) {
		if err == nil {
			t.FailNow()
		}
		t.Logf("Receive error, item size is %d", len(items))
	}

	aggregator := NewAggregator(context.Background(), batchProcess, func(option AggregatorOption) AggregatorOption {
		option.BatchSize = 4
		option.Workers = 2
		option.LingerTime = 8 * time.Millisecond
		option.ErrorHandler = errorHandler
		option.Logger = mylog.GetLogger()
		return option
	})

	aggregator.Start()

	for i := 0; i < 100; i++ {
		for !aggregator.TryEnqueue(i) {
			time.Sleep(10 * time.Millisecond)
		}
	}

	aggregator.SafeStop()
	wg.Wait()
}

func TestAggregator_LingerTimeOut(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(100)

	batchProcess := func(ctx context.Context, items []interface{}) error {
		defer wg.Add(-len(items))
		if len(items) != 4 {
			t.Log("linger time out")
		}
		return nil
	}

	aggregator := NewAggregator(context.Background(), batchProcess, func(option AggregatorOption) AggregatorOption {
		option.BatchSize = 4
		option.Workers = 1
		option.LingerTime = 150 * time.Millisecond
		option.Logger = mylog.GetLogger()
		return option
	})

	aggregator.Start()

	for i := 0; i < 100; i++ {
		aggregator.TryEnqueue(i)
		time.Sleep(100 * time.Millisecond)
	}

	aggregator.SafeStop()
	wg.Wait()
}
