package logger

import (
	"context"
	"log/slog"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var levelMap = map[slog.Level]zapcore.Level{
	slog.LevelDebug: zapcore.DebugLevel,
	slog.LevelInfo:  zapcore.InfoLevel,
	slog.LevelWarn:  zapcore.WarnLevel,
	slog.LevelError: zapcore.ErrorLevel,
}

type ZapHandler struct {
	config zap.Config
	logger *zap.Logger

	attrs []slog.Attr
	group string
}

func (h *ZapHandler) Enabled(_ context.Context, level slog.Level) bool {
	return levelMap[level] >= h.config.Level.Level()
}

func (h *ZapHandler) Handle(ctx context.Context, record slog.Record) error {
	fields := []zapcore.Field{}

	for _, attr := range h.attrs {
		fields = append(fields, zap.Any(attr.Key, attr.Value))
	}

	record.Attrs(func(a slog.Attr) bool {
		fields = append(fields, zap.Any(a.Key, a.Value))
		return true
	})

	if ce := h.logger.Check(levelMap[record.Level], record.Message); ce != nil {
		ce.Write(fields...)
	}

	return nil
}

func (h *ZapHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	groupPrepend := ""
	if len(h.group) > 0 {
		groupPrepend = h.group + "."
	}

	newattrs := append([]slog.Attr{}, h.attrs...)
	for _, attr := range attrs {
		newattrs = append(newattrs, slog.Attr{
			Key:   groupPrepend + attr.Key,
			Value: attr.Value,
		})
	}

	return &ZapHandler{
		config: h.config,
		logger: h.logger,

		attrs: newattrs,
		group: h.group,
	}
}

func (h *ZapHandler) WithGroup(name string) slog.Handler {
	return &ZapHandler{
		config: h.config,
		logger: h.logger,

		attrs: h.attrs,
		group: name,
	}
}
