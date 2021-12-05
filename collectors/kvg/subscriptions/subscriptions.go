package subscriptions

import (
	"encoding/json"
	"sync"

	"github.com/kiel-live/kiel-live/client"
	log "github.com/sirupsen/logrus"
)

type Subscriptions struct {
	client                          *client.Client
	mapConsumer2Subject             map[string]string // for easy deletion
	numberOfSubscriptionsPerSubject map[string]int    // keep track of duplicate subscriptions
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
	for subject := range s.numberOfSubscriptionsPerSubject {
		subscriptions = append(subscriptions, subject)
	}
	return subscriptions
}

func New(client *client.Client) *Subscriptions {
	return &Subscriptions{client: client}
}

func (s *Subscriptions) Subscribe(subscriptionCreatedCallback func()) {
	s.Lock()
	defer s.Unlock()

	s.mapConsumer2Subject = make(map[string]string)
	s.numberOfSubscriptionsPerSubject = make(map[string]int)

	// already existing consumers
	for consumerInfo := range s.client.JS.ConsumersInfo("data") {
		s.mapConsumer2Subject[consumerInfo.Name] = consumerInfo.Config.FilterSubject
		s.numberOfSubscriptionsPerSubject[consumerInfo.Config.FilterSubject]++
	}

	// new consumers
	err := s.client.Subscribe("$JS.EVENT.ADVISORY.CONSUMER.CREATED.>", func(msg *client.SubjectMessage) {
		var consumerEvent consumerEvent
		if err := json.Unmarshal([]byte(msg.Data), &consumerEvent); err != nil {
			log.Fatalf("Parse response failed, reason: %v \n", err)
		}
		consumerInfo, _ := s.client.JS.ConsumerInfo(consumerEvent.Stream, consumerEvent.Consumer)

		s.Lock()
		s.mapConsumer2Subject[consumerInfo.Name] = consumerInfo.Config.FilterSubject
		s.numberOfSubscriptionsPerSubject[consumerInfo.Config.FilterSubject]++
		s.Unlock()

		subscriptionCreatedCallback()
	})
	if err != nil {
		log.Errorf("Subscribe failed, reason: %v \n", err)
	}

	// remove consumers
	err = s.client.Subscribe("$JS.EVENT.ADVISORY.CONSUMER.DELETED.>", func(msg *client.SubjectMessage) {
		var consumerEvent consumerEvent
		if err := json.Unmarshal([]byte(msg.Data), &consumerEvent); err != nil {
			log.Fatalf("Parse response failed, reason: %v \n", err)
		}

		s.Lock()
		defer s.Unlock()

		if s.numberOfSubscriptionsPerSubject[s.mapConsumer2Subject[consumerEvent.Consumer]] > 1 {
			s.numberOfSubscriptionsPerSubject[s.mapConsumer2Subject[consumerEvent.Consumer]]--
		} else {
			delete(s.numberOfSubscriptionsPerSubject, s.mapConsumer2Subject[consumerEvent.Consumer])
		}
		delete(s.mapConsumer2Subject, consumerEvent.Consumer)
	})
	if err != nil {
		log.Errorf("Subscribe failed, reason: %v \n", err)
	}
}
