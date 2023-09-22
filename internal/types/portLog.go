package types

type ILog interface {
	LogTrace(format string, v ...any)
	LogDebug(format string, v ...any)
	LogInfo(format string, v ...any)
	LogWarn(format string, v ...any)
	LogError(err error, format string, v ...any)
	LogFatal(format string, v ...any)
	LogPanic(format string, v ...any)
}
