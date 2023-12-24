package router

import (
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
}