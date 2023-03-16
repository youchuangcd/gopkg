package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/youchuangcd/gopkg"
	"github.com/youchuangcd/gopkg/common/utils"
	"github.com/youchuangcd/gopkg/mylog"
	"io/ioutil"
	"runtime"
	"time"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

func Recover(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			rs, _ := c.Value(gopkg.ContextRequestSourceKey).(string)
			requestParam := utils.GetParams(c)
			responseErr := gopkg.ErrorInternalServer
			if e, ok := err.(*gopkg.Error); ok {
				responseErr = e
			}
			utils.ToJson(c, responseErr, nil)
			var latencyTime int64
			if requestStartTime, ok := c.Value(gopkg.ContextRequestStartTimeKey).(time.Time); ok {
				endTime := time.Now()
				// 执行时间
				latencyTime = endTime.Sub(requestStartTime).Milliseconds()
			}
			// 请求IP
			clientIP := c.ClientIP()
			mylog.WithError(c, gopkg.LogPanic, map[string]interface{}{
				"code":           c.Writer.Status(),
				"request_method": c.Request.Method,
				"request_domain": c.Request.Host,
				"request_uri":    c.Request.URL.Path,
				"request_url":    c.Request.Host + c.Request.RequestURI,
				"request_param":  requestParam,
				"request_header": c.Request.Header,
				"request_source": rs,
				"response_code":  responseErr.GetCode(),
				"response_msg":   responseErr.GetMsg(),
				"err":            err,
				"stack":          string(stack(3)),
				"log_type":       gopkg.RequestLogTypeFlag,
				"latency_time":   latencyTime,
				"client_ip":      clientIP,
			}, "捕获到panic错误")
		}
	}()
	c.Next()
}

// stack returns a nicely formatted stack frame, skipping skip frames.
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}
