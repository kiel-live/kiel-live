package client

import (
	"github.com/kiel-live/kiel-live/pkg/models"
)

type Client interface {
	Connect() error
	Disconnect() error
	IsConnected() bool
	SetOnConnectionChanged(handler func(connected bool))

	GetSubscribedTopics() []string
	SetOnTopicsChanged(handler func(topic string, subscribed bool))

	UpdateStop(stop *models.Stop) error
	UpdateVehicle(vehicle *models.Vehicle) error
	UpdateTrip(trip *models.Trip) error
	DeleteStop(stopID string) error
	DeleteVehicle(vehicleID string) error
	DeleteTrip(tripID string) error
}

type Message struct {
	Topic string `json:"topic,omitempty"`
	Data  string `json:"data,omitempty"`
}

type SubscribeCallback func(msg *Message)
