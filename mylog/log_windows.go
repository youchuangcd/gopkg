package mylog

import (
	"github.com/sirupsen/logrus"
	"github.com/youchuangcd/gopkg"
	"os"
)

func InitLog() {
	initLog.Do(func() {
		tmpLogger := logrus.New()
		writer := os.Stdout
		tmpLogger.SetOutput(writer)
		tmpLogger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat:   gopkg.DateMsTimeFormat,
			DisableHTMLEscape: true,
		})
		// 关闭安全并发写锁，如果文件是append模式打开的话，就不需要锁
		tmpLogger.SetNoLock()
		// 把初始化的对象赋值给全局变量
		logger = &MyLogger{l: tmpLogger}
		loggerSpecial = &MyLoggerSpecial{l: tmpLogger}
	})
}
