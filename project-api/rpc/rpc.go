package rpc

import (
	"log"
	userServiceV1 "com.levi/project-user/pkg/service/user.service.v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// UserServiceRpcClient 用户服务客户端
var UserServiceRpcClient userServiceV1.UserServiceClient

func InitRpcClient() {
	InitUserServiceGrpcClient()
}

func InitUserServiceGrpcClient() {
	conn, err := grpc.Dial("127.0.0.1:8881", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	UserServiceRpcClient = userServiceV1.NewUserServiceClient(conn)
}