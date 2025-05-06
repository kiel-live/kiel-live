package rpc

import "encoding/json"

const internalServiceName = "system"

type SubscribeRequest struct {
	Channel string `json:"c"`
}

type UnsubscribeRequest struct {
	Channel string `json:"c"`
}

type PublishRequest struct {
	Channel string           `json:"c"`
	Data    *json.RawMessage `json:"d"`
}
