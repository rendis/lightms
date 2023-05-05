package main

import (
	"github.com/rendis/lightms/v3"
	"github.com/rendis/lightms/v3/example/infra/config"
	_ "github.com/rendis/lightms/v3/example/infra/config"
	_ "github.com/rendis/lightms/v3/example/infra/secondary"
)

func main() {
	// Adding properties
	lightms.AddConf[config.PropsConfig]()

	//Set properties file path. Default is "./resources/properties.yml"
	lightms.SetPropFilePath("example/resources/properties.yml")
	//lightms.SetPropFilePath("example/resources/properties.json")

	// Run lightms
	lightms.Run()
}
