package middleware

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

// 初始化全局请求traceId
func RequestIDHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("trace")

		requestID := uuid.New().String()
		// 将requestID放入上下文中
		c.Set("RequestID", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)
		// 创建元数据
		md := metadata.New(map[string]string{
			"request-id": requestID,
		})
		// 创建带有超时的上下文
		timeout, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()
		// 将元数据和超时添加到新的上下文
		newCtx := metadata.NewOutgoingContext(timeout, md)
		// 将新的上下文设置到Gin的请求中
		c.Request = c.Request.WithContext(newCtx)
		c.Next()
	}
}