package mylog

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/youchuangcd/gopkg"
	"os"
	"strings"
	"time"
)

// WithField allocates a new entry and adds a field to it.
// Debug, Print, Info, Warn, Error, Fatal or Panic must be then applied to
// this new returned entry.
// If you want multiple fields, use `WithFields`.
func (logger *MyLoggerSpecial) WithField(key string, value interface{}) *logrus.Entry {
	return logger.l.WithField(key, value)
}

// Adds a struct of fields to the log entry. All it does is call `WithField` for
// each `Field`.
func (logger *MyLoggerSpecial) WithFields(fields Fields) *logrus.Entry {
	return logger.l.WithFields(logrus.Fields(fields))
}

// Add an error as single field to the log entry.  All it does is call
// `WithError` for the given `error`.
func (logger *MyLoggerSpecial) WithError(err error) *logrus.Entry {
	return logger.l.WithError(err)
}

// Add a context to the log entry.
func (logger *MyLoggerSpecial) WithContext(ctx context.Context) *logrus.Entry {
	return logger.l.WithContext(ctx)
}

// Overrides the time of the log entry.
func (logger *MyLoggerSpecial) WithTime(t time.Time) *logrus.Entry {
	return logger.l.WithTime(t)
}

//func (logger *MyLoggerSpecial) Debug(args ...interface{}) {
//	logger.l.Debug(args...)
//}
//
//func (logger *MyLoggerSpecial) Info(args ...interface{}) {
//	logger.l.Info(args...)
//}
//
//func (logger *MyLoggerSpecial) Warn(args ...interface{}) {
//	logger.l.Warn(args...)
//}
//
//func (logger *MyLoggerSpecial) Warning(args ...interface{}) {
//	logger.l.Warn(args...)
//}
//
//func (logger *MyLoggerSpecial) Error(args ...interface{}) {
//	logger.l.Error(args...)
//}
//
//func (logger *MyLoggerSpecial) Fatal(args ...interface{}) {
//	logger.l.Fatal(args...)
//}

func (logger *MyLoggerSpecial) Panic(args ...interface{}) {
	logger.l.Panic(args...)
}

func (logger *MyLoggerSpecial) Tracef(format string, args ...interface{}) {
	logger.l.Tracef(format, args...)
}

func (logger *MyLoggerSpecial) Debugf(format string, args ...interface{}) {
	logger.l.Debugf(format, args...)
}

func (logger *MyLoggerSpecial) Infof(format string, args ...interface{}) {
	logger.l.Infof(format, args...)
}

func (logger *MyLoggerSpecial) Printf(format string, args ...interface{}) {
	logger.l.Printf(format, args...)
}

func (logger *MyLoggerSpecial) Warnf(format string, args ...interface{}) {
	logger.l.Warnf(format, args...)
}

func (logger *MyLoggerSpecial) Warningf(format string, args ...interface{}) {
	logger.l.Warningf(format, args...)
}

func (logger *MyLoggerSpecial) Errorf(format string, args ...interface{}) {
	logger.l.Errorf(format, args...)
}

func (logger *MyLoggerSpecial) Fatalf(format string, args ...interface{}) {
	logger.l.Fatalf(format, args...)
}

func (logger *MyLoggerSpecial) Panicf(format string, args ...interface{}) {
	logger.l.Panicf(format, args...)
}

func (logger *MyLoggerSpecial) Debugln(args ...interface{}) {
	logger.l.Debugln(args...)
}

func (logger *MyLoggerSpecial) Infoln(args ...interface{}) {
	logger.l.Infoln(args...)
}

func (logger *MyLoggerSpecial) Println(args ...interface{}) {
	logger.l.Println(args...)
}

func (logger *MyLoggerSpecial) Warnln(args ...interface{}) {
	logger.l.Warnln(args...)
}

func (logger *MyLoggerSpecial) Warningln(args ...interface{}) {
	logger.l.Warnln(args...)
}

func (logger *MyLoggerSpecial) Errorln(args ...interface{}) {
	logger.l.Errorln(args...)
}

func (logger *MyLoggerSpecial) Fatalln(args ...interface{}) {
	logger.l.Fatalln(args...)
}

func (logger *MyLoggerSpecial) Panicln(args ...interface{}) {
	logger.l.Panicln(args...)
}

func (logger *MyLoggerSpecial) Exit(code int) {
	logger.l.Exit(code)
}

// AddHook adds a hook to the logger hooks.
func (logger *MyLoggerSpecial) AddHook(hook Hook) {
	logger.l.AddHook(logrus.Hook(hook))
}

func (logger *MyLoggerSpecial) Level(level string) {
	switch strings.ToLower(level) {
	case "debug":
		logger.l.SetLevel(logrus.DebugLevel)
	case "warn":
		logger.l.SetLevel(logrus.WarnLevel)
	case "error":
		logger.l.SetLevel(logrus.ErrorLevel)
	case "panic":
		logger.l.SetLevel(logrus.PanicLevel)
	case "fatal":
		logger.l.SetLevel(logrus.FatalLevel)
	case "trace":
		logger.l.SetLevel(logrus.TraceLevel)
	default:
		logger.l.SetLevel(logrus.InfoLevel)
	}
}

func (logger *MyLoggerSpecial) OutputPath(path string) (err error) {
	var file *os.File
	file, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	logger.l.Out = file
	return
}

func (logger *MyLoggerSpecial) Debug(msg string, fields map[string]interface{}) {
	if msg == "" && len(fields) == 0 {
		return
	}
	logger.l.WithFields(fieldsAttachCategory(fields)).Debug(msg)
}

func (logger *MyLoggerSpecial) Info(msg string, fields map[string]interface{}) {
	if msg == "" && len(fields) == 0 {
		return
	}
	logger.l.WithFields(fieldsAttachCategory(fields)).Info(msg)
}

func (logger *MyLoggerSpecial) Warn(msg string, fields map[string]interface{}) {
	if msg == "" && len(fields) == 0 {
		return
	}
	logger.l.WithFields(fieldsAttachCategory(fields)).Warn(msg)
}

func (logger *MyLoggerSpecial) Warning(msg string, fields map[string]interface{}) {
	if msg == "" && len(fields) == 0 {
		return
	}
	logger.l.WithFields(fieldsAttachCategory(fields)).Warn(msg)
}

func (logger *MyLoggerSpecial) Error(msg string, fields map[string]interface{}) {
	if msg == "" && len(fields) == 0 {
		return
	}
	logger.l.WithFields(fieldsAttachCategory(fields)).Error(msg)
}

func (logger *MyLoggerSpecial) Fatal(msg string, fields map[string]interface{}) {
	if msg == "" && len(fields) == 0 {
		return
	}
	logger.l.WithFields(fieldsAttachCategory(fields)).Fatal(msg)
}

func fieldsAttachCategory(fields map[string]interface{}) map[string]interface{} {
	if fields == nil {
		fields = map[string]interface{}{
			"category": gopkg.LogRocketMQTCP,
		}
	} else if _, ok := fields["category"]; !ok {
		// 目前就rocketmq-client-go包用这种方式写日志，所以默认记录为rocketmq-tcp
		fields["category"] = gopkg.LogRocketMQTCP
	}
	return fields
}
