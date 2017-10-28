package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const CONFIG_FILE = "./jenxt.json"

type RemoteServer struct {
	Name        string   `json:"name"`
	URL         string   `json:"url"`
	Username    string   `json:"username"`
	PasswordRaw string   `json:"password"`
	Labels      []string `json:"labels"`
}

type Configuration struct {
	Server struct {
		Port       int `json:"port"`
		HostString string
	} `json:"server"`
	Remotes     []RemoteServer `json:"remotes"`
	ServerCache map[string][]*RemoteServer
}

func (c *Configuration) GetServersForLabel(label string) []*RemoteServer {
	if c.ServerCache == nil {
		c.ServerCache = map[string][]*RemoteServer{}
	}

	if servers, ok := c.ServerCache[label]; ok {
		return servers
	}

	for i, s := range c.Remotes {
		for _, serverLabel := range s.Labels {
			if label == serverLabel {
				c.ServerCache[label] = append(c.ServerCache[label], &c.Remotes[i])
				break
			}
		}
	}

	return c.ServerCache[label]
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
