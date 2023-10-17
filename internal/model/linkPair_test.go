package model_test

import (
	"shortlink/internal/model"
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
		lp := model.LinkPair{}
		asrt.False(lp.IsValid())
		asrt.Equal("", lp.Long())
		asrt.Equal("", lp.Short())

		lp1 := model.NewLinkPair("")
		asrt.False(lp1.IsValid())
		asrt.Equal("", lp1.Long())
		asrt.Equal("", lp1.Short())

		lp2 := model.NewLinkPair("  ")
		asrt.False(lp2.IsValid())
		asrt.Equal("", lp2.Long())
		asrt.Equal("", lp2.Short())

		lp3 := model.NewLinkPair("abc")
		asrt.True(lp3.IsValid())
		asrt.Equal("abc", lp3.Long())
		asrt.Equal("eqtepu", lp3.Short())

		lp4 := model.NewLinkPair(" a b c ")
		asrt.True(lp4.IsValid())
		asrt.Equal("a b c", lp4.Long())
		asrt.Equal("jjjhe8", lp4.Short())

		lp5 := model.NewLinkPair("abd")
		asrt.True(lp5.IsValid())
		asrt.Equal("abd", lp5.Long())
		asrt.Equal("bilmqp", lp5.Short())

		lp6 := model.NewLinkPair(" \nHello, world !\t")
		asrt.True(lp6.IsValid())
		asrt.Equal("Hello, world !", lp6.Long())
		asrt.Equal("lf4w7t", lp6.Short())

		lp7 := model.NewLinkPair("http://lib.ru/PROZA/")
		asrt.True(lp7.IsValid())
		asrt.Equal("http://lib.ru/PROZA/", lp7.Long())
		asrt.Equal("8b4s29", lp7.Short())
	})
}
