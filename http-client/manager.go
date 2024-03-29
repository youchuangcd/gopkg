package httpclient

import (
	"bytes"
	"context"
	"github.com/youchuangcd/gopkg"
	"github.com/youchuangcd/gopkg/http-client/auth"
	"github.com/youchuangcd/gopkg/http-client/client"
	"io"
	"net/http"
	"sync"
	"time"
)

var (
	defaultManager *Manager
	once           sync.Once
	DefaultTimeout = gopkg.HttpClientTimeout // 客户端默认超时
)

type Manager struct {
	Client      *client.Client
	Credentials *auth.Credentials
}

// NewCredentials 获取认证
func NewCredentials(clientKey, clientSecret string) *auth.Credentials {
	return auth.New(clientKey, clientSecret)
}

// DefaultManager
// @Description: 获取默认http管理对象
// @return *Manager
func DefaultManager() *Manager {
	once.Do(func() {
		defaultManager = NewManager(nil)
	})
	return defaultManager
}

func NewManager(credentials *auth.Credentials) *Manager {
	c := client.DefaultClient
	//credentials := newCredentials(key, secret)
	c.Transport = newTransport(credentials, nil)
	c.Timeout = DefaultTimeout
	return &Manager{
		Client:      &c,
		Credentials: credentials,
	}
}

func NewManagerV2(credentials *auth.Credentials, timeout time.Duration, tr http.RoundTripper) *Manager {
	c := client.DefaultClient
	c.Transport = newTransport(credentials, tr)
	c.Timeout = timeout
	return &Manager{
		Client:      &c,
		Credentials: credentials,
	}
}

type transport struct {
	http.RoundTripper
	credentials *auth.Credentials
}

func newTransport(credentials *auth.Credentials, tr http.RoundTripper) *transport {
	if tr == nil {
		tr = gopkg.HttpClientDefaultTransport
	}
	return &transport{tr, credentials}
}

func CallWithJson(ctx context.Context, ret interface{}, method, reqUrl string, headers http.Header, param interface{}) (err error) {
	return DefaultManager().Client.CallWithJson(ctx, ret, method, reqUrl, headers, param)
}

func CallWithForm(ctx context.Context, ret interface{}, method, reqUrl string, headers http.Header, param interface{}) (err error) {
	return DefaultManager().Client.CallWithForm(ctx, ret, method, reqUrl, headers, param)
}

func CallWithJsonReturnResp(ctx context.Context, retBuf *bytes.Buffer, method, reqUrl string, headers http.Header, param interface{}) (err error) {
	return DefaultManager().Client.CallWithJsonReturnResp(ctx, retBuf, method, reqUrl, headers, param)
}

func DoRequestWith64(ctx context.Context, method, reqUrl string, headers http.Header, body io.Reader, bodyLength int64) (resp *http.Response, err error) {
	return DefaultManager().Client.DoRequestWith64(ctx, method, reqUrl, headers, body, bodyLength)
}

// CallWith64 请求
func CallWith64(ctx context.Context, ret interface{}, method, reqUrl string, headers http.Header, body io.Reader, bodyLength int64) (err error) {
	return DefaultManager().Client.CallWith64(ctx, ret, method, reqUrl, headers, body, bodyLength)
}
