package rpc

const serviceName = "service"
const internalServiceName = "internal"

// type Message struct {
// 	MessageType string `json:"t"`
// 	ID          string `json:"i"`
// 	Method      string `json:"m"`
// 	Data        []byte `json:"d"`
// }

type ChannelMessage struct {
	Channel string `json:"c"`
	Data    []byte `json:"d"`
}

type SubscribeRequest struct {
	Channel string `json:"c"`
}

type UnsubscribeRequest struct {
	Channel string `json:"c"`
}
