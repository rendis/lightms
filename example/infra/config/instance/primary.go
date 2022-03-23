package instance

import (
	"github.com/rendis/lightms"
	"github.com/rendis/lightms/example/infra/config/prop"
	"github.com/rendis/lightms/example/infra/primary"
)

// GetJohnDoeSubscription returns primary.JohnDoeSubscription instance
func GetJohnDoeSubscription() lightms.PrimaryProcess {
	return primary.GetJohnDoeSubscriptionInstance(
		GetJohnDoeUseCase(), prop.GetPubSubProp().Subscriptions.Sub2,
	)
}
