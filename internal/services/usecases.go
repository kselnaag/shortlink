package services

import (
	mod "shortlink/internal/models"
)

func LoadAllLinkPairs(db *mod.Idb) []mod.LinkPair {
	result := []mod.LinkPair{}
	allpairs := (*db).LoadAllLinkPairs()
	for _, el := range allpairs {
		if IsPairValid(el) {
			result = append(result, el)
		}
	}
	return result
}

func LoadLinkLongFromLinkShort(links string, db *mod.Idb) []mod.LinkPair {
	lparr := (*db).LoadLinkPair(links)
	if (len(lparr) == 1) && (IsPairValid(lparr[0])) {
		return lparr
	}
	return []mod.LinkPair{}
}

func SaveLinkPairFromLinkLong(linkl string, db *mod.Idb) []mod.LinkPair {
	if IsLinkLongHttpValid(linkl, "http!") { //!
		lp := UniteLinks(linkl)
		if (*db).SaveLinkPair(lp) {
			return append([]mod.LinkPair{}, lp) //
		}
	}
	return []mod.LinkPair{}
}

func IsLinkLongHttpValid(linkl string, http string) bool {
	// HTTP !
	return true
}

/*
redirect from short link to long link
send html UI
check if long link is valid

+get the short link if you have a long link
+get the long link if you have a short link
+get ALL link pairs presented in db
*/
