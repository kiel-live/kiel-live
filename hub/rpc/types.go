package rpc

import "encoding/json"

const defaultServiceName = "main"
const internalServiceName = "internal"

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
