package middleware

import (
	"fmt"
	"log"
	"net/http"

	"com.levi/project-common/base"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			log.Println("error")
			// recover() 函数捕获 panic
			if r := recover(); r != nil {
				// 日志记录错误
				zap.L().Error("Recovered", zap.Any("panic", r))

				// 通知客户端发生了错误
				c.JSON(http.StatusOK, &base.Result{
					Code:    999,
					Msg: "系统错误:" + fmt.Sprint(r),
					Data:   nil,
				})
				// 终止后续的处理函数调用
				c.Abort()
				return;
			}
		}()
		// 调用下一个中间件或处理函数
		c.Next()
	}
}