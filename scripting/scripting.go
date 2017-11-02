package scripting

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
)

const (
	SCRIPTS_LOCATION    = "./scripts"
	META_REGEXP         = "<jenxt>(?P<Meta>[\\S\\s]*)</jenxt>"
	MISSING_META_ERR    = "Invalid script - Meta is missing"
	BAD_META_ERR        = "Ignored script %s: Meta is malformed: %s"
	IGNORED_MISSING_ERR = "Ignored script %s: %s"
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
}

// Load reads all available scripts and attempts to read their
// meta information. If parsing this information fails for a file,
// it is ignored. A message for information is then printed to the console.
func Load() map[string]Meta {
	scripts := make(map[string]Meta)

	files, err := ioutil.ReadDir(SCRIPTS_LOCATION)
	if err != nil {
		fmt.Println(err.Error())
	}

	for _, f := range files {
		content := read(f.Name())

		metaRaw, err := extractMeta(content)
		if err != nil {
			fmt.Println(fmt.Sprintf(IGNORED_MISSING_ERR, f.Name(), err.Error()))
			continue
		}

		meta, err := parseMeta(metaRaw)
		if err != nil {
			fmt.Println(fmt.Sprintf(BAD_META_ERR, f.Name(), err.Error()))
			continue
		}

		meta.Script = content
		scripts[f.Name()] = meta
	}

	return scripts
}

// read opens a file and returns its contents as a string
func read(filename string) string {
	raw, err := ioutil.ReadFile(path.Join(SCRIPTS_LOCATION, filename))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return string(raw[:])
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
