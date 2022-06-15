package config

import (
	"github.com/rendis/lightms/example/core/port"
	"github.com/rendis/lightms/example/infra/secondary"
)

func (c *InstanceConfig) PersistencePortImpl() port.PersistencePort {
	return secondary.NewPersistencePort(c.DatabaseProp.Postgresql)
}
