package main

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	"github.com/heiing/logs"
)

func LoginHandler(rw http.ResponseWriter, req *http.Request) {
	token := req.FormValue("token")
	identity := tokenPool.Get(token)
	if nil == identity {
		rw.WriteHeader(401)
		logs.Info("[Access] token ", token, " not found, return 401")
		return
	}

	hostConfig, exists := config.LoginRequestes[req.Host]
	if !exists {
		rw.WriteHeader(404)
		logs.Info("[Access] host config NOT FOUND: ", req.Host, ", return 404")
		return
	}

	username := url.QueryEscape(identity.username)
	password := url.QueryEscape(identity.password)
	redirectUrl := url.QueryEscape(req.FormValue("url"))

	targetUrl := parseTemplate(hostConfig.Url, username, password, redirectUrl)
	body := parseTemplate(hostConfig.Body, username, password, redirectUrl)

	request, _ := http.NewRequest(hostConfig.Method, targetUrl, strings.NewReader(body))
	logs.Debug("[HTTP][Request] ", hostConfig.Method, " ", targetUrl, " ", request.Proto)
	for headerName, headerValues := range hostConfig.Header {
		for _, headerValue := range headerValues {
			request.Header.Add(headerName, headerValue)
			logs.Debug("[HTTP][Request] ", headerName, ": ", headerValue)
		}
	}
	logs.Debug("[HTTP][Request] ", body)

	client := &http.Client{
		CheckRedirect: StopRedirect,
	}
	res, err := client.Do(request)
	if err != nil {
		logs.Error("[HTTP][Request] Login Error, Do Request Error: ", err)
		return
	}

	logs.Debug("[HTTP][Response] ", res.Proto, " ", res.Status)
	for headerName, headerValues := range res.Header {
		for _, headerValue := range headerValues {
			rw.Header().Add(headerName, headerValue)
			logs.Debug("[HTTP][Response] ", headerName, ": ", headerValue)
		}
	}

	rw.WriteHeader(res.StatusCode)
	resBody := readString(res.Body)
	rw.Write([]byte(resBody))
	logs.Debug("[HTTP][Response] ", resBody)
}

// 每次请求，都会生成不一样的 token
func TokenHandler(rw http.ResponseWriter, req *http.Request) {
	username := req.PostFormValue("username")
	password := req.PostFormValue("password")
	token := tokenPool.Add(username, password)
	rw.WriteHeader(201)
	_, _ = io.WriteString(rw, token)
}

func TestLoginHandler(rw http.ResponseWriter, req *http.Request) {
	username := req.FormValue("username")
	url := req.FormValue("url")
	if "" == url {
		url = "/"
	}
	rw.Header().Set("Location", url)
	rw.Header().Add("Set-Cookie", "SESSID=loginOK")
	rw.WriteHeader(302)
	io.WriteString(rw, "User "+username+" Login OK")
}

func initHandlers() {
	r := mux.NewRouter()
	r.HandleFunc(config.LoginEndpoint, LoginHandler).Methods("GET")
	r.HandleFunc(config.TokenEndpoint, TokenHandler).Methods("POST")
	r.HandleFunc("/testLogin", TestLoginHandler).Methods("GET", "POST", "PUT")
	http.Handle("/", r)
}
