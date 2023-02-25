package services_test

import (
	"context"
	"testing"

	"shortlink/internal/adapters"
	"shortlink/internal/models"
	"shortlink/internal/services"

	"github.com/stretchr/testify/assert"
)

func TestServices(t *testing.T) {
	assert := assert.New(t)
	defer func() {
		err := recover()
		assert.Nil(err)
	}()

	t.Run("ServShortLink", func(t *testing.T) {

		ctx, _ := context.WithCancel(context.Background())
		db := adapters.NewMockDB(&ctx)
		hcli := adapters.NewHttpMockClient(&ctx)
		nssl := services.NewServShortLink(&ctx, &db, &hcli)

		// models.LinkPair{Short: "5clp60", Long: "http://lib.ru"}, models.LinkPair{Short: "8as3rb", Long: "http://lib.ru/abs"}
		assert.Equal([]models.LinkPair{}, nssl.GetAllLinkPairs())

		assert.True(nssl.IsLinkLongHttpValid("http://lib.ru"))
		assert.False(nssl.IsLinkLongHttpValid("http://lib.ru/abs"))

		assert.True(nssl.SetLinkPairFromLinkLong("http://lib.ru"))
		assert.False(nssl.SetLinkPairFromLinkLong("http://lib.ru/abs"))

		lp := nssl.GetLinkLongFromLinkShort("5clp60")
		assert.True(lp.IsValid())
		assert.Equal("http://lib.ru", lp.Long)
		assert.False(nssl.GetLinkLongFromLinkShort("8as3rb").IsValid())

		assert.Equal([]models.LinkPair{{Short: "5clp60", Long: "http://lib.ru"}}, nssl.GetAllLinkPairs())

		/* func NewServShortLink(ctx *context.Context, db *ports.Idb, hcli *ports.IHttpClient) ServShortLink
		func (ssl *ServShortLink) GetAllLinkPairs() []models.LinkPair
		func (ssl *ServShortLink) GetLinkLongFromLinkShort(linkshort string) models.LinkPair
		func (ssl *ServShortLink) SetLinkPairFromLinkLong(linklong string) bool
		func (ssl *ServShortLink) IsLinkLongHttpValid(linklong string) bool */

	})

}
