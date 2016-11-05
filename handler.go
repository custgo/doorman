package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
)

func LoginHandler(rw http.ResponseWriter, req *http.Request) {
	token := req.FormValue("token")
	identity := tokenPool.Get(token)
	if nil == identity {
		rw.WriteHeader(401)
		return
	}

	hostConfig, exists := config.LoginRequestes[req.Host]
	if !exists {
		log.Println("host config NOT FOUND:", req.Host)
		rw.WriteHeader(404)
		return
	}

	username := url.QueryEscape(identity.username)
	password := url.QueryEscape(identity.password)
	redirectUrl := url.QueryEscape(req.FormValue("url"))

	targetUrl := parseTemplate(hostConfig.Url, username, password, redirectUrl)
	body := parseTemplate(hostConfig.Body, username, password, redirectUrl)

	request, _ := http.NewRequest(hostConfig.Method, targetUrl, strings.NewReader(body))
	for headerName, headerValues := range hostConfig.Header {
		for _, headerValue := range headerValues {
			request.Header.Add(headerName, headerValue)
		}
	}

	client := &http.Client{
		CheckRedirect: StopRedirect,
	}
	res, err := client.Do(request)
	if err != nil {
		log.Println("[HTTP] Login Error, Do Request Error: ", err)
		return
	}

	for headerName, headerValues := range res.Header {
		for _, headerValue := range headerValues {
			rw.Header().Add(headerName, headerValue)
		}
	}

	rw.WriteHeader(res.StatusCode)
	io.Copy(rw, res.Body)
}

func TokenHandler(rw http.ResponseWriter, req *http.Request) {
	printRequest(req)
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
