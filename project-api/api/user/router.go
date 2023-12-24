package user

import (
	"github.com/gin-gonic/gin"
)

type RouterUser struct {
}

func (*RouterUser) Route(r *gin.Engine) {
	r.GET("/login", (&HandlerUser{}).login)
	r.GET("/getUser", (&HandlerUser{}).getUser)
}