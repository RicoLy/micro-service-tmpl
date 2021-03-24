package Ai

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/etcdv3"
	"github.com/go-kit/kit/sd/lb"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	"micro-service-tmpl/client/domain/global"
	"micro-service-tmpl/client/domain/pb"
	"micro-service-tmpl/client/domain/req"
	"micro-service-tmpl/client/domain/res"
	"micro-service-tmpl/client/domain/vo"
	endpoint2 "micro-service-tmpl/client/endpoint"
	"micro-service-tmpl/client/service"
	"micro-service-tmpl/utils/jaegerTracer"
	"micro-service-tmpl/utils/myLog"

	"time"
)

// AiRpc客户端
type AiAgent struct {
	instance *etcdv3.Instancer
	logger   log.Logger
	tracer   opentracing.Tracer
}

func NewAiAgentClient(addr []string, logger log.Logger, jaegerAddr string) (*AiAgent, error) {
	var (
		sEtcdAddr = addr
		serName   = "svc.Ai"
		ttl       = 5 * time.Second
	)
	options := etcdv3.ClientOptions{
		DialKeepAlive: ttl,
		DialTimeout:   ttl,
	}
	tracer, _, err := jaegerTracer.NewJaegerTracer("user_agent_client", jaegerAddr)
	if err != nil {
		return nil, err
	}
	etcdClient, err := etcdv3.NewClient(context.Background(), sEtcdAddr, options)
	if err != nil {
		return nil, err
	}
	instance, err := etcdv3.NewInstancer(etcdClient, serName, logger)
	if err != nil {
		return nil, err
	}
	return &AiAgent{
		instance: instance,
		logger:   logger,
		tracer:   tracer,
	}, err
}

func (a *AiAgent) AiAgentClient() (service.AIService, error) {
	var (
		retryMax     = 3
		retryTimeout = 5 * time.Second

		endpoints endpoint2.AiEndpoints
	)

	{
		factory := a.factoryFor(endpoint2.MakeShowAiDistortEndPoint)
		defaultEndPoint := sd.NewEndpointer(a.instance, factory, a.logger)
		balancer := lb.NewRoundRobin(defaultEndPoint)
		showAiDistortEndPoint := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.ShowAiDistortEndpoint = showAiDistortEndPoint
	}

	{
		factory := a.factoryFor(endpoint2.MakeAddAiDistortEndPoint)
		defaultEndPoint := sd.NewEndpointer(a.instance, factory, a.logger)
		balancer := lb.NewRoundRobin(defaultEndPoint)
		addAiDistortEndPoint := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.AddAiDistortEndpoint = addAiDistortEndPoint
	}

	{
		factory := a.factoryFor(endpoint2.MakeDeleteAiDistortEndPoint)
		defaultEndPoint := sd.NewEndpointer(a.instance, factory, a.logger)
		balancer := lb.NewRoundRobin(defaultEndPoint)
		deleteAiDistortEndPoint := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.DeleteAiDistortEndpoint = deleteAiDistortEndPoint
	}

	return &endpoints, nil
}

func (a *AiAgent) factoryFor(makeEndpoint func(service service.AIService) endpoint.Endpoint) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		fmt.Println("instance >>>>>>>>>>>>>>>>   ", instance)
		chainUnaryServer := grpcmiddleware.ChainUnaryClient(
			grpc_opentracing.UnaryClientInterceptor(grpc_opentracing.WithTracer(a.tracer)),
			grpc_zap.UnaryClientInterceptor(myLog.GetLogger()),
			jaegerTracer.JaegerClientMiddleware(a.tracer),
		)
		conn, err := grpc.Dial(
			instance,
			grpc.WithInsecure(),
			grpc.WithUnaryInterceptor(chainUnaryServer),
		)
		if err != nil {
			return nil, nil, err
		}
		srv := a.NewGRPCClient(conn)

		endpoints := makeEndpoint(srv)
		return endpoints, conn, err
	}
}

