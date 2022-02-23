package main

import (
	"bcm-analyzer/internal/app"
	"fmt"
)

var (
	Version string = "development"
	Build   string = "development"
)

func main() {
	fmt.Println("Version: ", Version)
	fmt.Println("Build time: ", Build)
	app.Run()
}
