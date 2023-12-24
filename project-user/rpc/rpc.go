package rpc

import (
	"log"
	"net"

	"com.levi/project-common/base"
	"com.levi/project-user/config"
	userServiceV1 "com.levi/project-user/pkg/service/user.service.v1"
	"google.golang.org/grpc"
)

type GrpcConfig  struct {
	Addr string
	RegisterFunc func(*grpc.Server)
}

func InitGrpcServer() *grpc.Server {
	c := &GrpcConfig{
		Addr: config.Conf.Grpc.Addr,
		RegisterFunc: func(server *grpc.Server) {
			userServiceV1.RegisterUserServiceServer(server, &userServiceV1.UserService{})
		},
	}
	s := grpc.NewServer(
		grpc.UnaryInterceptor(base.UnaryServerInterceptor()),
	)
	c.RegisterFunc(s)

	lis, err := net.Listen("tcp", config.Conf.Grpc.Addr)
	if err != nil {
		log.Fatalln("Fatal error listen: ", err);
	}
	go func() {
		if err = s.Serve(lis); err != nil {
			log.Fatalln("Fatal error serve: ", err);
		}
	} ()
	return s
}

