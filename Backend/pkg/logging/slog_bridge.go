package logging

import (
	"context"
	"log/slog"
)

type slogHandler struct {
	logger Logger
	attrs  []Field
	group  string
}

func NewSlogHandler(l Logger) slog.Handler {
	return &slogHandler{logger: l}
}

func SetupSlogBridge(l Logger) {
	slog.SetDefault(slog.New(NewSlogHandler(l)))
}

func (h *slogHandler) Enabled(_ context.Context, _ slog.Level) bool {

	return true
}

func (h *slogHandler) Handle(ctx context.Context, r slog.Record) error {
	fields := make([]Field, 0, len(h.attrs)+r.NumAttrs())
	fields = append(fields, h.attrs...)

	r.Attrs(func(a slog.Attr) bool {
		fields = append(fields, Field{Key: h.qualify(a.Key), Value: a.Value.Any()})
		return true
	})

	switch {
	case r.Level >= slog.LevelError:
		h.logger.Error(ctx, r.Message, fields...)
	case r.Level >= slog.LevelWarn:
		h.logger.Warn(ctx, r.Message, fields...)
	case r.Level >= slog.LevelInfo:
		h.logger.Info(ctx, r.Message, fields...)
	default:
		h.logger.Debug(ctx, r.Message, fields...)
	}
	return nil
}

func (h *slogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	next := &slogHandler{
		logger: h.logger,
		group:  h.group,
		attrs:  make([]Field, 0, len(h.attrs)+len(attrs)),
	}
	next.attrs = append(next.attrs, h.attrs...)
	for _, a := range attrs {
		next.attrs = append(next.attrs, Field{Key: h.qualify(a.Key), Value: a.Value.Any()})
	}
	return next
}

func (h *slogHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	group := name
	if h.group != "" {
		group = h.group + "." + name
	}
	return &slogHandler{logger: h.logger, attrs: h.attrs, group: group}
}

func (h *slogHandler) qualify(key string) string {
	if h.group == "" {
		return key
	}
	return h.group + "." + key
}
