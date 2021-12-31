package interceptor

import (
	"context"

	"go.uber.org/zap"

	"github.com/donech/pearl/pkg/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func Log() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		var traceId string
		// Read metadata from client.
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			log.S(ctx).Debugf("incoming md is: %v", md)
			traceId = log.GetTraceIDFromGrpcMetadata(md)
			if traceId != "" {
				log.S(ctx).Infof("grpc incoming header with trace id %s", traceId)
			}
		}
		if traceId == "" {
			log.S(ctx).Debug("trace-id not found, so generate it")
			traceId = log.NewTraceID()
		}
		ctx = context.WithValue(ctx, log.TraceKeyName, traceId)
		header := metadata.New(map[string]string{string(log.TraceKeyName): traceId})
		err = grpc.SetHeader(ctx, header)
		if err != nil {
			log.S(ctx).Errorf("grpc send header error %v", err)
		}

		log.L(ctx).Info("incoming grpc req", zap.Reflect("req", req))
		resp, err = handler(ctx, req)
		log.L(ctx).Info("output grpc resp", zap.Reflect("resp", resp), zap.Error(err))

		return resp, err
	}
}
