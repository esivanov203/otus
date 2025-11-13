package logger

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Conf struct {
	Level string `yaml:"level"` // DEBUG, INFO, WARN, ERROR
	Type  string `yaml:"type"`  // json, console
}

type Logger interface {
	Debug(msg string, fields ...Fields)
	Info(msg string, fields ...Fields)
	Warn(msg string, fields ...Fields)
	Error(msg string, fields ...Fields)
}

type CalLogger struct {
	*zap.Logger
}

type Fields map[string]any

func New(lc Conf) (*CalLogger, error) {
	var zl zapcore.Level

	switch strings.ToLower(lc.Level) {
	case "debug":
		zl = zapcore.DebugLevel
	case "info":
		zl = zapcore.InfoLevel
	case "warn":
		zl = zapcore.WarnLevel
	case "error":
		zl = zapcore.ErrorLevel
	default:
		return nil, fmt.Errorf("invalid log level: %s", lc.Level)
	}

	cfg := zap.Config{
		Encoding:         lc.Type,
		Level:            zap.NewAtomicLevelAt(zl),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}

	l, err := cfg.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}

	return &CalLogger{Logger: l}, nil
}

func (l *CalLogger) toZapFields(fields Fields) []zap.Field {
	zf := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zf = append(zf, zap.Any(k, v))
	}
	return zf
}

func (l *CalLogger) Debug(msg string, fields ...Fields) {
	if len(fields) > 0 {
		l.Logger.Debug(msg, l.toZapFields(fields[0])...)
	} else {
		l.Logger.Debug(msg)
	}
}

func (l *CalLogger) Info(msg string, fields ...Fields) {
	if len(fields) > 0 {
		l.Logger.Info(msg, l.toZapFields(fields[0])...)
	} else {
		l.Logger.Info(msg)
	}
}

func (l *CalLogger) Warn(msg string, fields ...Fields) {
	if len(fields) > 0 {
		l.Logger.Warn(msg, l.toZapFields(fields[0])...)
	} else {
		l.Logger.Warn(msg)
	}
}

func (l *CalLogger) Error(msg string, fields ...Fields) {
	if len(fields) > 0 {
		l.Logger.Error(msg, l.toZapFields(fields[0])...)
	} else {
		l.Logger.Error(msg)
	}
}
