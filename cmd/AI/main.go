package main

import (
	"context"
	"flag"
	"fmt"
	log2 "github.com/go-kit/kit/log"
	metricsprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/kit/sd/etcdv3"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"hash/crc32"
	"micro-service-tmpl/internal/AI/dao"
	"micro-service-tmpl/internal/AI/domain/pb"
	"micro-service-tmpl/internal/AI/endpoint"
	"micro-service-tmpl/internal/AI/service"
	"micro-service-tmpl/internal/AI/transport"
	"micro-service-tmpl/utils/jaegerTracer"
	"micro-service-tmpl/utils/myLog"
	"micro-service-tmpl/utils/viper"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

type ServerConfig struct {
	ServiceName    string   // 服务名
	Evn            string   // 环境
	GrpcAddr       string   // grpc地址
	HttpAddr       string   // http地址
	PrometheusAddr string   // Prometheus服务监控地址
	JaegerAddr     string   // Jaeger链路追踪上报地址
	EtcdAddr       []string // etcd地址
}

// 退出信号管道
var quitChan = make(chan error, 1)

func main() {
	//go run ./cmd/Ai/main.go -g 127.0.0.1:30001 -h 127.0.0.1:30002 -p 127.0.0.1:30003
	// 日志
	logger := myLog.GetLogger()

	var (
		ttl        = 5 * time.Second
		err        error
		srv        service.AIService
		serverConf ServerConfig
	)
	// 获取server配置
	if err = viper.ViperConf.UnmarshalKey("server", &serverConf); err != nil {
		myLog.GetLogger().Fatal("数据库获取配置文件失败" + err.Error())
	}
	var (
		// grpc 监听地址
		grpcAddr = flag.String("g", serverConf.GrpcAddr, "grpcAddr")

		// http 监听地址
		httpAddr = flag.String("h", serverConf.HttpAddr, "httpAddr")

		// prometheus 监听地址
		prometheusAddr = flag.String("p", serverConf.PrometheusAddr, "prometheus addr")
	)

	flag.Parse()
	// ETCD 配置
	options := etcdv3.ClientOptions{
		DialTimeout:   ttl,
		DialKeepAlive: ttl,
	}
	// 获取ETCD客户端
	etcdClient, err := etcdv3.NewClient(context.Background(), serverConf.EtcdAddr, options)
	if err != nil {
		logger.Error("[user_agent]  NewClient", zap.Error(err))
		return
	}

	// 服务注册
	etcdRegister := etcdv3.NewRegistrar(etcdClient, etcdv3.Service{
		Key:   fmt.Sprintf("%s/%d", serverConf.ServiceName, crc32.ChecksumIEEE([]byte(*grpcAddr))),
		Value: *grpcAddr,
	}, log2.NewNopLogger())

	count := metricsprometheus.NewCounterFrom(prometheus.CounterOpts{
		Subsystem: "AiService",
		Name:      "request_count",
		Help:      "Number of requests",
	}, []string{"method"})

	histogram := metricsprometheus.NewHistogramFrom(prometheus.HistogramOpts{
		Subsystem: "AiService",
		Name:      "request_consume",
		Help:      "Request consumes time",
	}, []string{"method"})

	// jaeger链路追踪
	tracer, _, err := jaegerTracer.NewJaegerTracer(serverConf.ServiceName, serverConf.JaegerAddr)
	if err != nil {
		logger.Warn(fmt.Sprint(zap.Any("jaeger tracer init failed:", err)))
		srv = service.NewAIServiceImpl(dao.AiDistort{}, count, histogram, logger, nil)
	} else {
		srv = service.NewAIServiceImpl(dao.AiDistort{}, count, histogram, logger, tracer)
	}

	// 令牌桶服务限流 每秒产生 cpu个数令牌，存储10 * cpu个数个令牌
	golangLimit := rate.NewLimiter(rate.Limit(runtime.NumCPU()), 10*runtime.NumCPU())
	aiEndpoints := endpoint.NewAiEndpoints(srv, logger, golangLimit)

	// 开启协程监听http请求
	go func() {
		httpServer := transport.MakeHttpHandler(aiEndpoints, logger)
		logger.Info("[Ai_Server] http run " + *httpAddr)
		quitChan <- http.ListenAndServe(*httpAddr, httpServer)
	}()

	// 开启协程监听grpc请求
	go func() {
		grpcServer := transport.NewGRPCServer(aiEndpoints, logger)
		grpcListener, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			logger.Warn("[Ai_Server] Listen", zap.Error(err))
			quitChan <- err
			return
		}

		chainUnaryServer := grpcmiddleware.ChainUnaryServer(
			grpctransport.Interceptor,
			grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(tracer)),
			grpc_zap.UnaryServerInterceptor(myLog.GetLogger()),
			jaegerTracer.JaegerServerMiddleware(tracer),
		)
		baseServer := grpc.NewServer(grpc.UnaryInterceptor(chainUnaryServer))
		pb.RegisterAiServer(baseServer, grpcServer)
		//服务注册
		etcdRegister.Register()
		logger.Info("[Ai_Server] grpc run " + *grpcAddr)
		quitChan <- baseServer.Serve(grpcListener)
	}()

	// 服务监控
	go func() {
		logger.Info("[Ai_Server] prometheus run " + *prometheusAddr)
		m := http.NewServeMux()
		m.Handle("/metrics", promhttp.Handler())
		quitChan <- http.ListenAndServe(*prometheusAddr, m)
	}()

	// 开启协程监听退出
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		quitChan <- fmt.Errorf("%s", <-c)
	}()

	// 主协程阻塞等待
	err = <-quitChan
	logger.Debug("err: ", zap.Any("error", err))
}
