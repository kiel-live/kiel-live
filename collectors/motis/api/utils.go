package api

import (
	"fmt"
	"strings"

	"github.com/kiel-live/kiel-live/pkg/models"
)

var nonBusModeSet = map[string]bool{
	"TRAM": true, "SUBWAY": true, "FERRY": true, "SUBURBAN": true,
	"RAIL": true, "REGIONAL_RAIL": true, "REGIONAL_FAST_RAIL": true,
	"LONG_DISTANCE": true, "NIGHT_RAIL": true, "HIGHSPEED_RAIL": true,
}

func TopicToID(topic string) string {
	prefix := fmt.Sprintf(models.TopicStop, IDPrefix)
	wildcard := fmt.Sprintf(models.TopicStop, ">")
	if strings.HasPrefix(topic, prefix) && topic != wildcard {
		return IDPrefix + strings.TrimPrefix(topic, prefix)
	}

	return ""
}

func FormatMotisID(motisStopID string) string {
	motisStopID = strings.ReplaceAll(motisStopID, ":", "-")
	motisStopID = strings.ReplaceAll(motisStopID, "_", "-")
	motisStopID = strings.ToLower(motisStopID)
	return IDPrefix + motisStopID
}

func skipStop(p place) bool {
	// Skip stops where modes are known but contain no non-bus mode.
	// If modes is empty we can't determine the type, so we skip too.
	if len(p.Modes) == 0 {
		return true
	}

	hasNonBus := false
	for _, m := range p.Modes {
		if nonBusModeSet[m] {
			hasNonBus = true
			break
		}
	}
	if !hasNonBus {
		return true
	}

	return false
}
