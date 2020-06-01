package logger

import (
	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
)

// Format specifies the log output format
type Format int8

// format enumeration
const (
	FormatLines Format = iota
	FormatGoogleCloud
)

// Default format is format lines
var (
	currentFormat = FormatLines
	debugLogging  = false
)

// SetFormat sets the log output format
func SetFormat(format Format) {
	currentFormat = format
}

// EnableDebugLogging outputs debug logs and trace logs
func EnableDebugLogging(exp trace.Exporter) {
	debugLogging = true
	trace.WithSampler(trace.AlwaysSample())
	trace.RegisterExporter(exp)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
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
