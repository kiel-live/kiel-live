package protocol

import (
	"encoding/json"
	"strings"
)

type Method string

const (
	MethodSubscribe   Method = "subscribe"
	MethodUnsubscribe Method = "unsubscribe"
	MethodViewport    Method = "viewport"
	MethodSnapshot    Method = "snapshot"
	MethodUpdate      Method = "update"
	MethodDelete      Method = "delete"
	MethodRequest     Method = "request"
	MethodResult      Method = "result"
	MethodError       Method = "error"
)

const TopicSubscriptions = "subscriptions"
const TopicSearch = "search"

type Envelope struct {
	RID    string           `json:"rid,omitempty"`
	Method Method           `json:"method"`
	Topic  string           `json:"topic,omitempty"`
	Data   *json.RawMessage `json:"data,omitempty"`
}

type ViewportData struct {
	North float64 `json:"north"`
	East  float64 `json:"east"`
	South float64 `json:"south"`
	West  float64 `json:"west"`
}

type SubscriptionEvent struct {
	Topic      string `json:"topic"`
	Subscribed bool   `json:"subscribed"`
}

type SearchRequest struct {
	Query  string        `json:"query"`
	Limit  int           `json:"limit,omitempty"` // defaults to 20 when zero
	Bounds *ViewportData `json:"bounds,omitempty"`
}

// ParseTopic parses "map.<kind>.<id>" or "<kind>.<id>".
// Dotted / colon GTFS ids in the id segment are preserved intact.
func ParseTopic(topic string) (isMap bool, kind, id string) {
	if strings.HasPrefix(topic, "map.") {
		parts := strings.SplitN(topic[4:], ".", 2)
		if len(parts) == 2 {
			return true, parts[0], parts[1]
		}
	}
	parts := strings.SplitN(topic, ".", 2)
	if len(parts) == 2 {
		return false, parts[0], parts[1]
	}
	return false, topic, ""
}

func FormatDetailTopic(kind, id string) string {
	return kind + "." + id
}

func FormatMapTopic(kind, id string) string {
	return "map." + kind + "." + id
}

func FormatMapKindTopic(kind string) string {
	return "map." + kind
}

// HubTopicToNats converts a hub detail topic ("stop.<id>") to a NATS topic ("data.map.stop.<id>").
func HubTopicToNats(topic string) string {
	return "data.map." + topic
}
