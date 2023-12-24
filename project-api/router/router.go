package router

import (
	"com.levi/project-api/api/user"
	"github.com/gin-gonic/gin"
)

type Router interface {
	Route(r *gin.Engine)
}

type RouterRegister struct {
}

func New() *RouterRegister {
	return &RouterRegister{}
}

func (*RouterRegister) Route(ro Router, r *gin.Engine) {
	ro.Route(r)
}

func InitRouter(r *gin.Engine) {
	rg := New()
	rg.Route(&user.RouterUser{}, r)
}