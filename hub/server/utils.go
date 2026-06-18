package server

import (
	"encoding/json"

	"github.com/kiel-live/kiel-live/pkg/models"
)

func toRaw(v any) *json.RawMessage {
	data, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	raw := json.RawMessage(data)
	return &raw
}

func stopsToAny(stops []*models.Stop) []any {
	out := make([]any, len(stops))
	for i, s := range stops {
		out[i] = s
	}
	return out
}

func vehiclesToAny(vehicles []*models.Vehicle) []any {
	out := make([]any, len(vehicles))
	for i, v := range vehicles {
		out[i] = v
	}
	return out
}
