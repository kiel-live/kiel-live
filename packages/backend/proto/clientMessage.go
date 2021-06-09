package proto

import "fmt"

type ClientMessage map[string]interface{}

func (c ClientMessage) ResultId() string {
	t := c.Type()
	if t == SubscribeOKMessage || t == SubscribeErrorMessage {
		t = SubscribeMessage
	}
	if t == UnsubscribeOKMessage {
		t = UnsubscribeMessage
	}
	return fmt.Sprintf("%s_%s", t, c["channel"])
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

func NewChannelMessage(t, channel string) ClientMessage {
	return ClientMessage{
		"__type":  t,
		"channel": channel,
	}
}

func NewBroadcastMessage(channel, body string) ClientMessage {
	return ClientMessage{
		"__type":  MessageMessage,
		"channel": channel,
		"body":    body,
	}
}

func NewChannelErrorMessage(t, channel string, err error) ClientMessage {
	return ClientMessage{
		"__type":  t,
		"channel": channel,
		"reason":  err.Error(),
	}
}
