package main

import (
	"fmt"
	"jenxt/config"
	"jenxt/scripting"
	"net/http"
)

func main() {
	conf := config.Load()
	scripts := scripting.Scripts{}

	http.HandleFunc("/", scripting.GetHandler(conf, &scripts))
	go scripting.FileWatch(&scripts)

	fmt.Println("Starting Jenxt server on port", conf.Server.Port)
	err := http.ListenAndServe(conf.Server.HostString, nil)
	if err != nil {
		fmt.Println("Server failure: ", err)
	}
}
