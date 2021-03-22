package req

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"micro-service-tmpl/internal/AI/domain/erro"
	"micro-service-tmpl/internal/AI/domain/global"
	"micro-service-tmpl/internal/AI/domain/pb"
	"micro-service-tmpl/utils/log"
	"net/http"
)

// 删除Ai误报信息
type DeleteAiDistortVO struct {
	Uin uint64 `json:"uin,string"` // Uin
}

// 解析HTTP获取删除Ai误报信息请求
func (*DeleteAiDistortVO) HTTPRequestDecode(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var vo DeleteAiDistortVO
	if err = json.NewDecoder(r.Body).Decode(&vo); err != nil {
		return nil, erro.ErrorBadRequest
	}
	log.GetLogger().Debug(fmt.Sprint(ctx.Value(global.ContextReqUUid)), zap.Any(" 开始解析请求数据", vo))
	if err = ValidatorStruct(vo); err != nil { // 校验参数
		return nil, erro.ErrorBadRequest
	}
	return &vo, nil
}

// 解析GRPC请求
func (*DeleteAiDistortVO) GRPCServerRequestDecode(ctx context.Context, grpcReq interface{}) (request interface{}, err error) {
	// 数据类型转换
	re := grpcReq.(*pb.DeleteAiDistortReq)
	var vo = DeleteAiDistortVO{Uin: re.Uin}
	log.GetLogger().Debug(fmt.Sprint(ctx.Value(global.ContextReqUUid)), zap.Any(" 开始解析请求数据", vo))
	if err = ValidatorStruct(vo); err != nil { // 校验参数
		return nil, erro.ErrorBadRequest
	}
	return &vo, nil
}
