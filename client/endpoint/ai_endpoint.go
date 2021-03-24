package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"micro-service-tmpl/client/domain/req"
	"micro-service-tmpl/client/domain/res"
	"micro-service-tmpl/client/service"
)

// 定义endpoints
type AiEndpoints struct {
	ShowAiDistortEndpoint   endpoint.Endpoint
	AddAiDistortEndpoint    endpoint.Endpoint
	DeleteAiDistortEndpoint endpoint.Endpoint
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
