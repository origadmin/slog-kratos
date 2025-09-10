package slog

import (
	"log/slog"
	"strings"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
)

type testWriteSyncer struct {
	output []string
}

func (x *testWriteSyncer) Write(p []byte) (n int, err error) {
	x.output = append(x.output, string(p))
	return len(p), nil
}

func (x *testWriteSyncer) Sync() error {
	return nil
}

func TestLogger(t *testing.T) {
	syncer := &testWriteSyncer{}
	handler := slog.NewJSONHandler(syncer, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == "time" {
				return slog.String("time", "2023-01-01 00:00:00")
			}
			if a.Key == "level" {
				return slog.String("level", strings.ToLower(a.Value.String()))
			}
			return a
		},
	})

	logger := NewLogger(WithLogger(slog.New(handler)))

	defer func() { _ = logger.Close() }()
	helper := log.NewHelper(logger)
	helper.Debugw("output", "debug")
	helper.Infow("output", "info")
	helper.Warnw("output", "warn")
	helper.Errorw("output", "error")
	helper.Errorw("output", "error", "except warn")
	helper.Info("hello world")

	except := []string{
		"{\"time\":\"2023-01-01 00:00:00\",\"level\":\"debug\",\"msg\":\"\",\"output\":\"debug\"}\n",
		"{\"time\":\"2023-01-01 00:00:00\",\"level\":\"info\",\"msg\":\"\",\"output\":\"info\"}\n",
		"{\"time\":\"2023-01-01 00:00:00\",\"level\":\"warn\",\"msg\":\"\",\"output\":\"warn\"}\n",
		"{\"time\":\"2023-01-01 00:00:00\",\"level\":\"error\",\"msg\":\"\",\"output\":\"error\"}\n",
		"{\"time\":\"2023-01-01 00:00:00\",\"level\":\"warn\",\"msg\":\"Key and values must appear in pairs: [output error except warn]\"}\n",
		"{\"time\":\"2023-01-01 00:00:00\",\"level\":\"info\",\"msg\":\"hello world\"}\n", // not {"level":"info","msg":"","msg":"hello world"}
	}
	for i, s := range except {
		if s != syncer.output[i] {
			t.Logf("except=%s, got=%s", s, syncer.output[i])
			t.Fail()
		}
	}
}
