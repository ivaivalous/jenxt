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

type Meta struct {
	Expose         string   `json:"expose"`            // The URL at which to expose the script
	Methods        []string `json:"methods,omitempty"` // The HTTP methods available for the script
	Authentication string   `json:"authentication"`    // The auth type used for the script
	Params         []struct {
		Name  string `json:"name"`            // The name of the parameter
		Regex string `json:"regex,omitempty"` // A regular expression the parameter must adhere to

	} `json:"params,omitempty"` // Expected URL parameters
}

func GetNames() {
	files, err := ioutil.ReadDir(SCRIPTS_LOCATION)
	if err != nil {
		fmt.Println(err.Error())
	}

	for _, f := range files {
		fmt.Println(f.Name())
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

		fmt.Println(fmt.Sprintf("%s", meta.Expose))
	}
}

func read(filename string) string {
	raw, err := ioutil.ReadFile(path.Join(SCRIPTS_LOCATION, filename))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return string(raw[:])
}

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

func parseMeta(metaRaw string) (meta Meta, err error) {
	err = json.Unmarshal([]byte(metaRaw), &meta)
	return meta, err
}
