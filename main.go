package main

import (
	"fmt"
	"jenxt/config"
	"jenxt/scripting"
)

func main() {
	configuration := config.Load()
	scripts := scripting.Load()

	fmt.Println(configuration)
	fmt.Println(scripts)
}
