package services_test

import (
	"hash/crc32"
	"shortlink/internal/models"
	"shortlink/internal/services"
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
		assert.Equal(uint32(0x352441c2), crc32.ChecksumIEEE([]byte("abc")))
		assert.Equal(uint32(0xab40d461), crc32.ChecksumIEEE([]byte("abd")))
		assert.Equal(uint32(0x69e69b98), crc32.ChecksumIEEE([]byte("   ")))
		assert.Equal(uint32(0xdc2573e1), crc32.ChecksumIEEE([]byte("Hello,world!")))

	})
	t.Run("CalcLinkShort", func(t *testing.T) {
		assert.Equal("eqtepu", services.CalcLinkShort("abc"))
		assert.Equal("1bilmqp", services.CalcLinkShort("abd"))
		assert.Equal("tdtajc", services.CalcLinkShort("   "))
		assert.Equal("1p2za3l", services.CalcLinkShort("Hello,world!"))
	})
	t.Run("IsPairValid", func(t *testing.T) {
		assert.True(services.IsPairValid(models.LinkPair{Short: "eqtepu", Long: "abc"}))
		assert.True(services.IsPairValid(models.LinkPair{Short: "1bilmqp", Long: "abd"}))
		assert.True(services.IsPairValid(models.LinkPair{Short: "1p2za3l", Long: "Hello,world!"}))

		assert.False(services.IsPairValid(models.LinkPair{Short: "tdtajc", Long: "   "}))
		assert.False(services.IsPairValid(models.LinkPair{Short: "  ", Long: "Hello,world!"}))
		assert.False(services.IsPairValid(models.LinkPair{Short: "", Long: " "}))
	})
}
