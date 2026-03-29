package api

import "time"

func timeToIsoDateTime(timeStr string, now time.Time) (string, error) {
	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		return "", err
	}

	dateTime := time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), now.Second(), now.Nanosecond(), time.Local)

	// If the parsed time is before the current time, assume it's for the next day
	if dateTime.Before(now) {
		dateTime = dateTime.Add(24 * time.Hour)
	}

	return dateTime.Format(time.RFC3339), nil
}
