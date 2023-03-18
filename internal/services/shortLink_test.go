package services_test

import (
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
		log := adapters.NewLogZero("localhost:8080", "testSL")
		db := adapters.NewDBMock()
		hcli := adapters.NewHttpMockClient()
		nssl := services.NewServShortLink(&db, &hcli, &log)

		// models.LinkPair{Short: "5clp60", Long: "http://lib.ru"}, models.LinkPair{Short: "8as3rb", Long: "http://lib.ru/abs"}, ("dhiu79", "http://google.ru")
		assert.Equal([]models.LinkPair{models.NewLinkPair("http://lib.ru"), models.NewLinkPair("http://google.ru")},
			nssl.GetAllLinkPairs())

		assert.True(nssl.IsLinkLongHttpValid("http://lib.ru/PROZA/"))
		assert.False(nssl.IsLinkLongHttpValid("http://lib.ru/abs"))

		assert.True(nssl.SetLinkPairFromLinkLong("http://lib.ru/PROZA/").IsValid())
		assert.False(nssl.SetLinkPairFromLinkLong("http://lib.ru/abs").IsValid())

		lp := nssl.GetLinkLongFromLinkShort("8b4s29")
		assert.True(lp.IsValid())
		assert.Equal("http://lib.ru/PROZA/", lp.Long())
		assert.False(nssl.GetLinkLongFromLinkShort("8as3rb").IsValid())

		assert.Equal([]models.LinkPair{models.NewLinkPair("http://lib.ru"), models.NewLinkPair("http://lib.ru/PROZA/"), models.NewLinkPair("http://google.ru")},
			nssl.GetAllLinkPairs())
	})

}
