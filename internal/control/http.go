package control

import (
	"encoding/json"
	"fmt"
	"regexp"
	T "shortlink/internal/apptype"
	"strings"
)

var _ T.ICtrlHTTP = (*CtrlHTTP)(nil)

type CtrlHTTP struct {
	servSL     T.ISvcShortLink
	ishashfunc func(s string) bool
}

func NewCtrlHTTP(servSL T.ISvcShortLink) *CtrlHTTP {
	isHash := regexp.MustCompile(`^[a-z0-9]{6}$`).MatchString
	return &CtrlHTTP{
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
	req := T.HTTPMessageDTO{}
	if err := json.Unmarshal(body, &req); (err != nil) || (req.IsResp) || (req.Body == "") {
		return "", fmt.Errorf("%w: %w: %s", T.ErrHTTPCtrl, T.ErrJSONUNMarshal, string(body))
	}
	lp := ctrl.servSL.GetLinkShortFromLinkLong(req.Body)
	if !lp.IsValid() {
		return "", fmt.Errorf("%w: %w: %s", T.ErrHTTPCtrl, T.ErrStructNotValid, string(body))
	}
	return lp.Short(), nil
}

func (ctrl *CtrlHTTP) Short(body []byte) (string, error) {
	req := T.HTTPMessageDTO{}
	if err := json.Unmarshal(body, &req); (err != nil) || (req.IsResp) || (!ctrl.ishashfunc(req.Body)) {
		return "", fmt.Errorf("%w: %w: %s", T.ErrHTTPCtrl, T.ErrJSONUNMarshal, string(body))
	}
	lp := ctrl.servSL.GetLinkLongFromLinkShort(req.Body)
	if !lp.IsValid() {
		return "", fmt.Errorf("%w: %w: %s", T.ErrHTTPCtrl, T.ErrStructNotValid, string(body))
	}
	return lp.Long(), nil
}

func (ctrl *CtrlHTTP) Save(body []byte) (string, error) {
	req := T.HTTPMessageDTO{}
	if err := json.Unmarshal(body, &req); (err != nil) || (req.IsResp) || (req.Body == "") {
		return "", fmt.Errorf("%w: %w: %s", T.ErrHTTPCtrl, T.ErrJSONUNMarshal, string(body))
	}
	lp := ctrl.servSL.SetLinkPairFromLinkLong(req.Body)
	if !lp.IsValid() {
		return "", fmt.Errorf("%w: %w: %s", T.ErrHTTPCtrl, T.ErrStructNotValid, string(body))
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
	return "", fmt.Errorf("%w: %w: %s", T.ErrHTTPCtrl, T.ErrHashCorrect, hash)
}
