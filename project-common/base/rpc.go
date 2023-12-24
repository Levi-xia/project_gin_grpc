package base

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// rpc请求拦截器
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		// 接口耗时统计起始
		start := time.Now()
		// 获取元数据
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			// 没有元数据的错误处理
			return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
		}
		requestIds, ok := md["request-id"]
		if !ok || len(requestIds) == 0 {
			// 没有requestId的错误处理
			return nil, status.Errorf(codes.InvalidArgument, "missing request ID")
		}
		requestId := requestIds[0]

		// 调用实际的处理函数
		resp, err = handler(ctx, req)

		// 在处理之后记录响应
		zap.L().Info(
			"RpcLog",
			zap.String("requestId", requestId),
			zap.Reflect("input", req),
			zap.Reflect("output", resp),
			zap.String("method", info.FullMethod),
			zap.String("cost", time.Since(start).String()),
			zap.Error(err),
		)
		return resp, err
	}
}
