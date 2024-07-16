package main

import (
	"time"

	"github.com/artonge/go-gtfs"
	"github.com/kiel-live/kiel-live/protocol"
)

func findInObjArr[T any, K comparable](arr []T, keyFunc func(T) K, value K) (int, bool) {
	for i, v := range arr {
		if keyFunc(v) == value {
			return i, true
		}
	}
	return -1, false
}

func weekdayIsActiveInCalendar(calendar gtfs.Calendar) bool {
	weekday := time.Now().Weekday()
	switch weekday {
	case time.Monday:
		return calendar.Monday == 1
	case time.Tuesday:
		return calendar.Tuesday == 1
	case time.Wednesday:
		return calendar.Wednesday == 1
	case time.Thursday:
		return calendar.Thursday == 1
	case time.Friday:
		return calendar.Friday == 1
	case time.Saturday:
		return calendar.Saturday == 1
	case time.Sunday:
		return calendar.Sunday == 1
	}
	return false
}

// 0 - Tram, Streetcar, Light rail. Any light rail or street level system within a metropolitan area.
// 1 - Subway, Metro. Any underground rail system within a metropolitan area.
// 2 - Rail. Used for intercity or long-distance travel.
// 3 - Bus. Used for short- and long-distance bus routes.
// 4 - Ferry. Used for short- and long-distance boat service.
// 5 - Cable tram. Used for street-level rail cars where the cable runs beneath the vehicle (e.g., cable car in San Francisco).
// 6 - Aerial lift, suspended cable car (e.g., gondola lift, aerial tramway). Cable transport where cabins, cars, gondolas or open chairs are suspended by means of one or more cables.
// 7 - Funicular. Any rail system designed for steep inclines.
// 11 - Trolleybus. Electric buses that draw power from overhead wires using poles.
// 12 - Monorail. Railway in which the track consists of a single rail or a beam.
func gtfsRouteTypeToProtocolStopType(stopType int) protocol.StopType {
	switch stopType {
	case 0:
		return "tram"
	case 1:
		return "subway"
	case 2:
		return "rail"
	case 3:
		return "bus"
	case 4:
		return "ferry"
	case 5:
		return "cable_tram"
	case 6:
		return "aerial_lift"
	case 7:
		return "funicular"
	case 11:
		return "trolleybus"
	case 12:
		return "monorail"
	}

	return "unknown"
}
