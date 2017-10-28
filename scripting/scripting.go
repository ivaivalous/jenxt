package scripting

import (
	"fmt"
	"io/ioutil"
)

const SCRIPTS_LOCATION = "./scripts"

func GetNames() {
	files, err := ioutil.ReadDir(SCRIPTS_LOCATION)
	if err != nil {
		fmt.Println(err.Error())
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}
}
