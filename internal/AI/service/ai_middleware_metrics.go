package service

import (
	"context"
	"github.com/go-kit/kit/metrics"
	"micro-service-tmpl/internal/AI/domain/req"
	"micro-service-tmpl/internal/AI/domain/res"
	"time"
)

// 服务监控中间件
type metricsMiddlewareServer struct {
	next      AIService
	counter   metrics.Counter
	histogram metrics.Histogram
}

func NewMetricsMiddlewareServer(counter metrics.Counter, histogram metrics.Histogram) NewMiddlewareServer {
	return func(service AIService) AIService {
		return metricsMiddlewareServer{
			next:      service,
			counter:   counter,
			histogram: histogram,
		}
	}
}

func (m metricsMiddlewareServer) ShowAiDistort(ctx context.Context, req *req.ShowVO) (rsp *res.ShowAiDistortRsp, err error) {
	defer func(start time.Time) {
		method := []string{"method", "ShowAiDistort"}
		m.counter.With(method...).Add(1)
		m.histogram.With(method...).Observe(time.Since(start).Seconds())
	}(time.Now())
	return m.next.ShowAiDistort(ctx, req)
}

func (m metricsMiddlewareServer) AddAiDistort(ctx context.Context, req *req.AddAiDistortVO) (rsp *res.Ack, err error) {
	defer func(start time.Time) {
		method := []string{"method", "AddAiDistort"}
		m.counter.With(method...).Add(1)
		m.histogram.With(method...).Observe(time.Since(start).Seconds())
	}(time.Now())
	return m.next.AddAiDistort(ctx, req)
}

func (m metricsMiddlewareServer) DeleteAiDistort(ctx context.Context, req *req.DeleteAiDistortVO) (rsp *res.Ack, err error) {
	defer func(start time.Time) {
		method := []string{"method", "DeleteAiDistort"}
		m.counter.With(method...).Add(1)
		m.histogram.With(method...).Observe(time.Since(start).Seconds())
	}(time.Now())
	return m.next.DeleteAiDistort(ctx, req)
}
