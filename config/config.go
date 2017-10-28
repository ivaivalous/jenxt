package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const CONFIG_FILE = "./jenxt.json"

type Configuration struct {
	Server struct {
		Port       int `json:"port"`
		HostString string
	} `json:"server"`
	Remotes []struct {
		Name        string   `json:"name"`
		URL         string   `json:"url"`
		Username    string   `json:"username"`
		PasswordRaw string   `json:"password"`
		Labels      []string `json:"labels"`
	} `json:"remotes"`
}

func (c Configuration) toString() string {
	return toJson(c)
}

func toJson(p interface{}) string {
	bytes, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return string(bytes)
}

func Load() Configuration {
	configuration, err := ioutil.ReadFile(CONFIG_FILE)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	conf := parseConfiguration(configuration)
	conf.Server.HostString = fmt.Sprintf(":%d", conf.Server.Port)
	return conf
}

func parseConfiguration(rawConf []byte) (c Configuration) {
	json.Unmarshal(rawConf, &c)
	return c
}
