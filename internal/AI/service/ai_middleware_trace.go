package service

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"micro-service-tmpl/internal/AI/domain/req"
	"micro-service-tmpl/internal/AI/domain/res"
)

type tracerMiddlewareServer struct {
	tracer opentracing.Tracer
	next   AIService
}

func NewTracerMiddlewareServer(tracer opentracing.Tracer) NewMiddlewareServer {
	return func(service AIService) AIService {
		return tracerMiddlewareServer{
			tracer: tracer,
			next:   service,
		}
	}
}

func (tm tracerMiddlewareServer) ShowAiDistort(ctx context.Context, req *req.ShowVO) (rsp *res.ShowAiDistortRsp, err error) {
	span, ctxContext := opentracing.StartSpanFromContextWithTracer(ctx, tm.tracer, "ShowAiDistort", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "NewTracerServerMiddleware",
	})
	defer func() {
		span.LogKV("req", req)
		span.Finish()
	}()
	return tm.next.ShowAiDistort(ctxContext, req)
}

func (tm tracerMiddlewareServer) AddAiDistort(ctx context.Context, req *req.AddAiDistortVO) (rsp *res.Ack, err error) {
	span, ctxContext := opentracing.StartSpanFromContextWithTracer(ctx, tm.tracer, "AddAiDistort", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "NewTracerServerMiddleware",
	})
	defer func() {
		span.LogKV("req", req)
		span.Finish()
	}()
	return tm.next.AddAiDistort(ctxContext, req)
}

func (tm tracerMiddlewareServer) DeleteAiDistort(ctx context.Context, req *req.DeleteAiDistortVO) (rsp *res.Ack, err error) {
	span, ctxContext := opentracing.StartSpanFromContextWithTracer(ctx, tm.tracer, "DeleteAiDistort", opentracing.Tag{
		Key:   string(ext.Component),
		Value: "NewTracerServerMiddleware",
	})
	defer func() {
		span.LogKV("req", req)
		span.Finish()
	}()
	return tm.next.DeleteAiDistort(ctxContext, req)
}