func (a *AiAgent) NewGRPCClient(conn *grpc.ClientConn) service.AIService {
	options := []grpctransport.ClientOption{
		grpctransport.ClientBefore(func(ctx context.Context, md *metadata.MD) context.Context {
			UUID := uuid.NewV5(uuid.Must(uuid.NewV4(), nil), "req_uuid").String()
			md.Set(global.ContextReqUUid, UUID)
			ctx = metadata.NewOutgoingContext(ctx, *md)
			return ctx
		}),
	}

	var showAiDistortEndpoint endpoint.Endpoint
	{
		showAiDistortEndpoint = grpctransport.NewClient(
			conn,
			"pb.Ai",
			"ShowAiDistort",
			a.RequestShowAiDistort,
			a.ResponseShowAiDistort,
			pb.ShowAiDistortRsp{},
			options...,
		).Endpoint()
	}

	var addAiDistortEndpoint endpoint.Endpoint
	{
		addAiDistortEndpoint = grpctransport.NewClient(
			conn,
			"pb.Ai",
			"AddAiDistort",
			a.RequestAddAiDistort,
			a.ResponseAddAiDistort,
			pb.Ack{},
			options...,
		).Endpoint()
	}

	var deleteAiDistortEndpoint endpoint.Endpoint
	{
		deleteAiDistortEndpoint = grpctransport.NewClient(
			conn,
			"pb.Ai",
			"DeleteAiDistort",
			a.RequestDeleteAiDistort,
			a.ResponseDeleteAiDistort,
			pb.Ack{},
			options...,
		).Endpoint()
	}

	return &endpoint2.AiEndpoints{
		ShowAiDistortEndpoint:   showAiDistortEndpoint,
		AddAiDistortEndpoint:    addAiDistortEndpoint,
		DeleteAiDistortEndpoint: deleteAiDistortEndpoint,
	}
}

// ---------------------- ShowAiDistort --------------------
// 请求参数转换为 vo -> pb
func (a *AiAgent) RequestShowAiDistort(_ context.Context, request interface{}) (interface{}, error) {
	rq := request.(*req.ShowVO)
	return &pb.ShowReq{
		Appid:    rq.Appid,
		Where:    rq.Where,
		OrderBy:  rq.OrderBy,
		Page:     rq.Page,
		PageSize: rq.PageSize,
	}, nil
}

// 请求响应转换 pb -> vo
func (a *AiAgent) ResponseShowAiDistort(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.ShowAiDistortRsp)
	var pbAiDistorts = make([]*vo.AiDistort, 0)

	if len(resp.AiDistorts) > 0 {
		for _, aiDistort := range resp.AiDistorts {
			pbAiDistorts = append(pbAiDistorts, &vo.AiDistort{
				Uin:       aiDistort.Uin,
				Appid:     aiDistort.Appid,
				Domain:    aiDistort.Domain,
				Payload:   aiDistort.Payload,
				From:      aiDistort.From,
				Remark:    aiDistort.Remark,
				Status:    uint8(aiDistort.Status),
				CreatedAt: aiDistort.CreatedAt,
			})
		}
	}
	return &res.ShowAiDistortRsp{
		AiDistorts: pbAiDistorts,
		Pagination: vo.Pagination{
			Page:      resp.Pagination.Page,
			PageSize:  resp.Pagination.PageSize,
			Total:     resp.Pagination.Total,
			TotalPage: resp.Pagination.TotalPage,
		},
	}, nil
}

// ---------------------- AddAiDistort -----------------------
// 请求参数转换 vo -> pb
func (a *AiAgent) RequestAddAiDistort(_ context.Context, request interface{}) (interface{}, error) {
	rq := request.(*req.AddAiDistortVO)
	return &pb.AddAiDistortReq{
		Appid:   rq.Appid,
		Domain:  rq.Domain,
		Payload: rq.Payload,
		From:    rq.From,
		Remark:  rq.Remark,
	}, nil
}

// 请求响应转换  pb -> vo
func (a *AiAgent) ResponseAddAiDistort(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.Ack)

	return &res.Ack{
		IsOk: resp.IsOk,
	}, nil
}

// ---------------------- DeleteAiDistort -----------------------
// 请求参数转换 vo -> pb
func (a *AiAgent) RequestDeleteAiDistort(_ context.Context, request interface{}) (interface{}, error) {
	rq := request.(*req.DeleteAiDistortVO)
	return &pb.DeleteAiDistortReq{
		Uin: rq.Uin,
	}, nil
}

// 请求响应转换  pb -> vo
func (a *AiAgent) ResponseDeleteAiDistort(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*res.Ack)

	return &pb.Ack{
		IsOk: resp.IsOk,
	}, nil
}
