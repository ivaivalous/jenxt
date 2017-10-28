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
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"server"`
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

func LoadConfig() Configuration {
	configuration, err := ioutil.ReadFile(CONFIG_FILE)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return parseConfiguration(configuration)
}

func parseConfiguration(rawConf []byte) (c Configuration) {
	json.Unmarshal(rawConf, &c)
	return c
}
