package scripting

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
)

const (
	SCRIPTS_LOCATION = "./scripts"
	META_REGEXP      = "<jenxt>(?P<Meta>[\\S\\s]*)</jenxt>"
)

func GetNames() {
	files, err := ioutil.ReadDir(SCRIPTS_LOCATION)
	if err != nil {
		fmt.Println(err.Error())
	}

	for _, f := range files {
		fmt.Println(f.Name())
		content := read(f.Name())
		meta := extractMeta(content)

		fmt.Println(meta)
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

func extractMeta(script string) string {
	re := regexp.MustCompile(META_REGEXP)
	res := re.FindStringSubmatch(script)

	fmt.Printf("%v", res[1])
	return ""
}
