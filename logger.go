package rtc

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
)

type RuntimeLogger interface {
	Trace(msg string, args ...any)
	TraceContext(ctx context.Context, msg string, args ...any)
	Fatal(msg string, args ...any)
	FatalContext(ctx context.Context, msg string, args ...any)
	With(a ...any) RuntimeLogger
}

func NewRuntimeLogger(l *slog.Logger) RuntimeLogger {
	return &tflogger{l}
}

type tflogger struct {
	*slog.Logger
}

func (t *tflogger) Trace(msg string, args ...any) {
	t.Logger.Log(context.Background(), LevelTrace, msg, args...)
}

func (t *tflogger) TraceContext(ctx context.Context, msg string, args ...any) {
	t.Logger.Log(ctx, LevelTrace, msg, args...)
}

func (t *tflogger) Fatal(msg string, args ...any) {
	t.Logger.Log(context.Background(), LevelFatal, msg, args...)
}

func (t *tflogger) FatalContext(ctx context.Context, msg string, args ...any) {
	t.Logger.Log(ctx, LevelFatal, msg, args...)
}

func (t *tflogger) With(a ...any) RuntimeLogger {
	return &tflogger{t.Logger.With(a...)}
}

func Trace(msg string, args ...any) {
	NewRuntimeLogger(slog.Default()).Trace(msg, args...)
}

func TraceContext(ctx context.Context, msg string, args ...any) {
	NewRuntimeLogger(slog.Default()).TraceContext(ctx, msg, args...)
}

func Fatal(msg string, args ...any) {
	NewRuntimeLogger(slog.Default()).Fatal(msg, args...)
}

func FatalContext(ctx context.Context, msg string, args ...any) {
	NewRuntimeLogger(slog.Default()).FatalContext(ctx, msg, args...)
}

func With(a ...any) RuntimeLogger {
	return &tflogger{slog.With(a...)}
}

func resolveLogLevel() {
	givenLogLevel := flagLogLevel

	if givenLogLevel == "" {
		if flagDebug {
			flagLogLevel = "debug"
		} else if flagTestMode {
			flagLogLevel = "debug"
		} else if flagVerbose {
			flagLogLevel = "info"
		} else {
			flagLogLevel = "error"
		}
	} else {
		_, ok := logLevelMap[strings.ToUpper(givenLogLevel)]
		if !ok {
			fmt.Printf("Unsupported log level: %v\n", givenLogLevel)
			flagLogLevel = slog.LevelError.String()
		}
	}

	level := logLevelMap[strings.ToUpper(flagLogLevel)]
	logLevel.Set(level)
	logOptions.AddSource = level <= slog.LevelDebug
}

func replaceAttributes(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.LevelKey {
		level := a.Value.Any().(slog.Level)
		levelLabel, exists := logLevelNames[level]
		if !exists {
			levelLabel = level.String()
		}

		a.Value = slog.StringValue(levelLabel)
	}

	return a
}
