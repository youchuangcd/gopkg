package utils

import (
	"github.com/youchuangcd/gopkg/mylog"
)

var PanicHandler = mylog.RecordGoroutineRecoverLog

// NewPoolWithFunc
// @Description: 创建固定方法goroutine池，用完要记得defer goPool.Release()
// @param poolSize
// @param pf
// @param options
// @return goPool
// @return err
func NewPoolWithFunc(poolSize int, pf func(interface{}), options ...ants.Option) (goPool *ants.PoolWithFunc, err error) {
	if len(options) == 0 {
		options = append(options, ants.WithPanicHandler(PanicHandler))
	}
	return ants.NewPoolWithFunc(poolSize, pf, options...)
}

// NewPool
//
//	@Description: 创建goroutine池，用完要记得defer goPool.Release()
//	@param poolSize
//	@param options
//	@return goPool
//	@return err
func NewPool(poolSize int, options ...ants.Option) (goPool *ants.Pool, err error) {
	if len(options) == 0 {
		options = append(options, ants.WithPanicHandler(PanicHandler))
	}
	return ants.NewPool(poolSize, options...)
}

func WithRecover(fn func()) {
	defer func() {
		handler := PanicHandler
		if handler != nil {
			if err := recover(); err != nil {
				handler(err)
			}
		}
	}()

	fn()
}
