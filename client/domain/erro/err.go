package erro

import "errors"

// 错误类型定义

type ErrorWrapper struct {
	Error string `json:"errors"`
}

// 请求参数错误
var ErrorBadRequest = errors.New("invalid request parameter")
