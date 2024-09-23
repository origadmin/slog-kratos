package slog

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

var _ log.Logger = (*Logger)(nil)

type Logger struct {
	ctx    context.Context
	log    *slog.Logger
	msgKey string
}

type Option func(*Logger)

func WithContext(ctx context.Context) Option {
	return func(log *Logger) {
		log.ctx = ctx
	}
}

func WithMessageKey(key string) Option {
	return func(log *Logger) {
		log.msgKey = key
	}
}

func NewLogger(logger *slog.Logger, opts ...Option) *Logger {
	slog := &Logger{
		ctx:    context.Background(),
		log:    logger,
		msgKey: log.DefaultMessageKey,
	}

	for _, opt := range opts {
		opt(slog)
	}

	return slog
}

func (l *Logger) Context() context.Context {
	return l.ctx
}

func (l *Logger) Log(level log.Level, keyvals ...interface{}) error {
	slevel := toLevel(level)
	// If logging at this level is completely disabled, skip the overhead of
	// string formatting.
	if !l.log.Enabled(l.ctx, slevel) {
		return nil
	}
	var (
		msg    = ""
		keylen = len(keyvals)
	)
	if keylen == 0 || keylen%2 != 0 {
		l.log.LogAttrs(l.Context(), slog.LevelWarn, fmt.Sprint("Keyvalues must appear in pairs: ", keyvals))
		return nil
	}

	data := make([]slog.Attr, 0, (keylen/2)+1)
	for i := 0; i < keylen; i += 2 {
		key, _ := keyvals[i].(string)
		if key == l.msgKey {
			msg, _ = keyvals[i+1].(string)
			continue
		}
		data = append(data, anyToAttr(key, keyvals[i+1]))
	}
	l.log.LogAttrs(l.Context(), slevel, msg, data...)
	return nil
}

func (l *Logger) Close() error {
	return nil
}

func toLevel(level log.Level) slog.Level {
	switch level {
	case log.LevelDebug:
		return slog.LevelDebug
	case log.LevelInfo:
		return slog.LevelInfo
	case log.LevelWarn:
		return slog.LevelWarn
	case log.LevelError:
		return slog.LevelError
	case log.LevelFatal:
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}

func anyToAttr(key string, value any) slog.Attr {
	switch v := value.(type) {
	case string:
		return slog.String(key, v)
	case int:
		return slog.Int(key, v)
	case int64:
		return slog.Int64(key, v)
	case uint64:
		return slog.Uint64(key, v)
	case float64:
		return slog.Float64(key, v)
	case bool:
		return slog.Bool(key, v)
	case time.Time:
		return slog.Time(key, v)
	case time.Duration:
		return slog.Duration(key, v)
	default:
		return slog.Any(key, v)
	}
}
