package adapters

import (
	"os"
	"shortlink/internal/ports"

	"github.com/rs/zerolog"
)

var _ ports.ILog = (*LogZero)(nil)

type LogZero struct {
	cfg    *CfgEnv
	logger zerolog.Logger
}

func NewLogZero(cfg *CfgEnv) LogZero {
	host := cfg.HTTP_IP + cfg.HTTP_PORT
	service := cfg.APP_NAME
	zerolog.TimestampFieldName = "T"
	zerolog.LevelFieldName = "L"
	zerolog.MessageFieldName = "M"
	zerolog.ErrorFieldName = "E" //nolint:reassign // standart reassigning ops for package
	newlogger := zerolog.New(os.Stderr).Level(zerolog.TraceLevel).With().
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

func (l *LogZero) LogFatal(format string, v ...any) {
	l.logger.WithLevel(zerolog.FatalLevel).Msgf(format, v...)
}

func (l *LogZero) LogPanic(format string, v ...any) {
	l.logger.WithLevel(zerolog.PanicLevel).Msgf(format, v...)
}
