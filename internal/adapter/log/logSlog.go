package adapterLog

import (
	"os"
	T "shortlink/internal/apptype"

	"log/slog"

	"context"
	"fmt"
)

var _ T.ILog = (*LogSlog)(nil)

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
	cfg    *T.CfgEnv
	logger *slog.Logger
	ctx    context.Context
}

func NewLogSlog(cfg *T.CfgEnv) *LogSlog {
	ctx := context.Background()
	host := cfg.SL_HTTP_IP + cfg.SL_HTTP_PORT
	svc := cfg.SL_APP_NAME + cfg.SL_APP_PROTOS
	replace := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.LevelKey {
			a.Key = "L"
			level := a.Value.Any().(slog.Level)
			switch {
			case level <= LevelTrace:
				a.Value = slog.StringValue(T.LevelTraceValue)
			case level <= LevelDebug:
				a.Value = slog.StringValue(T.LevelDebugValue)
			case level <= LevelInfo:
				a.Value = slog.StringValue(T.LevelInfoValue)
			case level <= LevelWarning:
				a.Value = slog.StringValue(T.LevelWarnValue)
			case level <= LevelError:
				a.Value = slog.StringValue(T.LevelErrorValue)
			case level <= LevelFatal:
				a.Value = slog.StringValue(T.LevelFatalValue)
			case level <= LevelPanic:
				a.Value = slog.StringValue(T.LevelPanicValue)
			default:
				a.Value = slog.StringValue(T.LevelDisabledValue)
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
	case T.LevelTraceValue:
		logLevel.Set(LevelTrace)
	case T.LevelDebugValue:
		logLevel.Set(LevelDebug)
	case T.LevelInfoValue:
		logLevel.Set(LevelInfo)
	case T.LevelWarnValue:
		logLevel.Set(LevelWarning)
	case T.LevelErrorValue:
		logLevel.Set(LevelError)
	case T.LevelFatalValue:
		logLevel.Set(LevelFatal)
	case T.LevelPanicValue:
		logLevel.Set(LevelPanic)
	default:
		logLevel.Set(LevelDisabled)
	}
	return &LogSlog{
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
