package rtc

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"strings"
	"time"

	"github.com/jwmwalrus/bnp/onerror"
)

type RuntimeLogger interface {
	Trace(msg string, args ...any)
	TraceContext(ctx context.Context, msg string, args ...any)
	Fatal(msg string, args ...any)
	FatalContext(ctx context.Context, msg string, args ...any)
	With(a ...any) RuntimeLogger
}

func NewRuntimeLogger(l *slog.Logger) RuntimeLogger {
	return &tflogger{Logger: l}
}

type tflogger struct {
	*slog.Logger
	pc uintptr
}

func (t *tflogger) Trace(msg string, args ...any) {
	t.addRecord(context.Background(), LevelTrace, msg, args...)
}

func (t *tflogger) TraceContext(ctx context.Context, msg string, args ...any) {
	t.addRecord(ctx, LevelTrace, msg, args...)
}

func (t *tflogger) Fatal(msg string, args ...any) {
	t.addRecord(context.Background(), LevelFatal, msg, args...)
}

func (t *tflogger) FatalContext(ctx context.Context, msg string, args ...any) {
	t.addRecord(ctx, LevelFatal, msg, args...)
}

func (t *tflogger) With(a ...any) RuntimeLogger {
	return &tflogger{Logger: t.Logger.With(a...)}
}

func (t *tflogger) addRecord(ctx context.Context, level slog.Level, msg string, args ...any) {
	if t.pc == 0 {
		t.pc, _, _, _ = runtime.Caller(2)
	}

	r := slog.NewRecord(time.Now(), level, msg, t.pc)
	r.Add(args...)
	onerror.Log(t.Logger.Handler().Handle(ctx, r))
}

func Trace(msg string, args ...any) {
	t := &tflogger{Logger: slog.Default()}
	t.pc, _, _, _ = runtime.Caller(1)
	t.Trace(msg, args...)
}

func TraceContext(ctx context.Context, msg string, args ...any) {
	t := &tflogger{Logger: slog.Default()}
	t.pc, _, _, _ = runtime.Caller(1)
	t.TraceContext(ctx, msg, args...)
}

func Fatal(msg string, args ...any) {
	t := &tflogger{Logger: slog.Default()}
	t.pc, _, _, _ = runtime.Caller(1)
	t.Fatal(msg, args...)
}

func FatalContext(ctx context.Context, msg string, args ...any) {
	t := &tflogger{Logger: slog.Default()}
	t.pc, _, _, _ = runtime.Caller(1)
	t.FatalContext(ctx, msg, args...)
}

func With(a ...any) RuntimeLogger {
	return &tflogger{Logger: slog.With(a...)}
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
