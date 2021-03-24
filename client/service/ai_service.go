package service

import (
	"context"
	"micro-service-tmpl/client/domain/req"
	"micro-service-tmpl/client/domain/res"
)

// Ai业务接口
type AIService interface {
	ShowAiDistort(ctx context.Context, req *req.ShowVO) (rsp *res.ShowAiDistortRsp, err error)
	AddAiDistort(ctx context.Context, req *req.AddAiDistortVO) (rsp *res.Ack, err error)
	DeleteAiDistort(ctx context.Context, req *req.DeleteAiDistortVO) (rsp *res.Ack, err error)
}


