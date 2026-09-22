package logging

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Backend string

const (
	BackendZerolog Backend = "zerolog"
	// BackendZap     Backend = "zap"
)

type Options struct {
	Backend  Backend
	Level    string
	Encoding string
	FilePath string
	AppName  string
	Env      string
}

func (o Options) withDefaults() Options {
	if o.Backend == "" {
		o.Backend = BackendZerolog
	}
	if o.Level == "" {
		o.Level = "info"
	}
	if o.Encoding == "" {
		o.Encoding = "json"
	}

	return o
}

func New(opts Options) (Logger, error) {
	opts = opts.withDefaults()

	writer, err := buildWriter(opts)
	if err != nil {
		return nil, err
	}

	switch opts.Backend {
	case BackendZerolog:
		return newZerologLogger(opts, writer)
	// case BackendZap:
	// 	return newZapLogger(opts, writer)
	default:
		return nil, fmt.Errorf("logging: unknown backend %q (allowed: zerolog, zap)", opts.Backend)
	}
}

func MustNew(opts Options) Logger {
	l, err := New(opts)
	if err != nil {
		panic(err)
	}
	return l
}

func buildWriter(opts Options) (io.Writer, error) {
	if opts.FilePath == "" {
		return os.Stdout, nil
	}

	dir := opts.FilePath
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return nil, fmt.Errorf("logging: failed to create log directory %q: %w", dir, err)
	}

	name := filepath.Join(dir, sanitize(opts.AppName)+".log")
	file, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return nil, fmt.Errorf("logging: failed to open log file %q: %w", name, err)
	}

	return io.MultiWriter(os.Stdout, file), nil
}

func sanitize(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, string(filepath.Separator), "-")
	s = strings.ReplaceAll(s, "/", "-")
	if s == "" {
		return "app"
	}
	return s
}