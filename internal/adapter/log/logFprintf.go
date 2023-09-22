package adapterLog

import (
	"fmt"
	"os"
	"shortlink/internal/types"
	"time"
)

var _ types.ILog = (*LogFprintf)(nil)

type LogLevel int8

const (
	TraceLevel LogLevel = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
	Disabled
)

type LogFprintf struct {
	cfg    *types.CfgEnv
	loglvl LogLevel
	host   string
	svc    string
}

func NewLogFprintf(cfg *types.CfgEnv) LogFprintf {
	host := cfg.SL_HTTP_IP + cfg.SL_HTTP_PORT
	svc := cfg.SL_APP_NAME
	var lvl LogLevel
	switch cfg.SL_LOG_LEVEL {
	case types.LevelTraceValue:
		lvl = TraceLevel
	case types.LevelDebugValue:
		lvl = DebugLevel
	case types.LevelInfoValue:
		lvl = InfoLevel
	case types.LevelWarnValue:
		lvl = WarnLevel
	case types.LevelErrorValue:
		lvl = ErrorLevel
	case types.LevelFatalValue:
		lvl = FatalLevel
	case types.LevelPanicValue:
		lvl = PanicLevel
	default:
		lvl = Disabled
	}
	return LogFprintf{
		cfg:    cfg,
		loglvl: lvl,
		host:   host,
		svc:    svc,
	}
}

func logMessage(lvl, host, svc, err, mess string) {
	timenow := time.Now().Format(time.RFC3339Nano)
	fmt.Fprintf(os.Stderr, "{\"L\":\"%s\",\"T\":\"%s\",\"H\":\"%s\",\"S\":\"%s\",\"M\":\"%s\",\"E\":\"%s\"}\n", lvl, timenow, host, svc, mess, err)
}

func (l *LogFprintf) LogTrace(format string, v ...any) {
	if l.loglvl <= TraceLevel {
		logMessage(types.LevelTraceValue, l.host, l.svc, "", fmt.Sprintf(format, v...))
	}
}

func (l *LogFprintf) LogDebug(format string, v ...any) {
	if l.loglvl <= DebugLevel {
		logMessage(types.LevelDebugValue, l.host, l.svc, "", fmt.Sprintf(format, v...))
	}
}

func (l *LogFprintf) LogInfo(format string, v ...any) {
	if l.loglvl <= InfoLevel {
		logMessage(types.LevelInfoValue, l.host, l.svc, "", fmt.Sprintf(format, v...))
	}
}

func (l *LogFprintf) LogWarn(format string, v ...any) {
	if l.loglvl <= WarnLevel {
		logMessage(types.LevelWarnValue, l.host, l.svc, "", fmt.Sprintf(format, v...))
	}
}

func (l *LogFprintf) LogError(err error, format string, v ...any) {
	if l.loglvl <= ErrorLevel {
		logMessage(types.LevelErrorValue, l.host, l.svc, err.Error(), fmt.Sprintf(format, v...))
	}
}

func (l *LogFprintf) LogFatal(err error, format string, v ...any) {
	if l.loglvl <= FatalLevel {
		logMessage(types.LevelFatalValue, l.host, l.svc, err.Error(), fmt.Sprintf(format, v...))
	}
}

func (l *LogFprintf) LogPanic(err error, format string, v ...any) {
	if l.loglvl <= PanicLevel {
		logMessage(types.LevelPanicValue, l.host, l.svc, err.Error(), fmt.Sprintf(format, v...))
	}
}
