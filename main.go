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

	server := config.RemoteServer{
		URL:      conf.Remotes[0].URL,
		Username: conf.Remotes[0].Username,
		Password: conf.Remotes[0].PasswordRaw,
	}

	for _, s := range scripts {
		http.HandleFunc(s.GetHandler(server))
		s.PrintInfo()
	}

	fmt.Println("Starting Jenxt server on port ", conf.Server.Port)
	err := http.ListenAndServe(conf.Server.HostString, nil)
	if err != nil {
		fmt.Println("Server failure: ", err)
	}
}
