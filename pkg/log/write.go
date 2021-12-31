package log

import (
	"context"

	"go.uber.org/zap"
)

var SystemKey = "system"
var SystemName = "system"
var SystemTraceName = "system"

//S SugaredLogger with log-id field
func S(ctx context.Context) *zap.SugaredLogger {
	return zap.S().With(TraceIDField(ctx))
}

//Logger with log-id field
func L(ctx context.Context) *zap.Logger {
	return zap.L().With(TraceIDField(ctx))
}

//TraceIDField TraceIDField
func TraceIDField(ctx context.Context) zap.Field {
	traceID := GetTraceIDFromContext(ctx)
	if traceID != "" {
		return zap.String(string(TraceKeyName), traceID)
	}
	return zap.Skip()
}
