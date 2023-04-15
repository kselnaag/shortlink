package adapters_test

import (
	"errors"
	"shortlink/internal/adapters"
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
		cfg := adapters.CfgEnv{
			APP_NAME:  "testSL",
			HTTP_IP:   "localhost",
			HTTP_PORT: ":8080",
			DB_IP:     "localhost",
			DB_PORT:   ":1313",
		}
		logger := adapters.NewLogZero(&cfg)
		logger.LogTrace("Hello, TRACE")
		logger.LogDebug("Hello, DEBUG")
		logger.LogInfo("Hello, INFO")
		logger.LogWarn("Hello, WARN")
		logger.LogError(errors.New("test error"), "Hello, ERROR")
		logger.LogFatal("Hello, FATAL")
		logger.LogPanic("Hello, PANIC")
	})

}
