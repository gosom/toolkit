package logger

import (
	"context"
	"os"
	"sync"

	"log/slog"

	"toolkit/pkg/errorsext"
)

type Logger interface {
	Info(ctx context.Context, msg string, args ...any)
	Debug(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)
	With(args ...any) Logger
}

type defaultLogger struct {
	l *slog.Logger
}

func (l *defaultLogger) Info(ctx context.Context, msg string, args ...any) {
	log(ctx, l.l, slog.LevelInfo, msg, args...)
}

func (l *defaultLogger) Debug(ctx context.Context, msg string, args ...any) {
	log(ctx, l.l, slog.LevelDebug, msg, args...)
}

func (l *defaultLogger) Error(ctx context.Context, msg string, args ...any) {
	log(ctx, l.l, slog.LevelError, msg, args...)
}

func (l *defaultLogger) With(args ...any) Logger {
	return &defaultLogger{l: l.l.With(args...)}
}

func Default() Logger {
	once.Do(func() {
		stdlog = LoggerInstance()
		stdreporter = Reporter()
	})

	return &defaultLogger{l: stdlog}
}

var LoggerInstance = func() *slog.Logger {
	ans := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	return ans
}

var Reporter = func() ErrorReporter {
	return &StubErrorReporter{}
}

func Info(ctx context.Context, msg string, args ...any) {
	log(ctx, nil, slog.LevelInfo, msg, args...)
}

func Debug(ctx context.Context, msg string, args ...any) {
	log(ctx, nil, slog.LevelDebug, msg, args...)
}

func Error(ctx context.Context, msg string, args ...any) {
	log(ctx, nil, slog.LevelError, msg, args...)
}

func ContextWithData(ctx context.Context, data ...any) context.Context {
	if len(data) == 0 {
		return ctx
	}

	current := ContextData(ctx)
	current = append(current, data...)

	ctx = context.WithValue(ctx, ctxDataKey, current)

	return ctx
}

func ContextData(ctx context.Context) []any {
	if data, ok := ctx.Value(ctxDataKey).([]any); ok {
		return data
	}

	return nil
}

func ReportError(ctx context.Context, args ...any) {
	once.Do(func() {
		stdlog = LoggerInstance()
		stdreporter = Reporter()
	})

	stdreporter.ReportError(ctx, args...)
}

func log(
	ctx context.Context,
	instance *slog.Logger,
	level slog.Level,
	msg string,
	args ...any,
) {
	once.Do(func() {
		stdlog = LoggerInstance()
		stdreporter = Reporter()
	})

	if instance == nil {
		instance = stdlog
	}

	fromctx := ContextData(ctx)

	args = append(fromctx, args...)

	var (
		attrs []slog.Attr
		key   string
		ok    bool
		val   slog.Value
	)

	for i := 0; i < len(args); i += 2 {
		key, ok = args[i].(string)
		if !ok {
			key = "invalid_key"
		}

		if i+1 < len(args) {
			v := args[i+1]
			if ste, ok := v.(errorsext.StackTracer); ok {
				attr := slog.Attr{Key: "stacktrace", Value: slog.AnyValue(ste.Stacktrace())}
				attrs = append(attrs, attr)
			}

			val = slog.AnyValue(args[i+1])
		} else {
			val = slog.AnyValue(nil)
		}

		attrs = append(attrs, slog.Attr{Key: key, Value: val})
	}

	instance.LogAttrs(ctx, level, msg, attrs...)
}

type ctxKey string

var (
	once        sync.Once
	stdlog      *slog.Logger
	stdreporter ErrorReporter
	ctxDataKey  = ctxKey("log_data")
)
