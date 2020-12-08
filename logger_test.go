package logger_test

import (
	"context"
	"github.com/lucidcube/logger-go"
	"go.uber.org/zap"
	"testing"
)

//This is more of a usage example than a test
func TestSetContext(t *testing.T) {
	logger.SetFormat(logger.FormatGoogleCloud)

	ctx := context.Background()
	ctx = logger.WithRequestID(ctx, "request-1234")
	logger.SetContext(ctx)

	logger.Instance().Info("debug message", zap.String("key", "key-1"))

	ctx = logger.WithRequestID(ctx, "request-34567")
	logger.SetContext(ctx)
	logger.Instance().Info("another debug message", zap.String("key", "key-2"))
}
