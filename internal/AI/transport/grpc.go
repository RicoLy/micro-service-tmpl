package transport

import (
	"context"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
	"micro-service-tmpl/internal/AI/domain/global"
	"micro-service-tmpl/internal/AI/domain/pb"
	"micro-service-tmpl/internal/AI/domain/req"
	"micro-service-tmpl/internal/AI/domain/res"
	"micro-service-tmpl/internal/AI/endpoint"
	"micro-service-tmpl/utils/log"
)

//实现protobuf中定义的接口
type grpcServer struct {
	showAiDistort   grpctransport.Handler
	addAiDistort    grpctransport.Handler
	deleteAiDistort grpctransport.Handler
}

func NewGRPCServer(endpoints endpoint.AiEndpoints, log *zap.Logger) pb.AiServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerBefore(func(ctx context.Context, md metadata.MD) context.Context {
			ctx = context.WithValue(ctx, global.ContextReqUUid, md.Get(global.ContextReqUUid))
			return ctx
		}),
		grpctransport.ServerErrorHandler(NewZapLogErrorHandler(log)),
	}

	return &grpcServer{
		showAiDistort: grpctransport.NewServer( //获取Ai 误报/漏报配置
			endpoints.ShowAiDistortEndpoint,
			(&req.ShowVO{}).GRPCServerRequestDecode,
			(&res.ShowAiDistortRsp{}).GRPCResponseEncode,
			options...,
		),
		addAiDistort: grpctransport.NewServer( //添加Ai误报信息
			endpoints.AddAiDistortEndpoint,
			(&req.AddAiDistortVO{}).GRPCServerRequestDecode,
			(&res.Ack{}).GRPCResponseEncode,
			options...,
		),
		deleteAiDistort: grpctransport.NewServer( //删除Ai误报信息
			endpoints.DeleteAiDistortEndpoint,
			(&req.DeleteAiDistortVO{}).GRPCServerRequestDecode,
			(&res.Ack{}).GRPCResponseEncode,
			options...,
		),
	}
}

func (s *grpcServer) ShowAiDistort(ctx context.Context, req *pb.ShowReq) (*pb.ShowAiDistortRsp, error) {
	_, rep, err := s.showAiDistort.ServeGRPC(ctx, req)
	if err != nil {
		log.GetLogger().Warn("s.ShowAiDistort.ServeGRPC", zap.Error(err))
		return nil, err
	}
	return rep.(*pb.ShowAiDistortRsp), nil
}

func (s *grpcServer) AddAiDistort(ctx context.Context, req *pb.AddAiDistortReq) (*pb.Ack, error) {
	_, rep, err := s.addAiDistort.ServeGRPC(ctx, req)
	if err != nil {
		log.GetLogger().Warn("s.AddAiDistort.ServeGRPC", zap.Error(err))
		return nil, err
	}
	return rep.(*pb.Ack), nil
}

func (s *grpcServer) DeleteAiDistort(ctx context.Context, req *pb.DeleteAiDistortReq) (*pb.Ack, error) {
	_, rep, err := s.deleteAiDistort.ServeGRPC(ctx, req)
	if err != nil {
		log.GetLogger().Warn("s.DeleteAiDistort.ServeGRPC", zap.Error(err))
		return nil, err
	}
	return rep.(*pb.Ack), nil
}
