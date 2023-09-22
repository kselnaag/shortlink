package adapterLog

import (
	"os"
	"shortlink/internal/types"

	"log/slog"

	"context"
	"fmt"
)

var _ types.ILog = (*LogSlog)(nil)

const (
	LevelTrace    = slog.Level(-8)
	LevelDebug    = slog.LevelDebug
	LevelInfo     = slog.LevelInfo
	LevelWarning  = slog.LevelWarn
	LevelError    = slog.LevelError
	LevelFatal    = slog.Level(12)
	LevelPanic    = slog.Level(16)
	LevelDisabled = slog.Level(20)
)

type LogSlog struct {
	cfg    *types.CfgEnv
	logger *slog.Logger
	ctx    context.Context
}

func NewLogSlog(cfg *types.CfgEnv) LogSlog {
	ctx := context.Background()
	host := cfg.SL_HTTP_IP + cfg.SL_HTTP_PORT
	svc := cfg.SL_APP_NAME
	replace := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.LevelKey {
			a.Key = "L"
			level := a.Value.Any().(slog.Level)
			switch {
			case level <= LevelTrace:
				a.Value = slog.StringValue(types.LevelTraceValue)
			case level <= LevelDebug:
				a.Value = slog.StringValue(types.LevelDebugValue)
			case level <= LevelInfo:
				a.Value = slog.StringValue(types.LevelInfoValue)
			case level <= LevelWarning:
				a.Value = slog.StringValue(types.LevelWarnValue)
			case level <= LevelError:
				a.Value = slog.StringValue(types.LevelErrorValue)
			case level <= LevelFatal:
				a.Value = slog.StringValue(types.LevelFatalValue)
			case level <= LevelPanic:
				a.Value = slog.StringValue(types.LevelPanicValue)
			default:
				a.Value = slog.StringValue(types.LevelDisabledValue)
			}
		}
		if a.Key == slog.TimeKey {
			a.Key = "T"
		}
		if a.Key == slog.MessageKey {
			a.Key = "M"
		}
		if a.Key == slog.SourceKey {
			a.Key = "P"
		}
		return a
	}
	var logLevel = new(slog.LevelVar)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: false, Level: logLevel, ReplaceAttr: replace}))
	logger = logger.With("H", host).With("S", svc)
	slog.SetDefault(logger)
	switch cfg.SL_LOG_LEVEL {
	case types.LevelTraceValue:
		logLevel.Set(LevelTrace)
	case types.LevelDebugValue:
		logLevel.Set(LevelDebug)
	case types.LevelInfoValue:
		logLevel.Set(LevelInfo)
	case types.LevelWarnValue:
		logLevel.Set(LevelWarning)
	case types.LevelErrorValue:
		logLevel.Set(LevelError)
	case types.LevelFatalValue:
		logLevel.Set(LevelFatal)
	case types.LevelPanicValue:
		logLevel.Set(LevelPanic)
	default:
		logLevel.Set(LevelDisabled)
	}
	return LogSlog{
		ctx:    ctx,
		cfg:    cfg,
		logger: logger,
	}
}

func (l *LogSlog) LogTrace(format string, v ...any) {
	l.logger.Log(l.ctx, LevelTrace, fmt.Sprintf(format, v...), "E", "")
}

func (l *LogSlog) LogDebug(format string, v ...any) {
	l.logger.Debug(fmt.Sprintf(format, v...), "E", "")
}

func (l *LogSlog) LogInfo(format string, v ...any) {
	l.logger.Info(fmt.Sprintf(format, v...), "E", "")
}

func (l *LogSlog) LogWarn(format string, v ...any) {
	l.logger.Warn(fmt.Sprintf(format, v...), "E", "")
}

func (l *LogSlog) LogError(err error, format string, v ...any) {
	l.logger.Error(fmt.Sprintf(format, v...), "E", err.Error())
}

func (l *LogSlog) LogFatal(err error, format string, v ...any) {
	l.logger.Log(l.ctx, LevelFatal, fmt.Sprintf(format, v...), "E", err.Error())
}

func (l *LogSlog) LogPanic(err error, format string, v ...any) {
	l.logger.Log(l.ctx, LevelPanic, fmt.Sprintf(format, v...), "E", err.Error())
}
