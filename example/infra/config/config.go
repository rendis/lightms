package config

import (
	"github.com/rendis/lightms"
	"github.com/rendis/lightms/example/core/impl"
	"github.com/rendis/lightms/example/core/port"
	"github.com/rendis/lightms/example/infra/config/prop"
	"github.com/rendis/lightms/example/infra/primary"
	"github.com/rendis/lightms/example/infra/secondary"
)

type PropsConfig struct {
	DatabaseProp *prop.DatabaseProp
	PubSubProp   *prop.PubSubProp
}

func (c *PropsConfig) Subs() {
	lightms.AddInstance(primary.NewJohnDoeSubscription)
}

var _ = lightms.AddInstance(primary.NewJaneDoeSubscription)

//var _ = lightms.AddInstance(primary.NewJohnDoeSubscription)

var _ = lightms.AddInstance(impl.NewJohnDoeImpl).
	AndInjections().
	WithInjection("postgres.PersistencePort", lightms.TOF[port.PersistencePort]())

var _ = lightms.AddInstance(secondary.NewPersistence1Port).
	WithAlias("postgres.PersistencePort", lightms.TOF[port.PersistencePort]())
