package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func StopRedirect(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}

func parseTemplate(tmpl, user, pass, url string) string {
	ret := strings.Replace(tmpl, "{username}", user, -1)
	ret = strings.Replace(ret, "{password}", pass, -1)
	ret = strings.Replace(ret, "{redirectUrl}", url, -1)
	return ret
}

func readString(r io.Reader) string {
	if nil == r {
		return ""
	}
	b, _ := ioutil.ReadAll(r)
	return string(b)
}
