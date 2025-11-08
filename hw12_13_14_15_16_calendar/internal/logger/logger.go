package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
)

type LoggerConf struct {
	Level string `yaml:"level"` // DEBUG, INFO, WARN, ERROR
	Type  string `yaml:"type"`  // json, console
}

type Logger struct {
	*zap.Logger
}

func New(lc LoggerConf) (*Logger, error) {
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

	l, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{Logger: l}, nil
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.Logger.Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.Logger.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.Logger.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.Logger.Error(msg, fields...)
}
