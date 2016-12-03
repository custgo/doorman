package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/heiing/logs"
)

type LoginRequest struct {
	Method string
	Url    string
	Header http.Header
	Body   string
}

type Config struct {
	Listen         string
	TokenEndpoint  string
	LoginEndpoint  string
	LoginRequestes map[string]LoginRequest
	Logs           *logs.LogsConfig
}

var defaultConfig = Config{
	Listen:         ":6252",
	TokenEndpoint:  "/token",
	LoginEndpoint:  "/doorman-login",
	LoginRequestes: make(map[string]LoginRequest),
	Logs: &logs.LogsConfig{
		Types: []string{"info", "warn", "error"},
		Files: map[string][]string{
			"STDOUT": []string{"info", "warn"},
			"STDERR": []string{"error"},
		},
	},
}

func readConfig() (Config, error) {
	var configFile string
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	} else {
		configFile = filepath.Join(logs.GetExecPath(), "config.json")
	}
	return readConfigFile(configFile)
}

func readConfigFile(configFile string) (Config, error) {
	log.Println("[Config] reading config file:", configFile)
	file, err := os.Open(configFile)
	if err != nil {
		log.Println("[Config] open config file error: ", err)
		return defaultConfig, err
	}
	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Println("[Config] parse config error: ", err)
		os.Exit(1)
	}
	return mergeConfig(config), nil
}

func mergeConfig(config Config) Config {
	if "" == config.Listen {
		config.Listen = defaultConfig.Listen
	}
	if "" == config.LoginEndpoint {
		config.LoginEndpoint = defaultConfig.LoginEndpoint
	}
	if "" == config.TokenEndpoint {
		config.TokenEndpoint = defaultConfig.TokenEndpoint
	}
	return config
}
