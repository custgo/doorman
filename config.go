package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
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
}

var DefaultConfig = Config{
	Listen:         ":6252",
	TokenEndpoint:  "/token",
	LoginEndpoint:  "/doorman-login",
	LoginRequestes: make(map[string]LoginRequest),
}

var execPath string

func ReadConfig() (Config, error) {
	var configFile string
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	} else {
		configFile = filepath.Join(GetExecPath(), "config.json")
	}
	return readConfigFile(configFile)
}

func readConfigFile(configFile string) (Config, error) {
	log.Println("[Config] reading config file:", configFile)
	file, err := os.Open(configFile)
	if err != nil {
		log.Println("[Config] open config file error: ", err)
		return DefaultConfig, err
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

func GetExecPath() string {
	if "" == execPath {
		execFile, _ := exec.LookPath(os.Args[0])
		execPath = filepath.Dir(execFile)
	}
	return execPath
}

func mergeConfig(config Config) Config {
	if "" == config.Listen {
		config.Listen = DefaultConfig.Listen
	}
	if "" == config.LoginEndpoint {
		config.LoginEndpoint = DefaultConfig.LoginEndpoint
	}
	if "" == config.TokenEndpoint {
		config.TokenEndpoint = DefaultConfig.TokenEndpoint
	}
	return config
}
