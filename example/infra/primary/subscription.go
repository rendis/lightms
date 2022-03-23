package primary

import (
	"github.com/rendis/lightms/example/core/usecase"
	"github.com/rendis/lightms/example/infra/config/prop"
	"log"
	"strconv"
	"sync"
)

var (
	subscription *JohnDoeSubscription
	once         sync.Once
)

func GetJohnDoeSubscriptionInstance(u usecase.JohnDoeUseCase, p prop.SubscriptionInfo) *JohnDoeSubscription {
	once.Do(func() {
		subscription = &JohnDoeSubscription{u, p}
	})
	return subscription
}

type JohnDoeSubscription struct {
	useCase usecase.JohnDoeUseCase
	prop    prop.SubscriptionInfo
}

// HandleSubEvent handles the event
func (e *JohnDoeSubscription) HandleSubEvent(event string) {
	err := e.useCase.Handle(event)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("event trigger success")
	}
}

// Start implementing the interface lightms.PrimaryProcess
func (e *JohnDoeSubscription) Start() {
	log.Printf("Starting John Doe subscription to '%s'\n", e.prop.SubscriptionName)
	for i := 0; i < 10; i++ {
		e.HandleSubEvent("John Doe event # " + strconv.Itoa(i))
	}
}
