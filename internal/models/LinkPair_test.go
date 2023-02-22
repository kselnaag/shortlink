package models_test

import (
	mod "shortlink/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModels(t *testing.T) {
	assert := assert.New(t)
	defer func() {
		err := recover()
		assert.Nil(err)
	}()

	t.Run("LinkPair", func(t *testing.T) {
		lp := mod.LinkPair{}
		assert.False(lp.IsValid())
		assert.Equal("", lp.Long)
		assert.Equal("", lp.Short)

		lp1 := mod.NewLinkPair("")
		assert.False(lp1.IsValid())
		assert.Equal("", lp1.Long)
		assert.Equal("", lp1.Short)

		lp2 := mod.NewLinkPair("  ")
		assert.False(lp2.IsValid())
		assert.Equal("", lp2.Long)
		assert.Equal("", lp2.Short)

		lp3 := mod.NewLinkPair("abc")
		assert.True(lp3.IsValid())
		assert.Equal("abc", lp3.Long)
		assert.Equal("eqtepu", lp3.Short)

		lp4 := mod.NewLinkPair(" a b c ")
		assert.True(lp4.IsValid())
		assert.Equal("a b c", lp4.Long)
		assert.Equal("jjjhe8", lp4.Short)

		lp5 := mod.NewLinkPair("abd")
		assert.True(lp5.IsValid())
		assert.Equal("abd", lp5.Long)
		assert.Equal("bilmqp", lp5.Short)

		lp6 := mod.NewLinkPair(" \nHello, world !\t")
		assert.True(lp6.IsValid())
		assert.Equal("Hello, world !", lp6.Long)
		assert.Equal("lf4w7t", lp6.Short)
	})

}
