package client

import "github.com/kiel-live/kiel-live/protocol"

type Client interface {
	Connect() error
	Disconnect() error
	IsConnected() bool
	SetOnConnectionChanged(handler func(connected bool))

	GetSubscribedTopics() []string
	SetOnTopicsChanged(handler func(topic string, added bool))

	UpdateStop(stop *protocol.Stop) error
	UpdateVehicle(vehicle *protocol.Vehicle) error
	UpdateTrip(trip *protocol.Trip) error
	DeleteStop(stopID string) error
	DeleteVehicle(vehicleID string) error
	DeleteTrip(tripID string) error
}

type Message struct {
	Topic string `json:"topic,omitempty"`
	Data  string `json:"data,omitempty"`
}

type SubscribeCallback func(msg *Message)
