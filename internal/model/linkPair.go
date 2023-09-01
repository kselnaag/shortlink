package model

import (
	"hash/crc32"
	"strconv"
	"strings"
)

type LinkPair struct {
	short string
	long  string
}

func (lp LinkPair) IsValid() bool {
	isLLvalid := isLinkValid(lp.long)
	isLSvalid := isLinkValid(lp.short)
	isHASHvalid := (calcLinkShort(lp.long) == lp.short)
	return (isLLvalid && isLSvalid && isHASHvalid)
}

func (lp LinkPair) Short() string {
	return lp.short
}

func (lp LinkPair) Long() string {
	return lp.long
}

func NewLinkPair(linklong string) LinkPair {
	link := formatLink(linklong)
	if !isLinkValid(link) {
		return LinkPair{}
	}
	return LinkPair{short: calcLinkShort(link), long: link}
}

func formatLink(link string) string {
	return strings.TrimSpace(link)
}

func isLinkValid(link string) bool {
	return link != ""
}

func calcLinkShort(linkl string) string {
	hashlen := 6
	radixlen := 36
	hash := crc32.ChecksumIEEE([]byte(linkl))
	str := strconv.FormatUint(uint64(hash), radixlen)
	if len(str) > hashlen {
		idx := len(str) - hashlen
		return str[idx:]
	}
	return str
}

/*
+ compute short link from long link
+ unite short link and long link
+ check if pair is valid
*/
