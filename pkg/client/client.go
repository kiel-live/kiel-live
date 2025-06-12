package client

type Client interface {
	Connect() error
	Disconnect() error
	IsConnected() bool
	Subscribe(topic string, cb SubscribeCallback) error
	Unsubscribe(topic string) error
	Publish(topic string, data string) error
	SetConnectionHandler(connectionHandler func(connected bool))
	SetTopicSubscriptionHandler(topicSubscriptionHandler func(topics []string))
}

type TopicMessage struct {
	Topic string `json:"topic,omitempty"`
	Data  string `json:"data,omitempty"`
}
type SubscribeCallback func(msg *TopicMessage)
