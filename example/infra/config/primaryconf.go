package config

import (
	"github.com/rendis/lightms"
	"github.com/rendis/lightms/example/core/usecase"
	"github.com/rendis/lightms/example/infra/config/prop"
	"github.com/rendis/lightms/example/infra/primary"
)

func (c *InstanceConfig) JohnDoeSubscription(uc usecase.JohnDoeUseCase, psProp *prop.PubSubProp) lightms.PrimaryProcess {
	return primary.NewJohnDoeSubscription(uc, psProp.Subscriptions.Sub2)
}
