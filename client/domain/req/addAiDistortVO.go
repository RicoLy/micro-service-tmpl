package req

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"micro-service-tmpl/client/domain/erro"
	"micro-service-tmpl/client/domain/global"
	"micro-service-tmpl/client/domain/pb"
	"micro-service-tmpl/utils/myLog"
	"net/http"
)

// 添加Ai误报信息
type AddAiDistortVO struct {
	Appid   uint64 `json:"appid,string"` // 用户ID
	Domain  string `json:"domain"`       // 域名
	Payload string `json:"Payload"`      // 载荷
	From    string `json:"from"`         // 来源
	Remark  string `json:"remark"`       // 备注
}

// 解析HTTP获取添加Ai误报信息请求
func (*AddAiDistortVO) HTTPRequestDecode(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var vo AddAiDistortVO
	if err = json.NewDecoder(r.Body).Decode(&vo); err != nil {
		return nil, erro.ErrorBadRequest
	}
	myLog.GetLogger().Debug(fmt.Sprint(ctx.Value(global.ContextReqUUid)), zap.Any(" 开始解析请求数据", vo))
	if err = ValidatorStruct(vo); err != nil {
		myLog.GetLogger().Debug(err.Error())
		return nil, erro.ErrorBadRequest
	}
	return &vo, nil
}

// 服务端解析GRPC请求 pb -> vo
func (*AddAiDistortVO) GRPCServerRequestDecode(ctx context.Context, grpcReq interface{}) (request interface{}, err error) {
	// 数据类型转换
	pbReq := grpcReq.(*pb.AddAiDistortReq)
	var vo = AddAiDistortVO{
		Appid:   pbReq.Appid,
		Domain:  pbReq.Domain,
		Payload: pbReq.Payload,
		From:    pbReq.From,
		Remark:  pbReq.Remark,
	}
	myLog.GetLogger().Debug(fmt.Sprint(ctx.Value(global.ContextReqUUid)), zap.Any(" 开始解析请求数据", vo))
	if err = ValidatorStruct(vo); err != nil { // 校验参数
		return nil, erro.ErrorBadRequest
	}
	return &vo, nil
}
