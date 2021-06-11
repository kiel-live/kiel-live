package protocol

import "fmt"

type ClientMessage map[string]interface{}

func (c ClientMessage) ResultID() string {
	t := c.Type()
	if t == SubscribeOKMessage || t == SubscribeErrorMessage {
		t = SubscribeMessage
	}
	if t == UnsubscribeOKMessage {
		t = UnsubscribeMessage
	}
	return fmt.Sprintf("%s_%s", t, c["__channel"])
}

func (c ClientMessage) Type() string {
	s, ok := c["__type"].(string)
	if !ok {
		return ""
	}
	return s
}

func (c ClientMessage) Channel() string {
	s, ok := c["__channel"].(string)
	if !ok {
		return ""
	}
	return s
}

func (c ClientMessage) Data() string {
	s, ok := c["data"].(string)
	if !ok {
		return ""
	}
	return s
}

func NewMessage(t string) ClientMessage {
	return ClientMessage{
		"__type": t,
	}
}

func NewErrorMessage(t string, err error) ClientMessage {
	return ClientMessage{
		"__type": t,
		"reason": err.Error(),
	}
}

func NewChannelMessage(t string, channel string) ClientMessage {
	return ClientMessage{
		"__type":    t,
		"__channel": channel,
	}
}

func NewBroadcastMessage(channel, data string) ClientMessage {
	return ClientMessage{
		"__type":    ChannelMessage,
		"__channel": channel,
		"data":      data,
	}
}

func NewChannelErrorMessage(t string, channel string, err error) ClientMessage {
	return ClientMessage{
		"__type":    t,
		"__channel": channel,
		"reason":    err.Error(),
	}
}

func NewSubscribeMessage(channel string) ClientMessage {
	return ClientMessage{
		"__type":    SubscribeMessage,
		"__channel": channel,
	}
}

func NewUnsubscribeMessage(channel string) ClientMessage {
	return ClientMessage{
		"__type":    UnsubscribeMessage,
		"__channel": channel,
	}
}

func NewAuthenticateMessage(data string) ClientMessage {
	return ClientMessage{
		"__type": AuthMessage,
		"data":   data,
	}
}

func NewPublishMessage(channel string, data string) ClientMessage {
	return ClientMessage{
		"__type":    PublishMessage,
		"__channel": channel,
		"data":      data,
	}
}
