package logger_test

import (
	"context"
	"github.com/lucidcube/logger-go"
	"go.uber.org/zap"
	"testing"
)

//This is more of a usage example than a test
func TestSetContext(t *testing.T) {
	ctx := context.Background()
	logger.SetContext(ctx, "request-1234")
	logger.SetFormat(logger.FormatGoogleCloud)

	logger.Instance().Info("debug message", zap.String("key", "key-1"))
}
