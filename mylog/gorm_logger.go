package mylog

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/youchuangcd/gopkg"
	gLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

// ErrRecordNotFound record not found error
var ErrRecordNotFound = gLogger.ErrRecordNotFound

const (
	gormCallerSkip = 4
)

// Writer log writer interface
type Writer interface {
	WithField(key string, value interface{}) *logrus.Entry
	WithFields(fields Fields) *logrus.Entry
	WithError(err error) *logrus.Entry
	WithContext(ctx context.Context) *logrus.Entry
	WithTime(t time.Time) *logrus.Entry

	Printf(string, ...interface{})
}

// NewGormLogger initialize logger
func NewGormLogger(writer Writer, config gLogger.Config) gLogger.Interface {
	var (
		infoStr      = "%s\n[info] "
		warnStr      = "%s\n[warn] "
		errStr       = "%s\n[error] "
		traceStr     = "%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
	)

	return &gormLogger{
		Writer:       writer,
		Config:       config,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

type gormLogger struct {
	Writer
	gLogger.Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

// LogMode log mode
func (l *gormLogger) LogMode(level gLogger.LogLevel) gLogger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info print info
func (l gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gLogger.Info {
		entry := loggerEntry(ctx, gopkg.LogDB, gormCallerSkip)
		entry.Printf(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Warn print warn messages
func (l gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gLogger.Warn {
		entry := loggerEntry(ctx, gopkg.LogDB, gormCallerSkip)
		entry.Warnf(l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Error print error messages
func (l gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gLogger.Error {
		entry := loggerEntry(ctx, gopkg.LogDB, gormCallerSkip)
		entry.Errorf(l.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Trace print sql message
func (l gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gLogger.Silent {
		return
	}
	entry := loggerEntry(ctx, gopkg.LogDB, gormCallerSkip)

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= gLogger.Error && (!errors.Is(err, ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		latencyTime := float64(elapsed.Nanoseconds()) / 1e6
		entry = entry.WithError(err).WithField("latency_time", latencyTime)
		if rows == -1 {
			entry.Errorf(l.traceErrStr, utils.FileWithLineNum(), err, latencyTime, "-", sql)
		} else {
			entry.Errorf(l.traceErrStr, utils.FileWithLineNum(), err, latencyTime, rows, sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gLogger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		latencyTime := float64(elapsed.Nanoseconds()) / 1e6
		entry = entry.WithField("latency_time", latencyTime)
		if rows == -1 {
			entry.Warnf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, latencyTime, "-", sql)
		} else {
			entry.Warnf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, latencyTime, rows, sql)
		}
	case l.LogLevel == gLogger.Info:
		sql, rows := fc()
		latencyTime := float64(elapsed.Nanoseconds()) / 1e6
		entry = entry.WithField("latency_time", latencyTime)
		if rows == -1 {
			entry.Printf(l.traceStr, utils.FileWithLineNum(), latencyTime, "-", sql)
		} else {
			entry.Printf(l.traceStr, utils.FileWithLineNum(), latencyTime, rows, sql)
		}
	}
}
