package user

import (
	"com.levi/project-api/rpc"
	"com.levi/project-common/base"
	"com.levi/project-common/utils"
	userServiceV1 "com.levi/project-user/pkg/service/user.service.v1"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HandlerUser struct {
}

func (*HandlerUser) login(ctx *gin.Context) {
	rsp := base.Result{}
	ctx.Set("checkLogin", true)
	username := ctx.Query("username")
	password := ctx.Query("password")

	response, err := rpc.UserServiceRpcClient.Login(ctx.Request.Context(), &userServiceV1.LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, rsp.Error(2001, err.Error()))
		return
	}
	token := response.Token
	ctx.SetCookie("AT", token, 86400, "/", "127.0.0.1", false, false)

	ctx.JSON(http.StatusOK, rsp.Success(response))
}

func (*HandlerUser) getUser(ctx *gin.Context) {
	rsp := base.Result{}
	uid := ctx.Query("userId")
	uid_int, _ := utils.StringToInt64(uid)
	response, err := rpc.UserServiceRpcClient.GetUser(ctx.Request.Context(), &userServiceV1.UserRequest{
		UserId: uid_int,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, rsp.Error(2001, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, rsp.Success(response))
}
