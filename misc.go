package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func StopRedirect(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}

func printRequest(req *http.Request) {
	log.Println(req.Method, req.URL.String(), req.Proto)
	for headerName, headerValues := range req.Header {
		for _, headerValue := range headerValues {
			log.Println(headerName, headerValue)
		}
	}
	log.Println()
}

func printResponse(res *http.Response) {
	if nil == res {
		log.Println("Error: nil response")
		return
	}
	log.Println(res.Proto, res.Status)
	for headerName, headerValues := range res.Header {
		for _, headerValue := range headerValues {
			log.Println(headerName, headerValue)
		}
	}
	log.Println()
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
