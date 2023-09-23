package adapterLog

import (
	"fmt"
	"os"
	T "shortlink/internal/apptype"
	"time"
)

var _ T.ILog = (*LogFprintf)(nil)

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
	cfg    *T.CfgEnv
	loglvl LogLevel
	host   string
	svc    string
}

func NewLogFprintf(cfg *T.CfgEnv) LogFprintf {
	host := cfg.SL_HTTP_IP + cfg.SL_HTTP_PORT
	svc := cfg.SL_APP_NAME
	var lvl LogLevel
	switch cfg.SL_LOG_LEVEL {
	case T.LevelTraceValue:
		lvl = TraceLevel
	case T.LevelDebugValue:
		lvl = DebugLevel
	case T.LevelInfoValue:
		lvl = InfoLevel
	case T.LevelWarnValue:
		lvl = WarnLevel
	case T.LevelErrorValue:
		lvl = ErrorLevel
	case T.LevelFatalValue:
		lvl = FatalLevel
	case T.LevelPanicValue:
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
		logMessage(T.LevelTraceValue, l.host, l.svc, "", fmt.Sprintf(format, v...))
	}
}

func (l *LogFprintf) LogDebug(format string, v ...any) {
	if l.loglvl <= DebugLevel {
		logMessage(T.LevelDebugValue, l.host, l.svc, "", fmt.Sprintf(format, v...))
	}
}

func (l *LogFprintf) LogInfo(format string, v ...any) {
	if l.loglvl <= InfoLevel {
		logMessage(T.LevelInfoValue, l.host, l.svc, "", fmt.Sprintf(format, v...))
	}
}

func (l *LogFprintf) LogWarn(format string, v ...any) {
	if l.loglvl <= WarnLevel {
		logMessage(T.LevelWarnValue, l.host, l.svc, "", fmt.Sprintf(format, v...))
	}
}

func (l *LogFprintf) LogError(err error, format string, v ...any) {
	if l.loglvl <= ErrorLevel {
		logMessage(T.LevelErrorValue, l.host, l.svc, err.Error(), fmt.Sprintf(format, v...))
	}
}

func (l *LogFprintf) LogFatal(err error, format string, v ...any) {
	if l.loglvl <= FatalLevel {
		logMessage(T.LevelFatalValue, l.host, l.svc, err.Error(), fmt.Sprintf(format, v...))
	}
}

func (l *LogFprintf) LogPanic(err error, format string, v ...any) {
	if l.loglvl <= PanicLevel {
		logMessage(T.LevelPanicValue, l.host, l.svc, err.Error(), fmt.Sprintf(format, v...))
	}
}
