package service

import (
	"context"
	"github.com/go-kit/kit/metrics"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"micro-service-tmpl/internal/AI/dao"
	"micro-service-tmpl/internal/AI/domain/req"
	"micro-service-tmpl/internal/AI/domain/res"
	"micro-service-tmpl/internal/AI/domain/vo"
	"micro-service-tmpl/utils/generatorId"
	"strings"
)

// Ai业务接口
type AIService interface {
	ShowAiDistort(ctx context.Context, req *req.ShowVO) (rsp *res.ShowAiDistortRsp, err error)
	AddAiDistort(ctx context.Context, req *req.AddAiDistortVO) (rsp *res.Ack, err error)
	DeleteAiDistort(ctx context.Context, req *req.DeleteAiDistortVO) (rsp *res.Ack, err error)
}

type AIServiceImpl struct {
	aiDistortDAO dao.AiDistort
	logger       *zap.Logger
}

func NewAIServiceImpl(aiDistortDAO dao.AiDistort, counter metrics.Counter, histogram metrics.Histogram, log *zap.Logger, tracer opentracing.Tracer) AIService {
	var service AIService
	service = &AIServiceImpl{
		aiDistortDAO: aiDistortDAO,
		logger:       log,
	}

	// 装饰日志中间件
	service = NewLogMiddlewareServer(log)(service)

	// 装饰服务监控中间件
	service = NewMetricsMiddlewareServer(counter, histogram)(service)

	// 装饰jaeger链路追踪中间件
	if tracer != nil {
		service = NewTracerMiddlewareServer(tracer)(service)
	}

	return service
}

//获取Ai 误报/漏报配置
func (s AIServiceImpl) ShowAiDistort(ctx context.Context, req *req.ShowVO) (rsp *res.ShowAiDistortRsp, err error) {
	rsp = new(res.ShowAiDistortRsp)
	var pagination vo.Pagination
	var whereStr = ""
	var param = make([]interface{}, 0)
	pagination.Page = req.Page
	pagination.PageSize = req.PageSize
	if req.Where != "" {
		kvs := strings.Split(req.Where, "&")
		for _, kv := range kvs {
			skv := strings.Split(kv, "=")
			if whereStr != "" {
				whereStr += " and "
			}
			whereStr = whereStr + skv[0] + " = ?"
			param = append(param, skv[1])
		}
	}

	aiDistort, err := s.aiDistortDAO.Find(whereStr, &pagination, param...)
	if err != nil {
		return
	}

	rsp.AiDistorts = res.BuildDistorts(aiDistort)
	rsp.Pagination = pagination

	return
}

//添加Ai误报信息
func (s AIServiceImpl) AddAiDistort(ctx context.Context, req *req.AddAiDistortVO) (rsp *res.Ack, err error) {
	rsp = new(res.Ack)
	if err = s.aiDistortDAO.Add(&dao.AiDistort{
		Uin:     generatorId.NextId(), // 雪花算法生成id
		Appid:   req.Appid,
		Domain:  req.Domain,
		Payload: req.Payload,
		From:    req.From,
		Remark:  req.Remark,
		Status:  dao.NotLearned, // 默认未学习
	}); err != nil {
		rsp.IsOk = false
		return
	}
	rsp.IsOk = true
	return
}

//删除Ai误报信息
func (s AIServiceImpl) DeleteAiDistort(ctx context.Context, req *req.DeleteAiDistortVO) (rsp *res.Ack, err error) {
	rsp = new(res.Ack)
	if err = s.aiDistortDAO.Delete("uin = ? ", req.Uin); err != nil {
		rsp.IsOk = false
		return
	}
	rsp.IsOk = true
	return
}
