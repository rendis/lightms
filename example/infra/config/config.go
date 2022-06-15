package config

import "github.com/rendis/lightms/example/infra/config/prop"

type InstanceConfig struct {
	DatabaseProp *prop.DatabaseProp
}

type PropConfig struct {
	PubSubProp *prop.PubSubProp
}
