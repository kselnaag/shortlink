package adapterLog

import (
	"os"
	"shortlink/internal/types"
	"time"

	"github.com/rs/zerolog"
)

var _ types.ILog = (*LogZero)(nil)

type LogZero struct {
	cfg    *types.CfgEnv
	logger zerolog.Logger
}

func NewLogZero(cfg *types.CfgEnv) LogZero {
	host := cfg.SL_HTTP_IP + cfg.SL_HTTP_PORT
	service := cfg.SL_APP_NAME
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.TimestampFieldName = "T"
	zerolog.LevelFieldName = "L"
	zerolog.MessageFieldName = "M"
	zerolog.ErrorFieldName = "E" //nolint:reassign // standart reassigning ops for package

	var lvl zerolog.Level
	switch cfg.SL_LOG_LEVEL {
	case zerolog.LevelTraceValue:
		lvl = zerolog.TraceLevel
	case zerolog.LevelDebugValue:
		lvl = zerolog.DebugLevel
	case zerolog.LevelInfoValue:
		lvl = zerolog.InfoLevel
	case zerolog.LevelWarnValue:
		lvl = zerolog.WarnLevel
	case zerolog.LevelErrorValue:
		lvl = zerolog.ErrorLevel
	case zerolog.LevelFatalValue:
		lvl = zerolog.FatalLevel
	case zerolog.LevelPanicValue:
		lvl = zerolog.PanicLevel
	default:
		lvl = zerolog.Disabled
	}

	newlogger := zerolog.New(os.Stderr).Level(lvl).With().
		Timestamp().Str("H", host).Str("S", service).
		Logger()
	return LogZero{
		cfg:    cfg,
		logger: newlogger,
	}
}

func (l *LogZero) LogTrace(format string, v ...any) {
	l.logger.Trace().Msgf(format, v...)
}

func (l *LogZero) LogDebug(format string, v ...any) {
	l.logger.Debug().Msgf(format, v...)
}

func (l *LogZero) LogInfo(format string, v ...any) {
	l.logger.Info().Msgf(format, v...)
}

func (l *LogZero) LogWarn(format string, v ...any) {
	l.logger.Warn().Msgf(format, v...)
}

func (l *LogZero) LogError(err error, format string, v ...any) {
	l.logger.Error().Err(err).Msgf(format, v...)
}

func (l *LogZero) LogFatal(err error, format string, v ...any) {
	l.logger.WithLevel(zerolog.FatalLevel).Err(err).Msgf(format, v...)
}

func (l *LogZero) LogPanic(err error, format string, v ...any) {
	l.logger.WithLevel(zerolog.PanicLevel).Err(err).Msgf(format, v...)
}
