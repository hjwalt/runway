package logger

import (
	"context"
	"log/slog"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapHandler struct {
	config zap.Config
	logger *zap.Logger

	attrs []slog.Attr
	group string
}

func (h *ZapHandler) Enabled(_ context.Context, level slog.Level) bool {
	return int(level) >= int(h.config.Level.Level())
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

	if ce := h.logger.Check(zapcore.Level(record.Level), record.Message); ce != nil {
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
