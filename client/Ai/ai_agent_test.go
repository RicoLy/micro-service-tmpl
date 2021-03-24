package Ai

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"micro-service-tmpl/client/domain/global"
	"micro-service-tmpl/client/domain/pb"
	"micro-service-tmpl/client/domain/req"
	"os"
	"testing"
)

func TestNewAiAgentClient(t *testing.T) {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
	client, err := NewAiAgentClient([]string{"127.0.0.1:2379"}, logger, "127.0.0.1:6831")
	if err != nil {
		t.Error(err)
		return
	}
	for i := 0; i < 5; i++ {
		aiAgent, err := client.AiAgentClient()
		if err != nil {
			t.Error(err)
			return
		}

		ack, err := aiAgent.AddAiDistort(context.Background(), &req.AddAiDistortVO{
			Appid:   222,
			Domain:  "/login",
			Payload: "/fdaf*",
			From:    "app",
			Remark:  "app",
		})
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(ack.IsOk)
	}
}

func TestGRPC(t *testing.T) {
	serviceAddress := "127.0.0.1:8881"
	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure())
	if err != nil {
		panic("connect error")
	}
	defer conn.Close()
	aiClient := pb.NewAiClient(conn)
	UUID := uuid.NewV5(uuid.Must(uuid.NewV4(), nil), "req_uuid").String()
	md := metadata.Pairs(global.ContextReqUUid, UUID)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	for i := 0; i < 20; i++ {
		res, err := aiClient.AddAiDistort(ctx, &pb.AddAiDistortReq{
			Appid:   123,
			Domain:  "/login",
			Payload: "/fdaf*",
			From:    "app",
			Remark:  "app",
		})
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(res.IsOk)
		}
	}
}