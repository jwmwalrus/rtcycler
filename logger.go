package rtc

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
)

func Trace(msg string, args ...any) {
	slog.Default().Log(context.Background(), LevelTrace, msg, args...)
}

func TraceContext(ctx context.Context, msg string, args ...any) {
	slog.Default().Log(ctx, LevelFatal, msg, args...)
}

func Fatal(msg string, args ...any) {
	slog.Default().Log(context.Background(), LevelFatal, msg, args...)
}

func FatalContext(ctx context.Context, msg string, args ...any) {
	slog.Default().Log(ctx, LevelFatal, msg, args...)
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
