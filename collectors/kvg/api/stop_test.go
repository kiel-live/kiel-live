package api

import (
	"testing"
	"time"
)

func TestDepartureParse(t *testing.T) {
	// Sample from the KVG API: HTTP Date header 13:15:07 GMT
	serverTime := time.Date(2026, 9, 25, 13, 15, 7, 0, time.UTC)
	rel := func(s int) *int { return &s }

	tests := []struct {
		name        string
		departure   departure
		wantActual  string
		wantPlanned string
	}{
		{
			name:        "predicted uses relative time",
			departure:   departure{Status: predicted, PlannedTime: "15:16", ActualTime: "15:21", ActualRelativeTime: rel(353)},
			wantActual:  "2026-09-25T15:21:00+02:00",
			wantPlanned: "2026-09-25T15:16:00+02:00",
		},
		{
			name:        "planned without actual time keeps actual empty",
			departure:   departure{Status: planned, PlannedTime: "15:22", ActualTime: "", ActualRelativeTime: rel(413)},
			wantActual:  "",
			wantPlanned: "2026-09-25T15:22:00+02:00",
		},
		{
			name:        "departed without actual time keeps actual empty",
			departure:   departure{Status: departed, PlannedTime: "15:12", ActualTime: "", ActualRelativeTime: rel(-247)},
			wantActual:  "",
			wantPlanned: "2026-09-25T15:12:00+02:00",
		},
		{
			name:        "missing relative time falls back to clock string",
			departure:   departure{Status: predicted, PlannedTime: "15:16", ActualTime: "15:21", ActualRelativeTime: nil},
			wantActual:  "2026-09-25T15:21:00+02:00",
			wantPlanned: "2026-09-25T15:16:00+02:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.departure.parse(serverTime)
			if got.Actual != tt.wantActual {
				t.Errorf("Actual = %v, want %v", got.Actual, tt.wantActual)
			}
			if got.Planned != tt.wantPlanned {
				t.Errorf("Planned = %v, want %v", got.Planned, tt.wantPlanned)
			}
		})
	}
}
