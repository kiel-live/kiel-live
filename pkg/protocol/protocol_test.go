package protocol_test

import (
	"testing"

	"github.com/kiel-live/kiel-live/pkg/protocol"
	"github.com/stretchr/testify/assert"
)

func TestParseTopic(t *testing.T) {
	tests := []struct {
		topic   string
		wantMap bool
		kind    string
		id      string
	}{
		{"stop.kvg-123", false, "stop", "kvg-123"},
		{"vehicle.ais-211865680", false, "vehicle", "ais-211865680"},
		{"trip.kvg-456", false, "trip", "kvg-456"},
		{"map.stop.kvg-123", true, "stop", "kvg-123"},
		{"map.vehicle.ais-211865680", true, "vehicle", "ais-211865680"},
		// GTFS ids with dots, colons, slashes must survive intact
		{"stop.gtfs-foo.bar:baz/qux", false, "stop", "gtfs-foo.bar:baz/qux"},
		{"map.stop.gtfs-agency:route.stop", true, "stop", "gtfs-agency:route.stop"},
	}

	for _, tt := range tests {
		t.Run(tt.topic, func(t *testing.T) {
			isMap, kind, id := protocol.ParseTopic(tt.topic)
			assert.Equal(t, tt.wantMap, isMap)
			assert.Equal(t, tt.kind, kind)
			assert.Equal(t, tt.id, id)
		})
	}
}
