// Copyright (c) 2024 OrigAdmin. All rights reserved.

// Package slog provides a structured logging system.
package slog

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/goexts/ggb/settings"
)

// Logger represents a structured logger.
type Logger struct {
	// ctx is the context associated with the logger.
	ctx context.Context
	// msgKey is the key used to identify the output message.
	msgKey string
	// output is the underlying slog.Logger instance.
	output *slog.Logger
}

// Option is a function that configures a Logger instance.
type Option = func(log *Logger)

// WithContext returns an Option that sets the context for the logger.
func WithContext(ctx context.Context) Option {
	return func(log *Logger) {
		log.ctx = ctx
	}
}

// WithMessageKey returns an Option that sets the message key for the logger.
func WithMessageKey(key string) Option {
	return func(log *Logger) {
		log.msgKey = key
	}
}

// WithLogger returns an Option that sets the underlying slog.Logger instance.
func WithLogger(logger *slog.Logger) Option {
	return func(log *Logger) {
		log.output = logger
	}
}

// NewLogger returns a new Logger instance with the given options.
func NewLogger(opts ...Option) *Logger {
	// Create a new Logger instance with default values.
	// Apply the given options to the logger.
	logger := settings.Apply(&Logger{
		ctx:    context.Background(),
		msgKey: log.DefaultMessageKey,
	}, opts)

	// Set the default slog.Logger instance if none is provided.
	if logger.output == nil {
		logger.output = slog.Default()
	}

	return logger
}

// Context returns the context associated with the logger.
func (l *Logger) Context() context.Context {
	return l.ctx
}

// Log logs a message at the given level with the given key-value pairs.
func (l *Logger) Log(level log.Level, keyvals ...interface{}) error {
	// Convert the log level to a slog.Level.
	slevel := toLevel(level)

	// Check if logging at this level is enabled.
	if !l.output.Enabled(l.ctx, slevel) {
		return nil
	}
	var (
		msg    = ""
		keylen = len(keyvals)
	)
	// Check if the key-value pairs are valid (i.e., they appear in pairs).
	if keylen == 0 || keylen%2 != 0 {
		// Log a warning message if the key-value pairs are invalid.
		l.output.LogAttrs(l.Context(), slog.LevelWarn, fmt.Sprint("Key and values must appear in pairs: ", keyvals))
		return nil
	}
	// Create a slice to store the log attributes.
	data := make([]slog.Attr, 0, (keylen/2)+1)
	// Iterate over the key-value pairs and create log attributes.
	for i := 0; i < keylen; i += 2 {
		key, _ := keyvals[i].(string)
		if key == l.msgKey {
			// Extract the log message if the key matches the message key.
			msg, _ = keyvals[i+1].(string)
			continue
		}
		// Create a log attribute for the key-value pair.
		data = append(data, anyToAttr(key, keyvals[i+1]))
	}
	// Log the message with the given attributes.
	l.output.LogAttrs(l.Context(), slevel, msg, data...)
	return nil
}

// Close closes the logger.
func (l *Logger) Close() error {
	return nil
}

// toLevel converts a log.Level to an slog.Level.
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

// anyToAttr converts a value to an slog.Attr.
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

// Ensure that Logger implements the output.Logger interface.
var _ log.Logger = (*Logger)(nil)
