package adapterLog_test

import (
	"errors"
	adapterCfg "shortlink/internal/adapters/cfg"
	adapterLog "shortlink/internal/adapters/log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogZero(t *testing.T) {
	asrt := assert.New(t)
	defer func() {
		err := recover()
		asrt.Nil(err)
	}()

	t.Run("LogZero", func(t *testing.T) {
		cfg := adapterCfg.CfgEnv{
			SL_APP_NAME:  "testSL",
			SL_HTTP_IP:   "localhost",
			SL_HTTP_PORT: ":8080",
			SL_DB_IP:     "localhost",
			SL_DB_PORT:   ":1313",
		}
		log := adapterLog.NewLogZero(&cfg)
		log.LogTrace("Hello, TRACE")
		log.LogDebug("Hello, DEBUG")
		log.LogInfo("Hello, INFO")
		log.LogWarn("Hello, WARN")
		log.LogError(errors.New("test error"), "Hello, ERROR")
		log.LogFatal("Hello, FATAL")
		log.LogPanic("Hello, PANIC")
	})

}
