package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/kiel-live/kiel-live/packages/client"
	"github.com/kiel-live/kiel-live/packages/collector-kvg/api"
)

type collector struct {
	client      *client.WebSocketClient
	channel     string
	channelType string
	provider    string
	entityID    string
}

func newCollector(client *client.WebSocketClient, channel string) (*collector, error) {
	parts := strings.Split(strings.TrimLeft(channel, "/"), "/")
	if len(parts) < 1 {
		return nil, errors.New("Channel format not supported")
	}

	channelType := parts[0]

	if channelType != "stops" && channelType != "vehicles" && channelType != "stop" && channelType != "vehicle" && channelType != "trip" && channelType != "route" {
		return nil, errors.New("Channel type not supported")
	}

	if channelType != "stops" && channelType != "vehicles" && len(parts) != 3 {
		return nil, errors.New("Channel format not supported")
	}

	provider := ""
	entityID := ""

	if len(parts) == 3 {
		provider = parts[1]
		entityID = parts[2]
	}

	return &collector{
		client:      client,
		channel:     channel,
		channelType: channelType,
		provider:    provider,
		entityID:    entityID,
	}, nil
}

func (c *collector) collect() interface{} {
	if c.channelType == "stops" {
		return api.GetStops()
	}

	if c.channelType == "stop" {
		return api.GetStop(c.entityID)
	}

	if c.channelType == "vehicles" {
		return api.GetVehicles()
	}

	if c.channelType == "vehicle" {
		return api.GetVehicle(c.entityID)
	}

	if c.channelType == "trip" {
		return api.GetTrip(c.entityID)
	}

	if c.channelType == "route" {
		return api.GetRoute(c.entityID)
	}

	return nil
}

func (c *collector) publish(data interface{}) error {
	encoded, err := json.Marshal(data)
	if err != nil {
		return err
	}

	c.client.Publish(c.channel, string(encoded))
	return nil
}

func (c *collector) run() {
	data := c.collect()
	err := c.publish(data)

	if err != nil {
		fmt.Printf("Error while publishing to channel '%s': %v", c.channel, err)
	}
}
