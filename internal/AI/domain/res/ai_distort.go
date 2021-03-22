package res

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"micro-service-tmpl/internal/AI/dao"
	"micro-service-tmpl/internal/AI/domain/global"
	"micro-service-tmpl/internal/AI/domain/pb"
	"micro-service-tmpl/internal/AI/domain/vo"
	"micro-service-tmpl/utils/log"
	"net/http"
)

// 获取Ai误报信息响应
type ShowAiDistortRsp struct {
	AiDistorts []*vo.AiDistort `json:"aiDistorts"`
	Pagination vo.Pagination   `json:"pagination"`
}

// HTTP响应数据转换
func (*ShowAiDistortRsp) HTTPResponseEncode(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	log.GetLogger().Debug(fmt.Sprint(ctx.Value(global.ContextReqUUid)), zap.Any("请求结束封装返回值", response))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// GRPC响应数据转换
func (*ShowAiDistortRsp) GRPCResponseEncode(ctx context.Context, resVo interface{}) (protoBuffRes interface{}, err error) {
	resp := resVo.(*ShowAiDistortRsp)
	pbAiDistorts := make([]*pb.AiDistortRsp, 0)
	if len(resp.AiDistorts) > 0 {
		for _, distort := range resp.AiDistorts {
			pbAiDistorts = append(pbAiDistorts, &pb.AiDistortRsp{
				Uin:       distort.Uin,
				Appid:     distort.Appid,
				Domain:    distort.Domain,
				Payload:   distort.Payload,
				From:      distort.From,
				Remark:    distort.Remark,
				Status:    uint32(distort.Status),
				CreatedAt: distort.CreatedAt,
			})
		}
	}

	return &pb.ShowAiDistortRsp{
		AiDistorts: pbAiDistorts,
		Pagination: &pb.PaginationRsp{
			Page:      resp.Pagination.Page,
			PageSize:  resp.Pagination.PageSize,
			Total:     resp.Pagination.Total,
			TotalPage: resp.Pagination.TotalPage,
		},
	}, nil
}

// BuildAiDistort Ai误报信息
func BuildAiDistort(item dao.AiDistort) vo.AiDistort {
	return vo.AiDistort{
		Uin:       item.Uin,
		Appid:     item.Appid,
		Domain:    item.Domain,
		Payload:   item.Payload,
		From:      item.From,
		Remark:    item.Remark,
		Status:    item.Status,
		CreatedAt: item.CreatedAt.Unix(),
	}
}

// BuildCarousels Ai误报信息列表
func BuildDistorts(items []dao.AiDistort) (distorts []*vo.AiDistort) {
	for _, item := range items {
		distort := BuildAiDistort(item)
		distorts = append(distorts, &distort)
	}
	return distorts
}
