package middleware

import (
	"bytes"
	"threads-service/global"
	"threads-service/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
)

type AccessLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w AccessLogWriter) Write(p []byte) (int, error) {
	n, err := w.body.Write(p)
	if err != nil {
		return n, err
	}
	return w.ResponseWriter.Write(p)
}

func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		bodyWriter := &AccessLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		// 只需写一个针对访问日志的Writer结构体，实现特定的Write方法就可以解决无法直接获取方法响应主体的问题了
		c.Writer = bodyWriter
		begin := time.Now().Unix()
		c.Next()
		end := time.Now().Unix()

		fields := logger.Fields{
			"request":  c.Request.PostForm.Encode(),
			"response": bodyWriter.body.String(),
		}
		global.Logger.WithFields(fields).Infof(c, "access log: method: %s, satus_code: %d, begin_time: %d, end_time: %d", c.Request.Method, bodyWriter.Status(), begin, end)
	}
}
