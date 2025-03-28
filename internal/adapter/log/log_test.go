package adapterLog_test

import (
	"fmt"
	adapterLog "shortlink/internal/adapter/log"
	T "shortlink/internal/apptype"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	asrt := assert.New(t)
	defer func() {
		err := recover()
		asrt.Nil(err)
	}()

	cfg := T.CfgEnv{
		SL_APP_NAME:  "testSL",
		SL_HTTP_IP:   "localhost",
		SL_HTTP_PORT: ":8080",
		SL_DB_IP:     "localhost",
		SL_DB_PORT:   ":1313",
		SL_LOG_LEVEL: "trace",
	}

	t.Run("LogZero", func(t *testing.T) {
		err := fmt.Errorf("%w", T.ErrTestLog)
		log := adapterLog.NewLogZero(&cfg)
		log.LogTrace("Hello, TRACE")
		log.LogDebug("Hello, DEBUG")
		log.LogInfo("Hello, INFO")
		log.LogWarn("Hello, WARN")
		log.LogError(err, "Hello,error")
		log.LogFatal(err, "Hello,fatal")
		log.LogPanic(err, "Hello,PANIC")
	})

	t.Run("LogFprintf", func(t *testing.T) {
		err := fmt.Errorf("%w", T.ErrTestLog)
		log := adapterLog.NewLogFprintf(&cfg)
		log.LogTrace("Hello, TRACE")
		log.LogDebug("Hello, DEBUG")
		log.LogInfo("Hello, INFO")
		log.LogWarn("Hello, WARN")
		log.LogError(err, "Hello,error")
		log.LogFatal(err, "Hello,fatal")
		log.LogPanic(err, "Hello,panic")
	})

	t.Run("LogSlog", func(t *testing.T) {
		err := fmt.Errorf("%w", T.ErrTestLog)
		log := adapterLog.NewLogSlog(&cfg)
		log.LogTrace("Hello, TRACE")
		log.LogDebug("Hello, DEBUG")
		log.LogInfo("Hello, INFO")
		log.LogWarn("Hello, WARN")
		log.LogError(err, "Hello,error")
		log.LogFatal(err, "Hello,fatal")
		log.LogPanic(err, "Hello,panic")
	})
}
