package api

import (
	"testing"
)

func TestFormatMotisID(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"de:01001:1", "motis-de-01001-1"},
		{"DE:01001:1", "motis-de-01001-1"}, // uppercased → lowercased
		{"de_01001_1", "motis-de-01001-1"}, // underscores → dashes
		{"de:01001_abc", "motis-de-01001-abc"},
		{"simple", "motis-simple"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := FormatMotisID(tt.input)
			if got != tt.want {
				t.Errorf("FormatMotisID(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestTopicToID(t *testing.T) {
	tests := []struct {
		topic string
		want  string
	}{
		{"data.map.stop.motis-de-01001-1", "motis-de-01001-1"},
		{"data.map.stop.motis-abc-xyz", "motis-abc-xyz"},
		// wildcard must return empty
		{"data.map.stop.>", ""},
		// wrong provider prefix
		{"data.map.stop.kvg-123", ""},
		// wrong topic kind
		{"data.map.vehicle.motis-123", ""},
		// unrelated
		{"ctrl.subscriptions", ""},
	}

	for _, tt := range tests {
		t.Run(tt.topic, func(t *testing.T) {
			got := TopicToID(tt.topic)
			if got != tt.want {
				t.Errorf("TopicToID(%q) = %q, want %q", tt.topic, got, tt.want)
			}
		})
	}
}
