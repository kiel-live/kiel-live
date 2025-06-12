package subscriptions

import (
	"encoding/json"
	"sync"

	"github.com/kiel-live/kiel-live/client"
	log "github.com/sirupsen/logrus"
)

type Subscriptions struct {
	client                        *client.Client
	mapConsumer2Topic             map[string]string // for easy deletion
	numberOfSubscriptionsPerTopic map[string]int    // keep track of duplicate subscriptions
	sync.Mutex
}

type consumerEvent struct {
	Stream   string `json:"stream"`
	Consumer string `json:"consumer"`
}

func (s *Subscriptions) GetSubscriptions() []string {
	subscriptions := []string{}
	s.Lock()
	defer s.Unlock()
	for topic := range s.numberOfSubscriptionsPerTopic {
		subscriptions = append(subscriptions, topic)
	}
	return subscriptions
}

func New(client *client.Client) *Subscriptions {
	return &Subscriptions{client: client}
}

func (s *Subscriptions) Subscribe(subscriptionCreatedCallback func(topic string)) {
	s.Lock()
	defer s.Unlock()

	s.mapConsumer2Topic = make(map[string]string)
	s.numberOfSubscriptionsPerTopic = make(map[string]int)

	// already existing consumers
	for consumerInfo := range s.client.JS.ConsumersInfo("data") {
		s.mapConsumer2Topic[consumerInfo.Name] = consumerInfo.Config.FilterSubject
		s.numberOfSubscriptionsPerTopic[consumerInfo.Config.FilterSubject]++
	}

	// new consumers
	err := s.client.Subscribe("$JS.EVENT.ADVISORY.CONSUMER.CREATED.>", func(msg *client.Message) {
		var consumerEvent consumerEvent
		if err := json.Unmarshal([]byte(msg.Data), &consumerEvent); err != nil {
			log.Errorf("Parse response failed, reason: %v \n", err)
			return
		}
		consumerInfo, err := s.client.JS.ConsumerInfo(consumerEvent.Stream, consumerEvent.Consumer)
		if err != nil {
			log.Errorf("Can't find consumer-info: %v", err)
			return
		}

		s.Lock()
		s.mapConsumer2Topic[consumerInfo.Name] = consumerInfo.Config.FilterSubject
		s.numberOfSubscriptionsPerTopic[consumerInfo.Config.FilterSubject]++
		s.Unlock()

		subscriptionCreatedCallback(consumerInfo.Config.FilterSubject)
	})
	if err != nil {
		log.Errorf("Subscribe failed, reason: %v \n", err)
	}

	// remove consumers
	err = s.client.Subscribe("$JS.EVENT.ADVISORY.CONSUMER.DELETED.>", func(msg *client.Message) {
		var consumerEvent consumerEvent
		if err := json.Unmarshal([]byte(msg.Data), &consumerEvent); err != nil {
			log.Errorf("Parse response failed, reason: %v \n", err)
			return
		}

		s.Lock()
		defer s.Unlock()

		if s.numberOfSubscriptionsPerTopic[s.mapConsumer2Topic[consumerEvent.Consumer]] > 1 {
			s.numberOfSubscriptionsPerTopic[s.mapConsumer2Topic[consumerEvent.Consumer]]--
		} else {
			delete(s.numberOfSubscriptionsPerTopic, s.mapConsumer2Topic[consumerEvent.Consumer])
		}
		delete(s.mapConsumer2Topic, consumerEvent.Consumer)
	})
	if err != nil {
		log.Errorf("Subscribe failed, reason: %v \n", err)
	}
}
