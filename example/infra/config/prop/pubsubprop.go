package prop

import "sync"

var pubSubProp *PubSubProp
var pubSubOnce sync.Once

func GetPubSubProp() *PubSubProp {
	pubSubOnce.Do(func() {
		pubSubProp = &PubSubProp{}
	})
	return pubSubProp
}

type PubSubProp struct {
	PubSub `yaml:"pubsub"`
}

type PubSub struct {
	Subscriptions Subscriptions `yaml:"subscriptions"`
}

type Subscriptions struct {
	Sub1 SubscriptionInfo `yaml:"sub1"`
	Sub2 SubscriptionInfo `yaml:"sub2"`
}

type SubscriptionInfo struct {
	ProjectId        string `yaml:"project-id"`
	SubscriptionName string `yaml:"subscription-name"`
}
