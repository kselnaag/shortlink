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
	case LevelTraceValue:
		lvl = TraceLevel
	case LevelDebugValue:
		lvl = DebugLevel
	case LevelInfoValue:
		lvl = InfoLevel
	case LevelWarnValue:
		lvl = WarnLevel
	case LevelErrorValue:
		lvl = ErrorLevel
	case LevelFatalValue:
		lvl = FatalLevel
	case LevelPanicValue:
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
		logMessage(LevelTraceValue, l.host, l.svc, "", fmt.Sprintf(format, v...))
	}
}

func (l *LogFprintf) LogDebug(format string, v ...any) {
	if l.loglvl <= DebugLevel {
		logMessage(LevelDebugValue, l.host, l.svc, "", fmt.Sprintf(format, v...))
	}
}

func (l *LogFprintf) LogInfo(format string, v ...any) {
	if l.loglvl <= InfoLevel {
		logMessage(LevelInfoValue, l.host, l.svc, "", fmt.Sprintf(format, v...))
	}
}

func (l *LogFprintf) LogWarn(format string, v ...any) {
	if l.loglvl <= WarnLevel {
		logMessage(LevelWarnValue, l.host, l.svc, "", fmt.Sprintf(format, v...))
	}
}

func (l *LogFprintf) LogError(err error, format string, v ...any) {
	if l.loglvl <= ErrorLevel {
		logMessage(LevelErrorValue, l.host, l.svc, err.Error(), fmt.Sprintf(format, v...))
	}
}

func (l *LogFprintf) LogFatal(err error, format string, v ...any) {
	if l.loglvl <= FatalLevel {
		logMessage(LevelFatalValue, l.host, l.svc, err.Error(), fmt.Sprintf(format, v...))
	}
}

func (l *LogFprintf) LogPanic(err error, format string, v ...any) {
	if l.loglvl <= PanicLevel {
		logMessage(LevelPanicValue, l.host, l.svc, err.Error(), fmt.Sprintf(format, v...))
	}
}
