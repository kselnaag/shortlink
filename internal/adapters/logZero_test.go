package adapters_test

import (
	"errors"
	"shortlink/internal/adapters"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogZero(t *testing.T) {
	assert := assert.New(t)
	defer func() {
		err := recover()
		assert.Nil(err)
	}()

	t.Run("LogZero", func(t *testing.T) {
		logger := adapters.NewLogZero("localhost", "testSL")
		logger.LogTrace("Hello, TRACE")
		logger.LogDebug("Hello, DEBUG")
		logger.LogInfo("Hello, INFO")
		logger.LogWarn("Hello, WARN")
		logger.LogError(errors.New("test error"), "Hello, ERROR")
		logger.LogFatal("Hello, FATAL")
		logger.LogPanic("Hello, PANIC")
	})

}
