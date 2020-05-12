/**
 totoroAgent project config.go
 author:feng
 since:2020-01-08
**/
package totoroAgent

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"
)

var AppConfig Config

type Config struct {
	LogPath    string `json:"logPath"`
	ListenPort string `json:"listenPort"`
	PidPath    string `json:"pidPath"`
	Version    string `json:"version"`
}

func loadConfig(configFile string) (config Config, err error) {
	file, e := ioutil.ReadFile(configFile)
	if e != nil {
		log.Fatal("File error: %v\n", e)
		return config, errors.New("read file error" + e.Error())
	}

	err1 := json.Unmarshal(file, &config)
	if err1 != nil {
		log.Fatal("File error: %v\n", err1)
		return config, errors.New("parse config error" + err1.Error())
	}

	return config, nil
}

func InitConfig() (config Config, s string, err error) {
	conf := flag.String("c", "NONE", "config file")
	sig := flag.String("s", "", "signal")

	flag.Parse()
	s = *sig

	if *conf != "NONE" {
		log.Println("config file: ", *conf)
		config, e := loadConfig(*conf)
		if e == nil {
			return config, s, nil
		}
	} else {
		log.Println("Usage: totoroAgent -c /xxx/config.json to run it with config file")
	}

	// loading default paramter
	config.LogPath = "/var/Logs/totoroAgent/agent.log"
	config.ListenPort = ":10099"
	config.PidPath = "/var/Data/totoroAgent/"
	config.Version = "1.0"

	return config, s, nil
}
