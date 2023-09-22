package adapterLog_test

import (
	"errors"
	adapterLog "shortlink/internal/adapter/log"
	"shortlink/internal/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	asrt := assert.New(t)
	defer func() {
		err := recover()
		asrt.Nil(err)
	}()

	cfg := types.CfgEnv{
		SL_APP_NAME:  "testSL",
		SL_HTTP_IP:   "localhost",
		SL_HTTP_PORT: ":8080",
		SL_DB_IP:     "localhost",
		SL_DB_PORT:   ":1313",
		SL_LOG_LEVEL: "trace",
	}

	t.Run("LogZero", func(t *testing.T) {
		log := adapterLog.NewLogZero(&cfg)
		log.LogTrace("Hello, TRACE")
		log.LogDebug("Hello, DEBUG")
		log.LogInfo("Hello, INFO")
		log.LogWarn("Hello, WARN")
		log.LogError(errors.New("test error"), "Hello, ERROR")
		log.LogFatal(errors.New("fatal error"), "Hello, FATAL")
		log.LogPanic(errors.New("panic error"), "Hello, PANIC")
	})

	t.Run("LogFprintf", func(t *testing.T) {
		log := adapterLog.NewLogFprintf(&cfg)
		log.LogTrace("Hello, TRACE")
		log.LogDebug("Hello, DEBUG")
		log.LogInfo("Hello, INFO")
		log.LogWarn("Hello, WARN")
		log.LogError(errors.New("test error"), "Hello, ERROR")
		log.LogFatal(errors.New("fatal error"), "Hello, FATAL")
		log.LogPanic(errors.New("panic error"), "Hello, PANIC")
	})

	/* 	t.Run("LogSlog", func(t *testing.T) {
		log := adapterLog.NewLogSlog(&cfg)
		log.LogTrace("Hello, TRACE")
		log.LogDebug("Hello, DEBUG")
		log.LogInfo("Hello, INFO")
		log.LogWarn("Hello, WARN")
		log.LogError(errors.New("test error"), "Hello, ERROR")
		log.LogFatal(errors.New("fatal error"), "Hello, FATAL")
		log.LogPanic(errors.New("panic error"), "Hello, PANIC")
	}) */

}
