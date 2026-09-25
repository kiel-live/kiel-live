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
			name:    "empty time",
			input:   "",
			now:     time.Date(2024, 6, 1, 12, 30, 0, 0, time.UTC),
			want:    "",
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
		{
			name:    "previous day time shortly after midnight",
			input:   "23:58",
			now:     time.Date(2024, 6, 1, 22, 5, 0, 0, time.UTC),
			want:    "2024-06-01T23:58:00+02:00",
			wantErr: false,
		},
		{
			name:    "seconds are dropped",
			input:   "14:30",
			now:     time.Date(2024, 6, 1, 12, 30, 45, 123456789, time.UTC),
			want:    "2024-06-01T14:30:00+02:00",
			wantErr: false,
		},
		{
			name:    "time in the current minute stays on the same day",
			input:   "14:30",
			now:     time.Date(2024, 6, 1, 12, 30, 45, 0, time.UTC),
			want:    "2024-06-01T14:30:00+02:00",
			wantErr: false,
		},
		{
			name:    "time one minute in the past stays on the same day",
			input:   "14:29",
			now:     time.Date(2024, 6, 1, 12, 30, 45, 0, time.UTC),
			want:    "2024-06-01T14:29:00+02:00",
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

func TestRelativeToIsoDateTime(t *testing.T) {
	tests := []struct {
		name    string
		anchor  time.Time
		seconds int
		want    string
	}{
		{
			name:    "lands on the full minute",
			anchor:  time.Date(2026, 9, 25, 13, 15, 7, 0, time.UTC),
			seconds: 353,
			want:    "2026-09-25T15:21:00+02:00",
		},
		{
			name:    "negative relative time lands in the past",
			anchor:  time.Date(2026, 9, 25, 13, 15, 7, 0, time.UTC),
			seconds: -247,
			want:    "2026-09-25T15:11:00+02:00",
		},
		{
			name:    "one second header skew is absorbed by rounding",
			anchor:  time.Date(2026, 9, 25, 13, 15, 8, 0, time.UTC),
			seconds: 353,
			want:    "2026-09-25T15:21:00+02:00",
		},
		{
			name:    "crosses midnight into the next day",
			anchor:  time.Date(2026, 9, 25, 21, 58, 30, 0, time.UTC),
			seconds: 210,
			want:    "2026-09-26T00:02:00+02:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := relativeToIsoDateTime(tt.anchor, tt.seconds); got != tt.want {
				t.Errorf("relativeToIsoDateTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
