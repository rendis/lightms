package prop

type PubSubProp struct {
	PubSub `yaml:"pubsub" json:"pubsub"`
}

type PubSub struct {
	Subscriptions Subscriptions `yaml:"subscriptions" json:"subscriptions"`
}

type Subscriptions struct {
	Sub1 SubscriptionInfo `yaml:"sub1" json:"sub1"`
	Sub2 SubscriptionInfo `yaml:"sub2" json:"sub2"`
}

type SubscriptionInfo struct {
	ProjectId        string `yaml:"project-id" json:"project-id"`
	SubscriptionName string `yaml:"subscription-name" json:"subscription-name"`
}
