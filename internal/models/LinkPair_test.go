package models_test

import (
	"shortlink/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModels(t *testing.T) {
	asrt := assert.New(t)
	defer func() {
		err := recover()
		asrt.Nil(err)
	}()

	t.Run("LinkPair", func(t *testing.T) {
		lp := models.LinkPair{}
		asrt.False(lp.IsValid())
		asrt.Equal("", lp.Long())
		asrt.Equal("", lp.Short())

		lp1 := models.NewLinkPair("")
		asrt.False(lp1.IsValid())
		asrt.Equal("", lp1.Long())
		asrt.Equal("", lp1.Short())

		lp2 := models.NewLinkPair("  ")
		asrt.False(lp2.IsValid())
		asrt.Equal("", lp2.Long())
		asrt.Equal("", lp2.Short())

		lp3 := models.NewLinkPair("abc")
		asrt.True(lp3.IsValid())
		asrt.Equal("abc", lp3.Long())
		asrt.Equal("eqtepu", lp3.Short())

		lp4 := models.NewLinkPair(" a b c ")
		asrt.True(lp4.IsValid())
		asrt.Equal("a b c", lp4.Long())
		asrt.Equal("jjjhe8", lp4.Short())

		lp5 := models.NewLinkPair("abd")
		asrt.True(lp5.IsValid())
		asrt.Equal("abd", lp5.Long())
		asrt.Equal("bilmqp", lp5.Short())

		lp6 := models.NewLinkPair(" \nHello, world !\t")
		asrt.True(lp6.IsValid())
		asrt.Equal("Hello, world !", lp6.Long())
		asrt.Equal("lf4w7t", lp6.Short())
	})
}
