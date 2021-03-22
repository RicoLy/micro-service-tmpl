package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"micro-service-tmpl/internal/AI/domain/req"
	"micro-service-tmpl/internal/AI/domain/res"
	"micro-service-tmpl/internal/AI/service"
)

// 定义endpoints
type AiEndpoints struct {
	ShowAiDistortEndpoint   endpoint.Endpoint
	AddAiDistortEndpoint    endpoint.Endpoint
	DeleteAiDistortEndpoint endpoint.Endpoint
}

// 构造endpoints
func NewAiEndpoints(srv service.AIService, log *zap.Logger, limit *rate.Limiter) AiEndpoints {
	var showAiDistortEndpoint endpoint.Endpoint
	{
		showAiDistortEndpoint = MakeShowAiDistortEndPoint(srv)
		// 装饰服务限流中间件
		showAiDistortEndpoint = GolangRateAllowMiddleware(limit)(showAiDistortEndpoint)
		// 装饰日志中间件
		showAiDistortEndpoint = LoggingMiddleware(log)(showAiDistortEndpoint)
	}

	var addAiDistortEndpoint endpoint.Endpoint
	{
		addAiDistortEndpoint = MakeAddAiDistortEndPoint(srv)
		// 装饰服务限流中间件
		addAiDistortEndpoint = GolangRateAllowMiddleware(limit)(addAiDistortEndpoint)
		// 装饰日志中间件
		addAiDistortEndpoint = LoggingMiddleware(log)(addAiDistortEndpoint)
	}

	var deleteAiDistortEndpoint endpoint.Endpoint
	{
		deleteAiDistortEndpoint = MakeDeleteAiDistortEndPoint(srv)
		// 装饰服务限流中间件
		deleteAiDistortEndpoint = GolangRateAllowMiddleware(limit)(deleteAiDistortEndpoint)
		// 装饰日志中间件
		deleteAiDistortEndpoint = LoggingMiddleware(log)(deleteAiDistortEndpoint)
	}

	return AiEndpoints{
		ShowAiDistortEndpoint:   showAiDistortEndpoint,
		AddAiDistortEndpoint:    addAiDistortEndpoint,
		DeleteAiDistortEndpoint: deleteAiDistortEndpoint,
	}
}

//获取Ai 误报/漏报配置endpoint
func (e *AiEndpoints) ShowAiDistort(ctx context.Context, req *req.ShowVO) (rsp *res.ShowAiDistortRsp, err error) {
	resp, err := e.ShowAiDistortEndpoint(ctx, req)
	if err != nil {
		return
	}
	return resp.(*res.ShowAiDistortRsp), err
}

//构造 获取Ai 误报/漏报配置endpoint
func MakeShowAiDistortEndPoint(svc service.AIService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		re := request.(*req.ShowVO)
		return svc.ShowAiDistort(ctx, re)
	}
}

//添加 Ai误报配置endpoint
func (e *AiEndpoints) AddAiDistort(ctx context.Context, req *req.AddAiDistortVO) (rsp *res.Ack, err error) {
	resp, err := e.AddAiDistortEndpoint(ctx, req)
	if err != nil {
		return
	}
	return resp.(*res.Ack), err
}

//构造 添加 Ai误报配置endpoint
func MakeAddAiDistortEndPoint(svc service.AIService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		re := request.(*req.AddAiDistortVO)
		return svc.AddAiDistort(ctx, re)
	}
}

//删除 Ai误报配置endpoint
func (e *AiEndpoints) DeleteAiDistort(ctx context.Context, req *req.DeleteAiDistortVO) (rsp *res.Ack, err error) {
	resp, err := e.DeleteAiDistortEndpoint(ctx, req)
	if err != nil {
		return
	}
	return resp.(*res.Ack), err
}

//构造 删除 Ai误报配置endpoint
func MakeDeleteAiDistortEndPoint(svc service.AIService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		re := request.(*req.DeleteAiDistortVO)
		return svc.DeleteAiDistort(ctx, re)
	}
}
