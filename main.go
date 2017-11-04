package main

import (
	"fmt"
	"jenxt/config"
	"jenxt/scripting"
	"net/http"
)

func main() {
	conf := config.Load()
	scripts := scripting.Load()

	for _, s := range scripts {
		http.HandleFunc(s.GetHandler(conf))
		s.PrintInfo()
	}

	fmt.Println("Starting Jenxt server on port", conf.Server.Port)
	err := http.ListenAndServe(conf.Server.HostString, nil)
	if err != nil {
		fmt.Println("Server failure: ", err)
	}
}
