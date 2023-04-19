package middleware

import (
	"fmt"
	"threads-service/global"
	"threads-service/pkg/app"
	"threads-service/pkg/email"
	"threads-service/pkg/errcode"
	"time"

	"github.com/gin-gonic/gin"
)

// func Recovery2() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		defer func() {
// 			err := recover()
// 			if err != nil {
// 				global.Logger.WithCallesFrames().Errf("panic recover err: %v", err)
// 				app.NewResponse(c).ToErrorResponse(errcode.ServerError)
// 				c.Abort()
// 			}
// 		}()
// 		c.Next()
// 	}
// }

func Recovery() gin.HandlerFunc {
	mail := email.NewEmail(&email.SMTPInfo{
		Host:     global.EmailSetting.Host,
		Port:     global.EmailSetting.Port,
		IsSSL:    global.EmailSetting.IsSSL,
		UserName: global.EmailSetting.UserName,
		Password: global.EmailSetting.Password,
		From:     global.EmailSetting.From,
	})
	return func(c *gin.Context) {
		defer func() {
			err := recover()
			if err != nil {
				global.Logger.WithCallesFrames().Errf(c, "err:=recover() err:%v", err)
				err := mail.SendEmail(global.EmailSetting.To, fmt.Sprintf("异常发生时间:%s", time.Now().Format("2008-05-12 14:14:12")), fmt.Sprintf("错误信息:%v", err))
				if err != nil {
					global.Logger.Panicf(c, "mail.SendEmail err:%v", err)
				}
				app.NewResponse(c).ToErrorResponse(errcode.ServerError)
				c.Abort()
			}
		}()
		c.Next()
	}
}
