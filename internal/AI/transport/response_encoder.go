package transport

import (
	"context"
	"net/http"
)

// HTTP响应数据编码 ResponseWriter -> json
type HTTPResponseEncoder interface {
	HTTPResponseEncode(ctx context.Context, w http.ResponseWriter, response interface{}) (err error)
}

// GRPC响应数据编码 vo -> pb
type GRPCResponseEncoder interface {
	GRPCResponseEncode(ctx context.Context, resVo interface{}) (protoBuffRes interface{}, err error)
}
