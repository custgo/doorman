package main

import (
	"log"
	"net/http"
)

var config Config

func main() {
	config, _ = ReadConfig()
	initHandlers()
	tokenPoolGC()
	err := http.ListenAndServe(config.Listen, nil)
	if nil != err {
		log.Fatalln("Error: ", err)
	}
}
