package utils

import (
	"net/url"
	"strings"
)

// Add params to a url string.
func UrlWithParams(url_ string, params map[string]string) string {
	if len(params) == 0 {
		return url_
	}

	if !strings.Contains(url_, "?") {
		url_ += "?"
	}

	if strings.HasSuffix(url_, "?") || strings.HasSuffix(url_, "&") {
		url_ += ParamsToString(params)
	} else {
		url_ += "&" + ParamsToString(params)
	}

	return url_
}

// Convert string map to url component.
func ParamsToString(params map[string]string) string {
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}
	return values.Encode()
}

