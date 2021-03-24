package res

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"micro-service-tmpl/client/domain/global"
	"micro-service-tmpl/client/domain/pb"
	"micro-service-tmpl/utils/myLog"
	"net/http"
)

type Ack struct {
	IsOk bool `json:"isOk"` // 成功与否
}

// HTTP响应数据转换
func (*Ack) HTTPResponseEncode(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	myLog.GetLogger().Debug(fmt.Sprint(ctx.Value(global.ContextReqUUid)), zap.Any("请求结束封装返回值", response))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// GRPC响应数据转换
func (*Ack) GRPCResponseEncode(ctx context.Context, resVo interface{}) (protoBuffRes interface{}, err error) {
	resp := resVo.(*Ack)

	return &pb.Ack{IsOk: resp.IsOk}, nil
}
