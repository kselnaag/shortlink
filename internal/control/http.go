package control

import (
	"encoding/json"
	"errors"
	"regexp"
	"shortlink/internal/i7e"
	"strings"
)

var _ i7e.ICtrlHTTP = (*CtrlHTTP)(nil)

type Message struct {
	IsResp bool
	Mode   string
	Body   string
}

type CtrlHTTP struct {
	servSL     i7e.ISvcShortLink
	ishashfunc func(s string) bool
}

func NewCtrlHTTP(servSL i7e.ISvcShortLink) CtrlHTTP {
	isHash := regexp.MustCompile(`^[a-z0-9]{6}$`).MatchString
	return CtrlHTTP{
		servSL:     servSL,
		ishashfunc: isHash,
	}
}

func (ctrl *CtrlHTTP) AllPairs() (string, error) {
	strarr := []string{}
	pairs := ctrl.servSL.GetAllLinkPairs()
	for _, el := range pairs {
		strarr = append(strarr, el.Short()+": "+el.Long())
	}
	return strings.Join(strarr, "; "), nil
}

func (ctrl *CtrlHTTP) Long(body []byte) (string, error) {
	req := Message{}
	if err := json.Unmarshal(body, &req); (err != nil) || (req.IsResp) || (req.Body == "") {
		return "", errors.New(string(body))
	}
	lp := ctrl.servSL.GetLinkShortFromLinkLong(req.Body)
	if !lp.IsValid() {
		return "", errors.New(string(body))
	}
	return lp.Short(), nil
}

func (ctrl *CtrlHTTP) Short(body []byte) (string, error) {
	req := Message{}
	if err := json.Unmarshal(body, &req); (err != nil) || (req.IsResp) || (!ctrl.ishashfunc(req.Body)) {
		return "", errors.New(string(body))
	}
	lp := ctrl.servSL.GetLinkLongFromLinkShort(req.Body)
	if !lp.IsValid() {
		return "", errors.New(string(body))
	}
	return lp.Long(), nil
}

func (ctrl *CtrlHTTP) Save(body []byte) (string, error) {
	req := Message{}
	if err := json.Unmarshal(body, &req); (err != nil) || (req.IsResp) || (req.Body == "") {
		return "", errors.New(string(body))
	}
	lp := ctrl.servSL.SetLinkPairFromLinkLong(string(body))
	if !lp.IsValid() {
		return "", errors.New(string(body))
	}
	return lp.Short(), nil
}

func (ctrl *CtrlHTTP) Hash(hash string) (string, error) {
	if ctrl.ishashfunc(hash) {
		lp := ctrl.servSL.GetLinkLongFromLinkShort(hash)
		if lp.IsValid() {
			return lp.Long(), nil
		}
	}
	return "", errors.New(hash)
}
