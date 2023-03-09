package utils

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/youchuangcd/gopkg"
)

type ResponseBodyWriter struct {
	gin.ResponseWriter
	Body      *bytes.Buffer
	RealWrite bool
}

func (r *ResponseBodyWriter) Write(b []byte) (int, error) {
	n, err := r.Body.Write(b)
	if r.RealWrite {
		return r.ResponseWriter.Write(b)
	}
	return n, err
}

func ReplaceResponseBodyWriter(c *gin.Context) *ResponseBodyWriter {
	var (
		writer *ResponseBodyWriter
		ok     bool
	)
	if writer, ok = c.Value(gopkg.ContextResponseBodyWriterKey).(*ResponseBodyWriter); !ok {
		writer = &ResponseBodyWriter{
			ResponseWriter: c.Writer,
			Body:           &bytes.Buffer{},
			RealWrite:      true,
		}
		c.Writer = writer
		c.Set(gopkg.ContextResponseBodyWriterKey, writer)
	}
	return writer
}
