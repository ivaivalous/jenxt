// Package config is meant to handle the application's configuration
package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const SettingsFile = "./jenxt.json"

// RemoteServer describes a Jenkins server Jenxt can connect to
type RemoteServer struct {
	Name        string   `json:"name"`
	URL         string   `json:"url"`
	Username    string   `json:"username"`
	PasswordRaw string   `json:"password"`
	Labels      []string `json:"labels"`
}

// Configuration describes the complete configuration of the Jenxt server
type Configuration struct {
	Server struct {
		Port       int `json:"port"`
		HostString string
	} `json:"server"`
	Remotes     []RemoteServer `json:"remotes"`
	ServerCache map[string][]*RemoteServer
}

// GetServersForLabel returns a list of pointers to all servers
// that have been labeled with a particular label.
// This operation utilizes a cache. The first attempt to get
// servers for a label will be slower as the cache is built.
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

// Load reads configuration from a file (jenxt.json)
// and parses it into a Configuration struct
func Load() Configuration {
	configuration, err := ioutil.ReadFile(SettingsFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	conf := parseConfiguration(configuration)
	conf.Server.HostString = fmt.Sprintf(":%d", conf.Server.Port)
	return conf
}

// parseConfiguration reads the contents of a configuration file
// into a Configuration object
func parseConfiguration(rawConf []byte) (c Configuration) {
	json.Unmarshal(rawConf, &c)
	return c
}
