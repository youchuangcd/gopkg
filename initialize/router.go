package initialize

import (
	"github.com/gin-gonic/gin"
	"github.com/youchuangcd/gopkg/middleware"
)

func Router(r *gin.Engine, serviceName string) {
	r.Use(middleware.Recover)
	//注册中间件
	r.Use(middleware.TraceId)
	//r.Use(otelgin.Middleware(serviceName, otelgin.WithFilter(func(r *http.Request) bool {
	//	if r.RequestURI == "/app/healthCheck" {
	//		return false
	//	}
	//	return true
	//})))
}
