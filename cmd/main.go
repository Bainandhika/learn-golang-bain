package main

import (
	"learn-golang-bain/configs"
	"learn-golang-bain/tools"
)

func main() {
	tools.Logger()
	tools.LogInfo.Printf("%#v", configs.GetConfig())
}