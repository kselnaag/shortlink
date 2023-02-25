package models

import (
	"hash/crc32"
	"strconv"
	"strings"
)

type LinkPair struct {
	Short string
	Long  string
}

func (lp LinkPair) IsValid() bool {
	isLLvalid := isLinkValid(lp.Long)
	isLSvalid := isLinkValid(lp.Short)
	isHASHvalid := (calcLinkShort(lp.Long) == lp.Short)
	return (isLLvalid && isLSvalid && isHASHvalid)
}

func NewLinkPair(linklong string) LinkPair {
	res := LinkPair{}
	link := formatLink(linklong)
	if !isLinkValid(link) {
		return res
	}
	res.Long = link
	res.Short = calcLinkShort(link)
	return res
}

func formatLink(link string) string {
	return strings.TrimSpace(link)
}

func isLinkValid(link string) bool {
	return len(link) != 0
}

func calcLinkShort(linkl string) string {
	hash := crc32.ChecksumIEEE([]byte(linkl))
	str := strconv.FormatUint(uint64(hash), 36)
	if len(str) > 6 {
		idx := len(str) - 6
		return str[idx:]
	}
	return str
}

/*
+ compute short link from long link
+ unite short link and long link
+ check if pair is valid
*/
