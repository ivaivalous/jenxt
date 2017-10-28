package main

import (
	"fmt"
	"jenxt/config"
	"jenxt/scripting"
)

func main() {
	var conf config.Configuration
	conf = config.LoadConfig()

	fmt.Println(conf.Server.Host)

	scripting.GetNames()
}
