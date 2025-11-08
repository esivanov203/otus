package logger

import (
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"testing"
)

func TestLoggerLevelType(t *testing.T) {
	t.Run("valid level and type", func(t *testing.T) {
		_, err := New(LoggerConf{Level: "debug", Type: "json"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	tests := []struct {
		name string
		conf LoggerConf
	}{
		{"invalid level", LoggerConf{Level: "INVALID", Type: "json"}},
		{"invalid type", LoggerConf{Level: "info", Type: "text"}},
	}

	for _, tt := range tests {
		tt := tt // захват переменной для t.Run
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.conf)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
		})
	}
}

func TestLoggerMethods(t *testing.T) {
	core, logs := observer.New(zapcore.DebugLevel)
	l := &Logger{Logger: zap.New(core)}

	l.Debug("debug message")
	l.Info("info message")
	l.Warn("warn message")
	l.Error("error message")

	entries := logs.All()
	require.Len(t, entries, 4, "expected 4 log entries")

	cases := []struct {
		i     int
		level zapcore.Level
		msg   string
	}{
		{0, zapcore.DebugLevel, "debug message"},
		{1, zapcore.InfoLevel, "info message"},
		{2, zapcore.WarnLevel, "warn message"},
		{3, zapcore.ErrorLevel, "error message"},
	}
	for _, tt := range cases {
		require.Equal(t, tt.level, entries[tt.i].Level)
		require.Equal(t, tt.msg, entries[tt.i].Message)
	}
}
