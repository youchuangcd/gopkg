package mylog

/**
Linux和Mac可以通过向进程发送自定义信号动态调整日志级别
kill -s USR1|USR2 程序的进程号
*/

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/youchuangcd/gopkg"
	"io/ioutil"
	"log"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Fields logrus.Fields

type Hook logrus.Hook

var (
	// 初始化在各自平台文件里的InitLog
	logger        *MyLogger
	loggerSpecial *MyLoggerSpecial
	initLog       sync.Once
	conf          LogConfig
)

type LogConfig struct {
	SavePath     string
	SaveName     string
	TimeFormat   string
	RotationTime int
	FileExt      string
	SaveMaxAge   int
	Level        string
	Env          string
}

// Logger
// @Description: 实现了Interface的结构体
type MyLogger struct {
	l *logrus.Logger
}

type MyLoggerSpecial struct {
	l *logrus.Logger
}

type LoggerInterface interface {
	WithField(key string, value interface{}) *logrus.Entry
	WithFields(fields Fields) *logrus.Entry
	WithError(err error) *logrus.Entry
	WithContext(ctx context.Context) *logrus.Entry
	WithTime(t time.Time) *logrus.Entry

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})

	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Println(args ...interface{})
	Warnln(args ...interface{})
	Warningln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})

	Exit(code int)

	AddHook(hook Hook)
	Level(level string)
	OutputPath(path string) (err error)
	SetPrefix(s string)

	// 自定义扩展方法
	LogDebug(ctx context.Context, logCategory string, logContent map[string]interface{}, msg string)
	LogInfo(ctx context.Context, logCategory string, logContent map[string]interface{}, msg string)
	LogWarn(ctx context.Context, logCategory string, logContent map[string]interface{}, msg string)
	LogError(ctx context.Context, logCategory string, logContent map[string]interface{}, msg string)
}

type LoggerSpecialInterface interface {
	WithField(key string, value interface{}) *logrus.Entry
	WithFields(fields Fields) *logrus.Entry
	WithError(err error) *logrus.Entry
	WithContext(ctx context.Context) *logrus.Entry
	WithTime(t time.Time) *logrus.Entry

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	//Debug(args ...interface{})
	//Info(args ...interface{})
	//Warn(args ...interface{})
	//Warning(args ...interface{})
	//Error(args ...interface{})
	//Fatal(args ...interface{})
	//Panic(args ...interface{})
	Debug(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
	Warn(msg string, fields map[string]interface{})
	Warning(msg string, fields map[string]interface{})
	Error(msg string, fields map[string]interface{})
	Fatal(msg string, fields map[string]interface{})

	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Println(args ...interface{})
	Warnln(args ...interface{})
	Warningln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})

	Exit(code int)

	AddHook(hook Hook)
	Level(level string)
	OutputPath(path string) (err error)
}

type ExtError map[any]interface{}

func (p ExtError) String() string {
	b, _ := json.Marshal(p)
	return string(b)
}

func setNull(logger *logrus.Logger) {
	//src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	//if err != nil {
	//	fmt.Println("err", err)
	//}
	//writer := bufio.NewWriter(src)
	//logger.SetOutput(src)
	logger.SetOutput(ioutil.Discard)
}

// WithDebug debug日志
func WithDebug(c context.Context, logCategory string, logContent map[string]interface{}, msg string) {
	loggerEntry(c, logCategory).WithFields(logContent).Debug(msg)
}

// WithInfo info日志
func WithInfo(c context.Context, logCategory string, logContent map[string]interface{}, msg string) {
	loggerEntry(c, logCategory).WithFields(logContent).Info(msg)
}

// WithWarn warning日志
func WithWarn(c context.Context, logCategory string, logContent map[string]interface{}, msg string) {
	loggerEntry(c, logCategory).WithFields(logContent).Warn(msg)
}

// WithError error日志
func WithError(c context.Context, logCategory string, logContent map[string]interface{}, msg string) {
	loggerEntry(c, logCategory).WithFields(logContent).Error(msg)
}

func Debug(c context.Context, logCategory string, msg string) {
	loggerEntry(c, logCategory).Debug(msg)
}

