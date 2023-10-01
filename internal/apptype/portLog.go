package types

const (
	LevelTraceValue    = "trace"
	LevelDebugValue    = "debug"
	LevelInfoValue     = "info"
	LevelWarnValue     = "warn"
	LevelErrorValue    = "error"
	LevelFatalValue    = "fatal"
	LevelPanicValue    = "panic"
	LevelDisabledValue = "disabled"
)

type ILog interface {
	LogTrace(format string, v ...any)
	LogDebug(format string, v ...any)
	LogInfo(format string, v ...any)
	LogWarn(format string, v ...any)
	LogError(err error, format string, v ...any)
	LogFatal(err error, format string, v ...any)
	LogPanic(err error, format string, v ...any)
}
