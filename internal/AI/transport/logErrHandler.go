package transport

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"micro-service-tmpl/internal/AI/domain/global"
)

type LogErrorHandler struct {
	logger *zap.Logger
}

func NewZapLogErrorHandler(logger *zap.Logger) *LogErrorHandler {
	return &LogErrorHandler{
		logger: logger,
	}
}

func (h *LogErrorHandler) Handle(ctx context.Context, err error) {
	h.logger.Warn(fmt.Sprint(ctx.Value(global.ContextReqUUid)), zap.Error(err))
}
