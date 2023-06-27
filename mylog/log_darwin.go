package mylog

/**
可以通过向进程发送自定义信号动态调整日志级别
kill -s USR1|USR2 程序的进程号
*/

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/youchuangcd/gopkg"
	"os"
	"syscall"
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

// 每收到一个SIGUSR1信号日志级别减一
// 每收到一个SIGUSR2信号日志级别加一
func watchAndUpdateLoglevel(c chan os.Signal, logger *logrus.Logger) {
	for {
		select {
		case sig := <-c:
			if sig == syscall.SIGUSR1 {
				level := logger.Level
				if level == logrus.PanicLevel {
					fmt.Println("Raise log level: It has been already the most top log level: panic level")
				} else {
					logger.SetLevel(level - 1)
					fmt.Println("Raise log level: the current level is", logger.Level)
				}

			} else if sig == syscall.SIGUSR2 {
				level := logger.Level
				if level == logrus.DebugLevel {
					fmt.Println("Reduce log level: It has been already the lowest log level: debug level")
				} else {
					logger.SetLevel(level + 1)
					fmt.Println("Reduce log level: the current level is", logger.Level)
				}

			} else {
				fmt.Println("receive unknown signal:", sig)
			}
		}
	}
}
