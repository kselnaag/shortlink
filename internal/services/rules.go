package services

import (
	"hash/crc32"
	"strconv"
	"strings"

	"shortlink/internal/models"
)

func IsPairValid(lp models.LinkPair) bool {
	isLLvalid := IsLinkValid(lp.Long)
	isLSvalid := IsLinkValid(lp.Short)
	isHASHvalid := (CalcLinkShort(lp.Long) == lp.Short)
	return isLLvalid && isLSvalid && isHASHvalid
}

func IsLinkValid(link string) bool {
	link = strings.Trim(link, " ")
	return len(link) != 0
}

func UniteLinks(linkl string) models.LinkPair {
	return models.LinkPair{
		Short: CalcLinkShort(linkl),
		Long:  strings.Clone(linkl),
	}
}

func CalcLinkShort(linkl string) string {
	hash := crc32.ChecksumIEEE([]byte(linkl))
	return strconv.FormatUint(uint64(hash), 36)
}

/*
+ compute short link from long link
+ unite short link and long link
+ check if pair is valid
*/
