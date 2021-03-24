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

// 获取Ai 误报/漏报 配置请求
type ShowVO struct {
	Appid    uint64 `json:"appid,string"`                                   // 用户ID
	Where    string `json:"where"`                                          // 查找条件 规定 key=value多个条件用&进行间隔
	OrderBy  string `json:"order_by"`                                       // 排序条件
	Page     int64  `json:"page" example:"0"`                               // 当前页
	PageSize int64  `json:"pageSize" example:"20" validate:"gt=0,lte=1000"` // 每页条数大于0 小于1000
	//PageSize int64 `json:"pageSize" example:"20" `
}

// 解析HTTP请求
func (*ShowVO) HTTPRequestDecode(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var vo ShowVO
	if err = json.NewDecoder(r.Body).Decode(&vo); err != nil {
		return nil, erro.ErrorBadRequest
	}
	myLog.GetLogger().Debug(fmt.Sprint(ctx.Value(global.ContextReqUUid)), zap.Any(" 开始解析请求数据", vo))
	if err = ValidatorStruct(vo); err != nil {
		return nil, erro.ErrorBadRequest
	}
	if vo.PageSize <= 0 || vo.PageSize >= 1000 { // 设置默认分页条数
		vo.PageSize = 20
	}
	return &vo, nil
}

// 解析GRPC请求
func (*ShowVO) GRPCServerRequestDecode(ctx context.Context, grpcReq interface{}) (request interface{}, err error) {
	// 数据类型转换
	re := grpcReq.(*pb.ShowReq)
	var vo = ShowVO{
		Appid:    re.Appid,
		Where:    re.Where,
		OrderBy:  re.OrderBy,
		Page:     re.Page,
		PageSize: re.PageSize,
	}

	myLog.GetLogger().Debug(fmt.Sprint(ctx.Value(global.ContextReqUUid)), zap.Any(" 开始解析请求数据", vo))

	if err = ValidatorStruct(vo); err != nil { // 校验参数
		return nil, erro.ErrorBadRequest
	}
	if vo.PageSize <= 0 || vo.PageSize >= 1000 { // 设置默认分页条数
		vo.PageSize = 20
	}
	return &vo, nil
}
