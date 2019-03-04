// Package config is meant to handle the application's configuration
package config /* import "ivo.qa/jenxt/config" */

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

// Secret defines a server name-password mapping stored
// in a separate file to the main configuration.
type Secret struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// SecretFile defines the structure of a JSON file
// storing server passwords
type SecretFile struct {
	Remotes []Secret `json:"remotes"`
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

// MergePasswords adds the passwords found in a secrets.json
// file into a Configuration object.
func (c *Configuration) MergePasswords(s SecretFile) {
	for _, secretServer := range s.Remotes {
		for i, server := range c.Remotes {
			if secretServer.Name == server.Name {
				c.Remotes[i].PasswordRaw = secretServer.Password
			}
		}
	}
}

// Load reads configuration from a file (jenxt.json)
// and parses it into a Configuration struct
func Load(settingsFilePath string) Configuration {
	configuration, err := ioutil.ReadFile(settingsFilePath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	conf := parseConfiguration(configuration)
	conf.Server.HostString = fmt.Sprintf(":%d", conf.Server.Port)
	return conf
}

// Load reads configuration from a the default location (jenxt.json)
// and parses it into a Configuration struct
func LoadDefault() Configuration {
	return Load(SettingsFile)
}

// LoadSecrets reads a file and parses it
// into a SecretFile
func LoadSecrets(secretsFilePath string) SecretFile {
	secretsConf, err := ioutil.ReadFile(secretsFilePath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return parseSecrets(secretsConf)
}

// LoadDynamic checks if the user has requested to load the configuration
// from a non-default location and if so, loads it.
// Otherwise it uses the default location.
func LoadDynamic() Configuration {
	if args := os.Args[1:]; len(args) >= 1 {
		conf := Load(args[0])

		if len(args) >= 2 {
			// A secrets file has been configured, too
			secrets := LoadSecrets(args[1])
			conf.MergePasswords(secrets)
		}

		return conf
	}

	return LoadDefault()
}

// parseConfiguration reads the contents of a configuration file
// into a Configuration object
func parseConfiguration(rawConf []byte) (c Configuration) {
	json.Unmarshal(rawConf, &c)
	return c
}

// parseSecrets reads a configuration file of server passwords
// into a SecretFile object
func parseSecrets(rawConf []byte) (s SecretFile) {
	json.Unmarshal(rawConf, &s)
	return s
}
