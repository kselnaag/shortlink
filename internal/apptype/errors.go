package types

import "errors"

var (
	ErrHTTPClientFast = errors.New("HTTPClientFast error")
	ErrHTTPClientNet  = errors.New("HTTPClientNet error")
	ErrHTTPClientMock = errors.New("HTTPClientMock error")
	ErrLinkNotAval    = errors.New("link is not available")
	ErrGetMethod      = errors.New("http get method error")

	ErrHTTPCtrl       = errors.New("HTTP Control error")
	ErrHashCorrect    = errors.New("wrong hash")
	ErrStructNotValid = errors.New("struct not valid")
	ErrJSONUNMarshal  = errors.New("JSON Marshal failed")

	ErrTestLog = errors.New("log_test error")
)
