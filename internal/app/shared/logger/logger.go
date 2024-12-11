package shared

import (
	"context"
	"os"

	"go.uber.org/zap"
)

type Logger struct {
	zapLogger *zap.Logger
}

func NewLogger() *Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	return &Logger{zapLogger: logger}
}
func (l *Logger) Error(ctx context.Context, msg string, err error, context ...map[string]interface{}) {
	var zapFields []zap.Field
	if len(context) > 0 && context[0] != nil {
		for key, value := range context[0] {
			zapFields = append(zapFields, zap.Any(key, value))
		}
	}
	zapFields = append(zapFields, zap.Error(err))
	l.zapLogger.Error(msg, zapFields...)
}

func (l *Logger) Info(ctx context.Context, msg string, context ...map[string]interface{}) {
	var zapFields []zap.Field

	if len(context) > 0 && context[0] != nil {
		for key, value := range context[0] {
			zapFields = append(zapFields, zap.Any(key, value))
		}
	}
	requestId, ok := ctx.Value("requestid").(string)
	if !ok {
		requestId = "unknown"
	}
	l.zapLogger.With(zap.String("requestid", requestId)).Info(msg, zapFields...)
}

func (l *Logger) Warn(ctx context.Context, msg string, context ...map[string]interface{}) {
	var zapFields []zap.Field
	if len(context) > 0 && context[0] != nil {
		for key, value := range context[0] {
			zapFields = append(zapFields, zap.Any(key, value))
		}
	}
	requestId, ok := ctx.Value("requestid").(string)
	if !ok {
		requestId = "unknown"
	}
	l.zapLogger.With(zap.String("requestid", requestId)).Warn(msg, zapFields...)
}

func (l *Logger) Fatal(ctx context.Context, msg string, err error, context ...map[string]interface{}) {
	var zapFields []zap.Field
	if len(context) > 0 && context[0] != nil {
		for key, value := range context[0] {
			zapFields = append(zapFields, zap.Any(key, value))
		}
	}
	requestId, ok := ctx.Value("requestid").(string)
	if !ok {
		requestId = "unknown"
	}
	// Add the error as a field
	zapFields = append(zapFields, zap.Error(err))
	l.zapLogger.With(zap.String("requestid", requestId)).Fatal(msg, zapFields...)

	// Terminate the application
	os.Exit(1)
}
func (l *Logger) Close() {
	_ = l.zapLogger.Sync()
}