func Info(c context.Context, logCategory string, msg string) {
	loggerEntry(c, logCategory).Info(msg)
}

func Warn(c context.Context, logCategory string, msg string) {
	loggerEntry(c, logCategory).Warn(msg)
}

func Error(c context.Context, logCategory string, msg string) {
	loggerEntry(c, logCategory).Error(msg)
}

// GetLogger
// @Description: 获取logger对象
// @return *LoggerInterface
func GetLogger() LoggerInterface {
	return logger
}

func GetSpecialLogger() LoggerSpecialInterface {
	return loggerSpecial
}

// LoggerEntry returns a logrus.Entry with as much context as possible
func loggerEntry(ctx context.Context, category string, args ...interface{}) *logrus.Entry {
	callerSkip := 3
	if len(args) > 0 {
		if v, ok := args[0].(int); ok {
			callerSkip = v
		}
	}
	if category == "" {
		if v, ok := ctx.Value(gopkg.ContextLogCategoryKey).(string); ok {
			category = v
		} else {
			category = "unknown"
		}
	}
	entry := logger.WithField("category", category)
	if ctx != nil {
		if traceId, ok := ctx.Value(gopkg.RequestHeaderTraceIdKey).(string); ok {
			entry = entry.WithField(gopkg.LogTraceIdKey, traceId)
		}
		if spanId, ok := ctx.Value(gopkg.RequestHeaderSpanIdKey).(string); ok {
			entry = entry.WithField(gopkg.LogParentSpanIdKey, spanId)
		}
		if spanId, ok := ctx.Value(gopkg.LogSpanIdKey).(string); ok {
			entry = entry.WithField(gopkg.LogSpanIdKey, spanId)
		}
		if userId, ok := ctx.Value(gopkg.ContextRequestSysUserIdKey).(uint); ok {
			entry = entry.WithField(gopkg.LogUserIdKey, userId)
		} else if userId, ok = ctx.Value(gopkg.ContextRequestUserIdKey).(uint); ok {
			entry = entry.WithField(gopkg.LogUserIdKey, userId)
		}
		if msgId, ok := ctx.Value(gopkg.LogMsgIdKey).(string); ok {
			entry = entry.WithField(gopkg.LogMsgIdKey, msgId)
		}
		if parentMsgId, ok := ctx.Value(gopkg.LogParentMsgIdKey).(string); ok {
			entry = entry.WithField(gopkg.LogParentMsgIdKey, parentMsgId)
		}
		if taskId, ok := ctx.Value(gopkg.LogTaskIdKey).(string); ok {
			entry = entry.WithField(gopkg.LogTaskIdKey, taskId)
		}
	}
	entry = entry.WithField("env", conf.Env)
	// 追加调用方法名称
	pc, _, _, _ := runtime.Caller(callerSkip)
	entry = entry.WithField("func", runtime.FuncForPC(pc).Name())
	return entry
}

// getLogFilePath get the log file save path
func getLogFilePath() string {
	//return fmt.Sprintf("%s%s", setting.AppSetting.RuntimeRootPath, setting.AppSetting.LogSavePath)
	return fmt.Sprintf("%s/", strings.TrimRight(conf.SavePath, "/"))
}

// getLogFileName get the save name of the log file
func getLogFileName() string {
	return fmt.Sprintf("%s%s.%s",
		conf.SaveName,
		conf.TimeFormat,
		conf.FileExt,
	)
}

// RecordGoroutineRecoverLog
// @Description: 记录goroutine recover日志
// @param r
func RecordGoroutineRecoverLog(r interface{}) {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	WithError(nil, gopkg.LogPanic, map[string]interface{}{
		"err":   r,
		"stack": string(buf[:n]),
	}, "捕获到一个goroutine panic错误")
}

func RecordRecoverLog(r interface{}) {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	withField := map[string]interface{}{
		"err":   r,
		"stack": string(buf[:n]),
	}
	if logger != nil {
		WithError(nil, gopkg.LogPanic, withField, "捕获到一个panic错误")
	} else {
		log.Printf("捕获到一个panic错误，且logger未初始化完成：%+v", withField)
	}
}
