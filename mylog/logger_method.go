package mylog

import (
	"context"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

// WithField allocates a new entry and adds a field to it.
// Debug, Print, Info, Warn, Error, Fatal or Panic must be then applied to
// this new returned entry.
// If you want multiple fields, use `WithFields`.
func (logger *MyLogger) WithField(key string, value interface{}) *logrus.Entry {
	return logger.l.WithField(key, value)
}

// Adds a struct of fields to the log entry. All it does is call `WithField` for
// each `Field`.
func (logger *MyLogger) WithFields(fields Fields) *logrus.Entry {
	return logger.l.WithFields(logrus.Fields(fields))
}

// Add an error as single field to the log entry.  All it does is call
// `WithError` for the given `error`.
func (logger *MyLogger) WithError(err error) *logrus.Entry {
	return logger.l.WithError(err)
}

// Add a context to the log entry.
func (logger *MyLogger) WithContext(ctx context.Context) *logrus.Entry {
	return logger.l.WithContext(ctx)
}

// Overrides the time of the log entry.
func (logger *MyLogger) WithTime(t time.Time) *logrus.Entry {
	return logger.l.WithTime(t)
}

func (logger *MyLogger) Debug(args ...interface{}) {
	logger.l.Debug(args...)
}

func (logger *MyLogger) Info(args ...interface{}) {
	logger.l.Info(args...)
}

func (logger *MyLogger) Warn(args ...interface{}) {
	logger.l.Warn(args...)
}

func (logger *MyLogger) Warning(args ...interface{}) {
	logger.l.Warn(args...)
}

func (logger *MyLogger) Error(args ...interface{}) {
	logger.l.Error(args...)
}

func (logger *MyLogger) Fatal(args ...interface{}) {
	logger.l.Fatal(args...)
}

func (logger *MyLogger) Panic(args ...interface{}) {
	logger.l.Panic(args...)
}

func (logger *MyLogger) Tracef(format string, args ...interface{}) {
	logger.l.Tracef(format, args...)
}

func (logger *MyLogger) Debugf(format string, args ...interface{}) {
	logger.l.Debugf(format, args...)
}

func (logger *MyLogger) Infof(format string, args ...interface{}) {
	logger.l.Infof(format, args...)
}

func (logger *MyLogger) Printf(format string, args ...interface{}) {
	logger.l.Printf(format, args...)
}

func (logger *MyLogger) Warnf(format string, args ...interface{}) {
	logger.l.Warnf(format, args...)
}

func (logger *MyLogger) Warningf(format string, args ...interface{}) {
	logger.l.Warningf(format, args...)
}

func (logger *MyLogger) Errorf(format string, args ...interface{}) {
	logger.l.Errorf(format, args...)
}

func (logger *MyLogger) Fatalf(format string, args ...interface{}) {
	logger.l.Fatalf(format, args...)
}

func (logger *MyLogger) Panicf(format string, args ...interface{}) {
	logger.l.Panicf(format, args...)
}

func (logger *MyLogger) Debugln(args ...interface{}) {
	logger.l.Debugln(args...)
}

func (logger *MyLogger) Infoln(args ...interface{}) {
	logger.l.Infoln(args...)
}

func (logger *MyLogger) Println(args ...interface{}) {
	logger.l.Println(args...)
}

func (logger *MyLogger) Warnln(args ...interface{}) {
	logger.l.Warnln(args...)
}

func (logger *MyLogger) Warningln(args ...interface{}) {
	logger.l.Warnln(args...)
}

func (logger *MyLogger) Errorln(args ...interface{}) {
	logger.l.Errorln(args...)
}

func (logger *MyLogger) Fatalln(args ...interface{}) {
	logger.l.Fatalln(args...)
}

func (logger *MyLogger) Panicln(args ...interface{}) {
	logger.l.Panicln(args...)
}

func (logger *MyLogger) Exit(code int) {
	logger.l.Exit(code)
}

// AddHook adds a hook to the logger hooks.
func (logger *MyLogger) AddHook(hook Hook) {
	logger.l.AddHook(logrus.Hook(hook))
}

func (logger *MyLogger) Level(level string) {
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

func (logger *MyLogger) OutputPath(path string) (err error) {
	var file *os.File
	file, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	logger.l.Out = file
	return
}

func (logger *MyLogger) SetPrefix(s string) {

}

// LogDebug debug日志
func (logger *MyLogger) LogDebug(ctx context.Context, logCategory string, logContent map[string]interface{}, msg string) {
	loggerEntry(ctx, logCategory).WithFields(logContent).Debug(msg)
}

// LogInfo info日志
func (logger *MyLogger) LogInfo(ctx context.Context, logCategory string, logContent map[string]interface{}, msg string) {
	loggerEntry(ctx, logCategory).WithFields(logContent).Info(msg)
}

// LogWarn warning日志
func (logger *MyLogger) LogWarn(ctx context.Context, logCategory string, logContent map[string]interface{}, msg string) {
	loggerEntry(ctx, logCategory).WithFields(logContent).Warn(msg)
}

// LogError error日志
func (logger *MyLogger) LogError(ctx context.Context, logCategory string, logContent map[string]interface{}, msg string) {
	loggerEntry(ctx, logCategory).WithFields(logContent).Error(msg)
}
