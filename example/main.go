package main

import (
	"github.com/rendis/lightms"
	"github.com/rendis/lightms/example/infra/config"
	_ "github.com/rendis/lightms/example/infra/config"
)

func main() {
	// Adding properties
	lightms.AddConf(&config.InstanceConfig{})
	lightms.AddConf(&config.PropConfig{})

	// Set properties file path. Default is "./resources/properties.yml"
	lightms.SetPropFilePath("example/resources/properties.yml")

	// Run lightms
	lightms.Run()
}
