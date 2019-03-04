package scripting /* import "ivo.qa/jenxt/scripting" */

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
	FileWatchIntervalSeconds = 10
	ScriptsLocation          = "./scripts"
	MetaRegexp               = "<jenxt>(?P<Meta>[\\S\\s]*)</jenxt>"
	MissingMetaError         = "invalid script - Meta is missing"
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

type Scripts map[string]Meta

// Load reads all available scripts and attempts to read their
// meta information. If parsing this information fails for a file,
// it is ignored. A message for information is then printed to the console.
func (currentScripts *Scripts) Load() {
	newScripts := Scripts{}

	files, err := ioutil.ReadDir(ScriptsLocation)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, f := range files {
		content, err := read(f.Name())
		if err != nil {
			fmt.Println("File", f.Name(), "could not be read. Removing.")
			continue
		}

		// File has already been loaded, check for updates
		if existingMeta, ok := (*currentScripts)[f.Name()]; ok {
			if config.GetFileHash(content) == existingMeta.Hash {
				newScripts[f.Name()] = existingMeta
				continue
			}

			newMeta, err := LoadContent(existingMeta.FileName, content)
			if err != nil {
				fmt.Println("Change detected for", existingMeta.getEndpoint(), "but reload failed")
				continue
			}

			fmt.Println("Resource", existingMeta.getEndpoint(), "reloaded due to file change")
			newScripts[f.Name()] = newMeta
			continue
		}

		// New file
		meta, err := LoadContent(f.Name(), content)
		if err != nil {
			fmt.Println("Could not load", meta.FileName)
			continue
		}

		meta.PrintInfo()
		newScripts[f.Name()] = meta
	}

	*currentScripts = newScripts
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
func FileWatch(scripts *Scripts) {
	for {
		scripts.Load()
		time.Sleep(FileWatchIntervalSeconds * time.Second)
	}
}

// read opens a file and returns its contents as a string
func read(filename string) (content string, err error) {
	raw, err := ioutil.ReadFile(path.Join(ScriptsLocation, filename))
	if err != nil {
		return
	}

	return string(raw[:]), nil
}

// extractMeta reads a script's meta information from a Groovy script
func extractMeta(script string) (result string, err error) {
	re := regexp.MustCompile(MetaRegexp)
	res := re.FindStringSubmatch(script)

	if len(res) != 2 {
		err = errors.New(MissingMetaError)
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
