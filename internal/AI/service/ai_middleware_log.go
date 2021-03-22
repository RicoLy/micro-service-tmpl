package service

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"micro-service-tmpl/internal/AI/domain/global"
	"micro-service-tmpl/internal/AI/domain/req"
	"micro-service-tmpl/internal/AI/domain/res"
)

type NewMiddlewareServer func(service AIService) AIService

// 服务日志中间件
type logMiddlewareServer struct {
	logger *zap.Logger
	next   AIService
}

func NewLogMiddlewareServer(log *zap.Logger) NewMiddlewareServer {
	return func(service AIService) AIService {
		return logMiddlewareServer{
			logger: log,
			next:   service,
		}
	}
}

func (l logMiddlewareServer) ShowAiDistort(ctx context.Context, req *req.ShowVO) (rsp *res.ShowAiDistortRsp, err error) {
	defer func() {
		l.logger.Debug(fmt.Sprint(ctx.Value(global.ContextReqUUid)), zap.Any("调用 ShowAiDistort logMiddlewareServer", "ShowAiDistort"), zap.Any("req", req), zap.Any("res", rsp))
	}()
	return l.next.ShowAiDistort(ctx, req)
}

func (l logMiddlewareServer) AddAiDistort(ctx context.Context, req *req.AddAiDistortVO) (rsp *res.Ack, err error) {
	defer func() {
		l.logger.Debug(fmt.Sprint(ctx.Value(global.ContextReqUUid)), zap.Any("调用 AddAiDistort logMiddlewareServer", "AddAiDistort"), zap.Any("req", req), zap.Any("res", err))
	}()
	return l.next.AddAiDistort(ctx, req)
}

func (l logMiddlewareServer) DeleteAiDistort(ctx context.Context, req *req.DeleteAiDistortVO) (rsp *res.Ack, err error) {
	defer func() {
		l.logger.Debug(fmt.Sprint(ctx.Value(global.ContextReqUUid)), zap.Any("调用 DeleteAiDistort logMiddlewareServer", "DeleteAiDistort"), zap.Any("req", req), zap.Any("res", err))
	}()
	return l.next.DeleteAiDistort(ctx, req)
}
