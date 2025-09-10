package slog

import (
	"io"
	"log/slog"
	"os"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
)

var (
	benchLogger   *Logger
	discardLogger = NewLogger(
		WithLogger(slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))),
	)
)

func BenchmarkLogger_Log(b *testing.B) {
	// Setup test cases
	tests := []struct {
		name     string
		level    log.Level
		keyvals  []interface{}
		setup    func() *Logger
		teardown func()
	}{
		{
			name:  "simple debug log",
			level: log.LevelDebug,
			keyvals: []interface{}{
				"msg", "test message",
				"key1", "value1",
				"key2", 42,
			},
			setup: func() *Logger {
				return discardLogger
			},
		},
		{
			name:  "complex log with multiple fields",
			level: log.LevelInfo,
			keyvals: []interface{}{
				"msg", "complex log message",
				"user_id", 12345,
				"email", "test@example.com",
				"active", true,
				"score", 98.76,
			},
			setup: func() *Logger {
				return discardLogger
			},
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			logger := tt.setup()
			if tt.teardown != nil {
				defer tt.teardown()
			}

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_ = logger.Log(tt.level, tt.keyvals...)
			}
		})
	}
}

func BenchmarkLogger_Concurrent(b *testing.B) {
	logger := discardLogger

	b.Run("concurrent logging", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = logger.Log(
					log.LevelInfo,
					"msg", "concurrent log message",
					"goroutine", "test",
					"count", 1,
				)
			}
		})
	})
}

func TestMain(m *testing.M) {
	// Setup benchmark logger
	handler := slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	benchLogger = NewLogger(WithLogger(slog.New(handler)))

	// Run tests
	code := m.Run()
	os.Exit(code)
}
