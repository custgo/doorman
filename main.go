package main

import (
	"net/http"
	"strings"

	"github.com/heiing/logs"
)

var config Config

func main() {
	config, _ = readConfig()

	if config.Logs != nil {
		logs.SetDefaultLoggerForConfig(config.Logs)
		logs.Info("[Logger] ", strings.Join(config.Logs.Types, ","))
		for file, types := range config.Logs.Files {
			logs.Info("[Logger] ", file, " - ", strings.Join(types, ","))
		}
	} else {
		logs.Warn("[Logger] log config not found, use DefaultLogger.")
	}
	initHandlers()
	tokenPoolGC()
	err := http.ListenAndServe(config.Listen, nil)
	if nil != err {
		logs.Error("Error: ", err)
	}
}
