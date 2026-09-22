package logging

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/rs/zerolog"
)


type zerologLogger struct {
	logger zerolog.Logger
}

var zerologLevels = map[string]zerolog.Level{
	"debug": zerolog.DebugLevel,
	"info":  zerolog.InfoLevel,
	"warn":  zerolog.WarnLevel,
	"error": zerolog.ErrorLevel,
	"fatal": zerolog.FatalLevel,
}

func newZerologLogger(opts Options, w io.Writer) (Logger, error) {
	level, ok := zerologLevels[opts.Level]
	if !ok {
		return nil, fmt.Errorf("logging: سطح لاگ ناشناخته %q", opts.Level)
	}

	if opts.Encoding == "console" {
		w = zerolog.ConsoleWriter{Out: w, TimeFormat: time.RFC3339}
	}

	zerolog.TimeFieldFormat = time.RFC3339Nano

	l := zerolog.New(w).
		Level(level).
		With().
		Timestamp().
		Str("app", opts.AppName).
		Str("env", opts.Env).
		Logger()

	return &zerologLogger{logger: l}, nil
}

func (l *zerologLogger) log(ctx context.Context, ev *zerolog.Event, msg string, fields []Field) {
	if ev == nil {
		return
	}
	for _, f := range fieldsFromContext(ctx) {
		ev = ev.Interface(f.Key, f.Value)
	}
	for _, f := range fields {
		ev = ev.Interface(f.Key, f.Value)
	}
	ev.Msg(msg)
}

func (l *zerologLogger) Debug(ctx context.Context, msg string, fields ...Field) {
	l.log(ctx, l.logger.Debug(), msg, fields)
}

func (l *zerologLogger) Info(ctx context.Context, msg string, fields ...Field) {
	l.log(ctx, l.logger.Info(), msg, fields)
}

func (l *zerologLogger) Warn(ctx context.Context, msg string, fields ...Field) {
	l.log(ctx, l.logger.Warn(), msg, fields)
}

func (l *zerologLogger) Error(ctx context.Context, msg string, fields ...Field) {
	l.log(ctx, l.logger.Error(), msg, fields)
}

func (l *zerologLogger) Fatal(ctx context.Context, msg string, fields ...Field) {
	l.log(ctx, l.logger.Fatal(), msg, fields)
}

func (l *zerologLogger) With(fields ...Field) Logger {
	c := l.logger.With()
	for _, f := range fields {
		c = c.Interface(f.Key, f.Value)
	}
	return &zerologLogger{logger: c.Logger()}
}
