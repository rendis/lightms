package config

import (
	"github.com/rendis/lightms/example/core/impl"
	"github.com/rendis/lightms/example/core/port"
	"github.com/rendis/lightms/example/core/usecase"
	"github.com/rendis/lightms/example/infra/config/prop"
	"log"
)

func (c *InstanceConfig) JohnDoeImpl(pport port.PersistencePort, dbProp *prop.DatabaseProp) usecase.JohnDoeUseCase {
	log.Printf("PersistencePort and DatabaseProp received: %+v, %+v", pport, dbProp)
	return impl.NewJohnDoeImpl(pport)
}
