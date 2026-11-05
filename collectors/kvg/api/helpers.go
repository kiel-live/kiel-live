package api

import "time"

// The KVG API reports clock strings in German local time, independent of the timezone the collector runs in.
var kvgLocation = func() *time.Location {
	location, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		panic(err)
	}
	return location
}()

func timeToIsoDateTime(timeStr string, ref time.Time) (string, error) {
	if timeStr == "" {
		return "", nil
	}

	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		return "", err
	}

	ref = ref.In(kvgLocation)
	dateTime := time.Date(ref.Year(), ref.Month(), ref.Day(), t.Hour(), t.Minute(), 0, 0, kvgLocation)

	// The nearest day is chosen instead of always rolling forward so that a departure a minute in the past is not pushed a day ahead.
	if diff := dateTime.Sub(ref); diff > 12*time.Hour {
		dateTime = dateTime.AddDate(0, 0, -1)
	} else if diff < -12*time.Hour {
		dateTime = dateTime.AddDate(0, 0, 1)
	}

	return dateTime.Format(time.RFC3339), nil
}

func relativeToIsoDateTime(anchor time.Time, seconds int) string {
	// Rounding to the minute absorbs a possible one-second offset between the Date header and the server clock the relative time was computed against.
	return anchor.Add(time.Duration(seconds) * time.Second).Round(time.Minute).In(kvgLocation).Format(time.RFC3339)
}
