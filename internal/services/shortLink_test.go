package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRules(t *testing.T) {
	assert := assert.New(t)
	defer func() {
		err := recover()
		assert.Nil(err)
	}()

	t.Run("HashCRC32", func(t *testing.T) {

	})

}
