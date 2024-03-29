package client

import (
	"bytes"
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/youchuangcd/gopkg"
	"github.com/youchuangcd/gopkg/common/utils"
	"github.com/youchuangcd/gopkg/mylog"
	"io"
	"net/http"
	"net/http/httptrace"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

var (
	// UserAgent UA
	UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"
	// DefaultClient 默认Client
	DefaultClient = Client{&http.Client{Transport: http.DefaultTransport}}
	// DebugMode 用来打印调试信息
	DebugMode = gopkg.HttpClientDebugMode
	// DeepDebugInfo 调试信息
	DeepDebugInfo      = gopkg.HttpClientDeepDebugInfo
	insReqStartTimeKey = reqStartTimeKey{}
	insReqUrlKey       = reqUrlKey{}
	json               = jsoniter.ConfigCompatibleWithStandardLibrary
)

type reqStartTimeKey struct{}
type reqUrlKey struct{}

// --------------------------------------------------------------------

// Client 负责发送HTTP请求
type Client struct {
	*http.Client
}

// WithTraceId 把traceId加入context中
func WithTraceId(ctx context.Context, traceId string) context.Context {
	return context.WithValue(ctx, gopkg.RequestHeaderTraceIdKey, traceId)
}

// ContextGetTraceId 从context中获取traceId
func ContextGetTraceId(ctx context.Context) (traceId string, ok bool) {
	traceId, ok = ctx.Value(gopkg.RequestHeaderTraceIdKey).(string)
	return
}

func newRequest(ctx context.Context, method, reqUrl string, headers http.Header, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(method, reqUrl, body)
	if err != nil {
		return
	}

	if headers == nil {
		headers = http.Header{}
	}

	req.Header = headers
	req = req.WithContext(ctx)

	return
}

// DoRequestWith 请求
func (r Client) DoRequestWith(ctx context.Context, method, reqUrl string, headers http.Header, body io.Reader,
	bodyLength int) (resp *http.Response, err error) {

	req, err := newRequest(ctx, method, reqUrl, headers, body)
	if err != nil {
		return
	}
	req.ContentLength = int64(bodyLength)
	return r.Do(ctx, req)
}

// DoRequestWithForm Form请求
func (r Client) DoRequestWithForm(ctx context.Context, method, reqUrl string, headers http.Header,
	data interface{}) (resp *http.Response, err error) {

	var reqBody []byte
	if data != nil {
		// 直接传url.Values encode后的字节
		if v, ok := data.([]byte); ok {
			reqBody = v
		} else if v2, ok2 := data.(string); ok2 {
			reqBody = []byte(v2)
		} else if v3, ok3 := data.(url.Values); ok3 {
			reqBody = []byte(v3.Encode())
		} else if v4, ok4 := data.(map[string]any); ok4 {
			reqBody = []byte(utils.MapUrlEncode(v4))
		}
	}

	if headers == nil {
		headers = http.Header{}
	}

	if headers.Get("Content-Type") == "" {
		headers.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	if method == http.MethodGet {
		if strings.Index(reqUrl, "?") > 0 {
			reqUrl += "&" + string(reqBody)
		} else {
			reqUrl += "?" + string(reqBody)
		}
		reqBody = nil
	}
	return r.DoRequestWith(ctx, method, reqUrl, headers, bytes.NewReader(reqBody), len(reqBody))
}

// DoRequestWithJson JSON请求
func (r Client) DoRequestWithJson(ctx context.Context, method, reqUrl string, headers http.Header,
	data interface{}) (resp *http.Response, err error) {

	var reqBody []byte
	if data != nil {
		// 直接传json序列化后的字节
		if v, ok := data.([]byte); ok {
			reqBody = v
		} else {
			reqBody, err = json.Marshal(data)
			if err != nil {
				return
			}
		}
	}

	if headers == nil {
		headers = http.Header{}
	}

	if headers.Get("Content-Type") == "" {
		headers.Add("Content-Type", "application/json")
	}
	return r.DoRequestWith(ctx, method, reqUrl, headers, bytes.NewReader(reqBody), len(reqBody))
}

// Do 请求
func (r Client) Do(ctx context.Context, req *http.Request) (resp *http.Response, err error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if traceId, ok := ContextGetTraceId(ctx); ok {
		req.Header.Set(gopkg.RequestHeaderTraceIdKey, traceId)
	}

	if _, ok := req.Header["User-Agent"]; !ok {
		req.Header.Set("User-Agent", UserAgent)
	}
	// 追加istio B3 请求头
	for _, key := range gopkg.RequestB3Headers {
		if val, ok := ctx.Value(key).(string); ok && val != "" {
			req.Header.Set(key, val)
		}
	}
	//var newCtx = ctx
	//// http请求的话，要提取request里面的上下文才可以获取到b3请求头
	//if ginCtx, ok := ctx.Value(gin.ContextKey).(*gin.Context); ok {
	//	newCtx = ginCtx.Request.Context()
	//}
	// 把上游传递过来的追踪参数从上下文中提取出来注入到请求中
	//otel.GetTextMapPropagator().Inject(newCtx, propagation.HeaderCarrier(req.Header))
	if DebugMode {
		trace := &httptrace.ClientTrace{
			//GotConn: func(connInfo httptrace.GotConnInfo) {
			//	remoteAddr := connInfo.Conn.RemoteAddr()
			//	mylog.Info(ctx, LogCategoryHttp, fmt.Sprintf("Network: %s, Remote ip:%s, URL: %s", remoteAddr.Network(), remoteAddr.String(), req.URL))
			//},
		}
		req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
		bs, bErr := httputil.DumpRequest(req, DeepDebugInfo)
		if bErr != nil {
			err = bErr
			return
		}
		mylog.Info(ctx, gopkg.LogHttp, string(bs))
	}

	transport := r.Transport // don't change r.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	// avoid cancel() is called before Do(req), but isn't accurate
	select {
	case <-ctx.Done():
		err = ctx.Err()
		return
	default:
	}

	if tr, ok := getRequestCanceler(transport); ok {
		// support CancelRequest
		reqC := make(chan bool, 1)
		go func() {
			resp, err = r.Client.Do(req)
			reqC <- true
		}()
		select {
		case <-reqC:
		case <-ctx.Done():
			tr.CancelRequest(req)
			<-reqC
			err = ctx.Err()
		}
	} else {
		resp, err = r.Client.Do(req)
	}
	return
}

// ErrorInfo 错误信息
type ErrorInfo struct {
	Err     string `json:"error,omitempty"`
	Key     string `json:"key,omitempty"`
	TraceId string `json:"traceId,omitempty"`
	Errno   int    `json:"errno,omitempty"`
	Code    int    `json:"code"`
}

// Error 错误
func (r *ErrorInfo) Error() string {
	return r.Err
}

func parseError(e *ErrorInfo, r io.Reader) {
	body, err1 := io.ReadAll(r)
	if err1 != nil {
		e.Err = err1.Error()
		return
	}

	var ret struct {
		Err   string `json:"error"`
		Key   string `json:"key"`
		Errno int    `json:"errno"`
	}
	if json.Unmarshal(body, &ret) == nil && ret.Err != "" {
		e.Err, e.Key, e.Errno = ret.Err, ret.Key, ret.Errno
		return
	}
	e.Err = string(body)
}

// ResponseError 错误响应
func ResponseError(resp *http.Response) error {
	e := &ErrorInfo{
		TraceId: resp.Header.Get(gopkg.LogHttp),
		Code:    resp.StatusCode,
	}
	if resp.StatusCode > 299 {
		if resp.ContentLength != 0 {
			ct, ok := resp.Header["Content-Type"]
			if ok && strings.HasPrefix(ct[0], "application/json") {
				parseError(e, resp.Body)
			} else {
				bs, rErr := io.ReadAll(resp.Body)
				if rErr != nil {
					e.Err = rErr.Error()
				}
				e.Err = strings.TrimRight(string(bs), "\n")
			}
		} else if resp.Status != "" {
			e.Err = resp.Status
		}
	}
	return e
}

// CallRet Http请求
func CallRet(ctx context.Context, ret interface{}, resp *http.Response) (err error) {
	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	if DebugMode {
		var (
			latencyTime time.Duration
			reqUrl      string
		)
		if startTime, ok := ctx.Value(insReqStartTimeKey).(time.Time); ok {
			// 执行时间
			latencyTime = time.Since(startTime)
		}
		if url, ok := ctx.Value(insReqUrlKey).(string); ok {
			reqUrl = url
		}
		bs, dErr := httputil.DumpResponse(resp, DeepDebugInfo)
		if dErr != nil {
			err = dErr
			return
		}
		mylog.WithInfo(ctx, gopkg.LogHttp, map[string]interface{}{
			"latency_time_str": latencyTime.String(),
			"latency_time":     float64(latencyTime.Nanoseconds()) / 1e6,
			"req_url":          reqUrl,
		}, string(bs))
	}
	if resp.StatusCode/100 == 2 {
		if ret != nil && resp.ContentLength != 0 {
			err = json.NewDecoder(resp.Body).Decode(ret)
			if err != nil {
				return
			}
		}
		if resp.StatusCode == 200 {
			return nil
		}
	}
	return ResponseError(resp)
}

// CallRetResp Http请求
func CallRetResp(ctx context.Context, retBuf *bytes.Buffer, resp *http.Response) (err error) {
	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()
	if DebugMode {
		var (
			latencyTime time.Duration
			reqUrl      string
		)
		if startTime, ok := ctx.Value(insReqStartTimeKey).(time.Time); ok {
			// 执行时间
			latencyTime = time.Since(startTime)
		}
		if url, ok := ctx.Value(insReqUrlKey).(string); ok {
			reqUrl = url
		}
		bs, dErr := httputil.DumpResponse(resp, DeepDebugInfo)
		if dErr != nil {
			err = dErr
			return
		}
		mylog.WithInfo(ctx, gopkg.LogHttp, map[string]interface{}{
			"latency_time_str": latencyTime.String(),
			"latency_time":     float64(latencyTime.Nanoseconds()) / 1e6,
			"req_url":          reqUrl,
		}, string(bs))
	}
	if resp.StatusCode/100 == 2 {
		if retBuf != nil && resp.ContentLength != 0 {
			_, err = io.Copy(retBuf, resp.Body)
			if err != nil {
				return
			}
		}
		if resp.StatusCode == 200 {
			return nil
		}
	}
	return ResponseError(resp)
}

// CallWithJson JSON请求
func (r Client) CallWithJson(ctx context.Context, ret interface{}, method, reqUrl string, headers http.Header,
	param interface{}) (err error) {
	if DebugMode {
		ctx = context.WithValue(ctx, insReqStartTimeKey, time.Now())
		ctx = context.WithValue(ctx, insReqUrlKey, reqUrl)
	}
	resp, err := r.DoRequestWithJson(ctx, method, reqUrl, headers, param)
	if err != nil {
		return err
	}
	return CallRet(ctx, ret, resp)
}

// CallWithForm Form请求
func (r Client) CallWithForm(ctx context.Context, ret interface{}, method, reqUrl string, headers http.Header,
	param interface{}) (err error) {
	if DebugMode {
		ctx = context.WithValue(ctx, insReqStartTimeKey, time.Now())
		ctx = context.WithValue(ctx, insReqUrlKey, reqUrl)
	}
	resp, err := r.DoRequestWithForm(ctx, method, reqUrl, headers, param)
	if err != nil {
		return err
	}
	return CallRet(ctx, ret, resp)
}

// CallWithJsonReturnResp JSON请求返回resp
func (r Client) CallWithJsonReturnResp(ctx context.Context, retBuf *bytes.Buffer, method, reqUrl string, headers http.Header,
	param interface{}) (err error) {
	if DebugMode {
		ctx = context.WithValue(ctx, insReqStartTimeKey, time.Now())
		ctx = context.WithValue(ctx, insReqUrlKey, reqUrl)
	}
	resp, err := r.DoRequestWithJson(ctx, method, reqUrl, headers, param)
	if err != nil {
		return err
	}
	return CallRetResp(ctx, retBuf, resp)
}

// DoRequestWith64 请求
func (r Client) DoRequestWith64(ctx context.Context, method, reqUrl string, headers http.Header, body io.Reader,
	bodyLength int64) (resp *http.Response, err error) {

	req, err := newRequest(ctx, method, reqUrl, headers, body)
	if err != nil {
		return
	}
	req.ContentLength = bodyLength
	return r.Do(ctx, req)
}

// CallWith64 请求
func (r Client) CallWith64(ctx context.Context, ret interface{}, method, reqUrl string, headers http.Header, body io.Reader,
	bodyLength int64) (err error) {
	if DebugMode {
		ctx = context.WithValue(ctx, insReqStartTimeKey, time.Now())
		ctx = context.WithValue(ctx, insReqUrlKey, reqUrl)
	}
	resp, err := r.DoRequestWith64(ctx, method, reqUrl, headers, body, bodyLength)
	if err != nil {
		return err
	}
	return CallRet(ctx, ret, resp)
}

type requestCanceler interface {
	CancelRequest(req *http.Request)
}

type nestedObjectGetter interface {
	NestedObject() interface{}
}

func getRequestCanceler(tp http.RoundTripper) (rc requestCanceler, ok bool) {
	if rc, ok = tp.(requestCanceler); ok {
		return
	}

	p := interface{}(tp)
	for {
		getter, ok1 := p.(nestedObjectGetter)
		if !ok1 {
			return
		}
		p = getter.NestedObject()
		if rc, ok = p.(requestCanceler); ok {
			return
		}
	}
}
