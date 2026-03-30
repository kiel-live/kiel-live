package api

import (
	"testing"
	"time"
)

func TestTimeToIsoDateTime(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		now     time.Time
		want    string
		wantErr bool
	}{
		{
			name:    "valid time",
			input:   "14:30",
			now:     time.Date(2024, 6, 1, 12, 30, 0, 0, time.UTC),
			want:    "2024-06-01T14:30:00+02:00",
			wantErr: false,
		},
		{
			name:    "invalid time format",
			input:   "2:30 PM",
			now:     time.Date(2024, 6, 1, 12, 30, 0, 0, time.UTC),
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid time value",
			input:   "25:00",
			now:     time.Date(2024, 6, 1, 12, 30, 0, 0, time.UTC),
			want:    "",
			wantErr: true,
		},
		{
			name:    "next day time",
			input:   "01:44",
			now:     time.Date(2024, 6, 1, 23, 0, 0, 0, time.UTC),
			want:    "2024-06-02T01:44:00+02:00",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := timeToIsoDateTime(tt.input, tt.now)
			if (err != nil) != tt.wantErr {
				t.Errorf("timeToIsoDateTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("timeToIsoDateTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
