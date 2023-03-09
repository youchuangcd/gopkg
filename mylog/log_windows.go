package mylog

import (
	"github.com/sirupsen/logrus"
	"github.com/youchuangcd/gopkg"
	"os"
)

func InitLog(logConf LogConfig) {
	conf = logConf
	initLog.Do(func() {
		tmpLogger := logrus.New()
		//logFilePath := getLogFilePath()
		//if err := os.MkdirAll(logFilePath, 0777); err != nil {
		//	log.Fatalf("cannot makedir logFilePath: %s; stack: %v", logFilePath, errors.WithStack(err))
		//}
		//logFileName := path.Join(logFilePath, getLogFileName())
		//
		//writer, err := rotatelogs.New(
		//	logFileName,
		//	//rotatelogs.WithLinkName(getWithLinkLatestLogFileName()),                                    // 生成软链，指向最新日志文件
		//	rotatelogs.WithMaxAge(time.Duration(conf.SaveMaxAge)*time.Duration(24)*time.Hour), // 文件最大保存时间
		//	rotatelogs.WithRotationTime(time.Duration(conf.RotationTime)*time.Hour),           // 日志切割时间间隔
		//)
		//if err != nil {
		//	log.Fatalf("config local file system logger error. %v", errors.WithStack(err))
		//}

		//switch level := conf.Level; level {
		///*
		//   如果日志级别不是debug就不要打印日志到控制台了
		//*/
		//case "debug":
		//	tmpLogger.SetLevel(logrus.DebugLevel)
		//	tmpLogger.SetOutput(os.Stderr)
		//case "info":
		//	//setNull(tmpLogger)
		//	tmpLogger.SetLevel(logrus.InfoLevel)
		//case "warn":
		//	//setNull(tmpLogger)
		//	tmpLogger.SetLevel(logrus.WarnLevel)
		//case "error":
		//	//setNull(tmpLogger)
		//	tmpLogger.SetLevel(logrus.ErrorLevel)
		//default:
		//	//setNull(tmpLogger)
		//	tmpLogger.SetLevel(logrus.InfoLevel)
		//}
		writer := os.Stdout
		tmpLogger.SetOutput(writer)
		tmpLogger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat:   gopkg.DateMsTimeFormat,
			DisableHTMLEscape: true,
		})
		// 关闭安全并发写锁，如果文件是append模式打开的话，就不需要锁
		tmpLogger.SetNoLock()

		//lfHook := lfshook.NewHook(lfshook.WriterMap{
		//	logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
		//	logrus.InfoLevel:  writer,
		//	logrus.WarnLevel:  writer,
		//	logrus.ErrorLevel: writer,
		//	logrus.FatalLevel: writer,
		//	logrus.PanicLevel: writer,
		//}, tmpLogger.Formatter)
		//tmpLogger.AddHook(lfHook)

		// 把初始化的对象赋值给全局变量
		logger = &MyLogger{l: tmpLogger}
		loggerSpecial = &MyLoggerSpecial{l: tmpLogger}
	})
}
