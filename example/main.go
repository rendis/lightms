package main

import (
	"github.com/rendis/lightms"
	"github.com/rendis/lightms/example/infra/config/instance"
	"github.com/rendis/lightms/example/infra/config/prop"
)

func main() {
	// Adding properties
	lightms.AddProperty(prop.GetDatabaseProp())
	lightms.AddProperty(prop.GetPubSubProp())

	// Adding primary processes
	lightms.AddPrimary(instance.GetJohnDoeSubscription)

	// Set properties file path. Default is "./resources/properties.yml"
	lightms.SetPropFilePath("example/resources/properties.yml")

	// Run lightms
	lightms.Run()
}
