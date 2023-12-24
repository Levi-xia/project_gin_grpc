package middleware

import (
	"fmt"
	"log"
	"net/http"

	"com.levi/project-common/base"
	"com.levi/project-common/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func JWTAuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			log.Println("jwt")

			tokenStr, _ := c.Cookie("AT")
			if tokenStr != "" {
				token, err := jwt.ParseWithClaims(tokenStr, &base.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return jwt.ErrInvalidKey, nil
					}
					return []byte(config.GlobalConf.Jwt.Secret), nil
				})
				if err != nil {
					zap.L().Error("Jwt", zap.Any("panic", err))

					c.JSON(http.StatusOK, &base.Result{
						Code: 999,
						Msg:  "系统错误:" + fmt.Sprint(err),
						Data: nil,
					})
					c.Abort()
					return
				}
				claims := token.Claims.(*base.CustomClaims)
				c.Set("userId", claims.Id)
			}
		}()
		c.Next()
	}
}
