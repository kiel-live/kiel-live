package client

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

type consumerEvent struct {
	Stream   string `json:"stream"`
	Consumer string `json:"consumer"`
}

func (n *natsClient) GetSubscribedTopics() []string {
	n.subscriptionsMu.Lock()
	defer n.subscriptionsMu.Unlock()

	subscriptions := []string{}
	for topic := range n.topicSubscriptions {
		subscriptions = append(subscriptions, topic)
	}
	return subscriptions
}

func (n *natsClient) addSubscription(topic, consumerName string) {
	n.subscriptionsMu.Lock()
	defer n.subscriptionsMu.Unlock()

	n.topicSubscriptions[topic] = append(n.topicSubscriptions[topic], consumerName)
}

func (n *natsClient) removeSubscription(topic, consumerName string) {
	n.subscriptionsMu.Lock()
	defer n.subscriptionsMu.Unlock()

	consumers, ok := n.topicSubscriptions[topic]
	if !ok {
		return
	}

	// remove element from array
	for i, consumer := range consumers {
		if consumer == consumerName {
			consumers = append(consumers[:i], consumers[i+1:]...)
			break
		}
	}

	if len(consumers) == 0 {
		delete(n.topicSubscriptions, topic)
	} else {
		n.topicSubscriptions[topic] = consumers
	}
}

func (n *natsClient) init() {
	n.subscriptionsMu.Lock()
	defer n.subscriptionsMu.Unlock()

	// init with already existing consumers
	for consumerInfo := range n.JS.ConsumersInfo("data") {
		topic := consumerInfo.Config.FilterSubject
		n.addSubscription(topic, consumerInfo.Name)
	}

	// new consumer
	err := n.Subscribe("$JS.EVENT.ADVISORY.CONSUMER.CREATED.>", func(msg *Message) {
		var consumerEvent consumerEvent
		if err := json.Unmarshal([]byte(msg.Data), &consumerEvent); err != nil {
			log.Errorf("Parse response failed, reason: %v \n", err)
			return
		}
		consumerInfo, err := n.JS.ConsumerInfo(consumerEvent.Stream, consumerEvent.Consumer)
		if err != nil {
			log.Errorf("Can't find consumer-info: %v", err)
			return
		}
		topic := consumerInfo.Config.FilterSubject

		n.addSubscription(topic, consumerInfo.Name)

		if n.topicSubscriptionHandler != nil {
			n.topicSubscriptionHandler(topic, true)
		}
	})
	if err != nil {
		log.Errorf("Subscribe failed, reason: %v \n", err)
	}

	// remove consumer
	err = n.Subscribe("$JS.EVENT.ADVISORY.CONSUMER.DELETED.>", func(msg *Message) {
		var consumerEvent consumerEvent
		if err := json.Unmarshal([]byte(msg.Data), &consumerEvent); err != nil {
			log.Errorf("Parse response failed, reason: %v \n", err)
			return
		}
		consumerInfo, err := n.JS.ConsumerInfo(consumerEvent.Stream, consumerEvent.Consumer)
		if err != nil {
			log.Errorf("Can't find consumer-info: %v", err)
			return
		}
		topic := consumerInfo.Config.FilterSubject

		n.removeSubscription(topic, consumerEvent.Consumer)

		if n.topicSubscriptionHandler != nil {
			n.topicSubscriptionHandler(topic, false)
		}
	})
	if err != nil {
		log.Errorf("Subscribe failed, reason: %v \n", err)
	}
}
