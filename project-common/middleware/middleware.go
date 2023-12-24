package middleware

import "github.com/gin-gonic/gin"

func InitMiddleware(r *gin.Engine) {
	r.Use(ErrorHandler(), RequestIDHandler(), ActionHandler(), JWTAuthHandler())
}
