package main

import (
	"fmt"
	"jenxt/config"
	"jenxt/scripting"
	"net/http"
)

func main() {
	configuration := config.Load()
	scripts := scripting.Load()

	for _, s := range scripts {
		http.HandleFunc(s.GetHandler())
	}

	err := http.ListenAndServe(configuration.Server.HostString, nil)
	if err != nil {
		fmt.Println("Server failure: ", err)
	}
}
