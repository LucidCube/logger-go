package logger

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Format specifies the log output format
type Format int8

// format enumeration
const (
	FormatLines Format = iota
	FormatGoogleCloud
)

type correlationIdType int

const (
	requestIdKey correlationIdType = iota
)

// Default format is format lines
var (
	logger        *zap.Logger
	isInitialized = false

	currentFormat = FormatLines
	debugLogging  = false
	logCtx        context.Context
)

//WithRequestID sets a requestID value in ctx
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIdKey, requestID)
}

//SetContext sets context on logger
func SetContext(ctx context.Context) {
	logCtx = ctx
}

// SetFormat sets the log output format
func SetFormat(format Format) {
	if isInitialized {
		logger.Warn("logger already initialized when setting format")
	}
	currentFormat = format
}

// EnableDebugLogging outputs debug logs and trace logs
func EnableDebugLogging() {
	if isInitialized {
		logger.Warn("logger already initialized when enabling debug")
	}
	debugLogging = true
}

// Instance is the logger instance
func Instance() *zap.Logger {
	if logger == nil {

		logLevel := zap.InfoLevel
		if debugLogging {
			logLevel = zap.DebugLevel
		}

		cfg := zap.Config{
			Development:      true,                           // more liberal stack traces.
			Level:            zap.NewAtomicLevelAt(logLevel), // lowest log level
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
			Sampling: &zap.SamplingConfig{
				Initial:    100,
				Thereafter: 100,
			},
		}

		switch currentFormat {
		case FormatGoogleCloud:
			cfg.Encoding = "json"
			cfg.EncoderConfig = googleEncoderConfig()
		case FormatLines:
			fallthrough
		default:
			cfg.Encoding = "console"
			cfg.EncoderConfig = zap.NewDevelopmentEncoderConfig()
		}
		logger, _ = cfg.Build()
		isInitialized = true
	}

	if logCtx != nil {
		if ctxRequestID, ok := logCtx.Value(requestIdKey).(string); ok {
			return logger.With(zap.String("REQUEST_ID", ctxRequestID))
		}
	}

	return logger
}

// encoder to match GCP payloads
// https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry
func googleEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:       "timestamp",
		LevelKey:      "severity",
		NameKey:       "logName",
		CallerKey:     "caller",
		MessageKey:    "textPayload",
		StacktraceKey: "trace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel: func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
			switch l {
			case zapcore.DebugLevel:
				enc.AppendString("DEBUG")
			case zapcore.InfoLevel:
				enc.AppendString("INFO")
			case zapcore.WarnLevel:
				enc.AppendString("WARNING")
			case zapcore.ErrorLevel:
				enc.AppendString("ERROR")
			case zapcore.DPanicLevel:
				enc.AppendString("CRITICAL")
			case zapcore.PanicLevel:
				enc.AppendString("ALERT")
			case zapcore.FatalLevel:
				enc.AppendString("EMERGENCY")
			}
		},
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}
