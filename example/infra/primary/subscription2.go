package primary

import (
	"github.com/rendis/lightms/v3"
	"github.com/rendis/lightms/v3/example/core/usecase"
	"github.com/rendis/lightms/v3/example/infra/config/prop"
	"log"
	"strconv"
	"time"
)

func NewJaneDoeSubscription(johnDoeUseCase usecase.JohnDoeUseCase, psProp *prop.PubSubProp) lightms.PrimaryProcess {
	return &JaneDoeSubscription{johnDoeUseCase, psProp.Subscriptions.Sub1}
}

type JaneDoeSubscription struct {
	useCase usecase.JohnDoeUseCase
	prop    prop.SubscriptionInfo
}

// HandleSubEvent handles the event
func (e *JaneDoeSubscription) HandleSubEvent(event string) {
	log.Printf("John Doe subscription handling event '%s'\n", event)
	err := e.useCase.Handle(event)
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("Event '%s' handled.\n\n", event)
	}
}

// Start implementing the interface lightms.PrimaryProcess
func (e *JaneDoeSubscription) Start() {
	log.Printf("Starting John Doe subscription to '%s'\n\n", e.prop.SubscriptionName)
	var i int
	for {
		i++
		time.Sleep(1 * time.Second)
		e.HandleSubEvent("John Doe event # " + strconv.Itoa(i))
	}
}
