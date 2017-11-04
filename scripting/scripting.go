package scripting

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"jenxt/config"
	"path"
	"regexp"
	"strings"
	"time"
)

const (
	FILE_WATCH_INTERVAL_S = 10
	SCRIPTS_LOCATION      = "./scripts"
	META_REGEXP           = "<jenxt>(?P<Meta>[\\S\\s]*)</jenxt>"
	MISSING_META_ERR      = "Invalid script - Meta is missing"
	BAD_META_ERR          = "Ignored script %s: Meta is malformed: %s"
	IGNORED_MISSING_ERR   = "Ignored script %s: %s"
)

// Meta represents the system information given for configuration
// at the beginning of every Groovy script file.
type Meta struct {
	Expose         string `json:"expose"`         // The URL at which to expose the script
	Authentication string `json:"authentication"` // The auth type used for the script
	Params         []struct {
		Name  string `json:"name"`            // The name of the parameter
		Regex string `json:"regex,omitempty"` // A regular expression the parameter must adhere to

	} `json:"params,omitempty"` // Expected URL parameters
	JSONResponse bool   `json:"jsonResponse"` // Whether to return the Jenkins response as JSON or bool
	Script       string // The content of the script
	Hash         string // A hash of the file's contents
	FileName     string // The name of the file the script was loaded from
}

// Load reads all available scripts and attempts to read their
// meta information. If parsing this information fails for a file,
// it is ignored. A message for information is then printed to the console.
func Load() map[string]*Meta {
	scripts := make(map[string]*Meta)

	files, err := ioutil.ReadDir(SCRIPTS_LOCATION)
	if err != nil {
		fmt.Println(err.Error())
		return scripts
	}

	for _, f := range files {
		meta, err := LoadFile(f.Name())
		if err != nil {
			fmt.Println(fmt.Sprintf(BAD_META_ERR, f.Name(), err.Error()))
			continue
		}

		scripts[f.Name()] = &meta
	}

	return scripts
}

// Reload reads through loaded scripts and reloads ones that
// have been changed on disk.
func Reload(currentScripts map[string]*Meta) {
	for name, meta := range currentScripts {
		content, err := read(meta.FileName)
		if err != nil {
			fmt.Println("Could not load file", meta.FileName)
			continue
		}

		hash := config.GetFileHash(content)

		if hash != meta.Hash {
			newMeta, err := LoadContent(meta.FileName, content)
			if err != nil {
				fmt.Println("Change detected for", meta.FileName, "but reload failed")
				continue
			}

			currentScripts[name] = &newMeta
			fmt.Println("Script", meta.FileName, "updated due to file change")
		}
	}
}

// LoadFile reads and parses a script file
func LoadFile(filename string) (meta Meta, err error) {
	content, err := read(filename)
	if err != nil {
		return
	}

	return LoadContent(filename, content)
}

// LoadContent parses a string into a Meta object
func LoadContent(filename, content string) (meta Meta, err error) {
	metaRaw, err := extractMeta(content)
	if err != nil {
		return
	}

	meta, err = parseMeta(metaRaw)
	if err != nil {
		return
	}

	meta.Script = content
	meta.Hash = config.GetFileHash(content)
	meta.FileName = filename

	return
}

// FileWatch runs forever, checking for sript file
// changes. It should be called as a goroutine.
func FileWatch(scripts map[string]*Meta) {
	for {
		Reload(scripts)
		time.Sleep(FILE_WATCH_INTERVAL_S * time.Second)
	}
}

// read opens a file and returns its contents as a string
func read(filename string) (content string, err error) {
	raw, err := ioutil.ReadFile(path.Join(SCRIPTS_LOCATION, filename))
	if err != nil {
		return
	}

	return string(raw[:]), nil
}

// extractMeta reads a script's meta information from a Groovy script
func extractMeta(script string) (result string, err error) {
	re := regexp.MustCompile(META_REGEXP)
	res := re.FindStringSubmatch(script)

	if len(res) != 2 {
		err = errors.New(MISSING_META_ERR)
		return
	}

	result = strings.TrimSpace(res[1])
	return
}

// parseMeta reads a Meta formatted string into an actual Meta object.
// It's usually fed with the output form extractMeta.
func parseMeta(metaRaw string) (meta Meta, err error) {
	err = json.Unmarshal([]byte(metaRaw), &meta)
	return meta, err
}
