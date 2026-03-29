package api

import (
	"net/url"
	"strings"
	"time"

	"github.com/kiel-live/kiel-live/pkg/models"
)

// Kiel region bounding box: SW corner (min) and NE corner (max)
const (
	kielerBboxMin = "54.20,9.90"
	kielerBboxMax = "54.50,10.30"
)

// nonBusModes are the transit modes we collect — excludes BUS and COACH
// since KVG provides first-class realtime data for those.
var nonBusModes = []string{
	"TRAM", "SUBWAY", "FERRY", "SUBURBAN",
	"REGIONAL_RAIL", "REGIONAL_FAST_RAIL",
	"LONG_DISTANCE", "NIGHT_RAIL", "HIGHSPEED_RAIL",
}

type StopWithProviderID struct {
	*models.Stop
	ProviderIDs []string `json:"-"`
}

// place mirrors the MOTIS Place schema (subset of fields we use)
type place struct {
	Name               string     `json:"name"`
	StopID             string     `json:"stopId"`
	Lat                float64    `json:"lat"`
	Lon                float64    `json:"lon"`
	Departure          *time.Time `json:"departure"`
	ScheduledDeparture *time.Time `json:"scheduledDeparture"`
	Track              string     `json:"track"`
	ScheduledTrack     string     `json:"scheduledTrack"`
	Modes              []string   `json:"modes"`
	Cancelled          bool       `json:"cancelled"`
}

func modeToStopType(modes []string) models.StopType {
	for _, m := range modes {
		switch m {
		case "FERRY":
			return models.StopTypeFerryStop
		case "SUBWAY":
			return models.StopTypeSubwayStop
		case "SUBURBAN", "RAIL", "REGIONAL_RAIL", "REGIONAL_FAST_RAIL",
			"LONG_DISTANCE", "NIGHT_RAIL", "HIGHSPEED_RAIL":
			return models.StopTypeTrainStop
		}
	}
	return models.StopTypeBusStop
}

func modeToVehicleType(mode string) models.VehicleType {
	switch mode {
	case "TRAM":
		return models.VehicleTypeTram
	case "SUBWAY":
		return models.VehicleTypeSubway
	case "FERRY":
		return models.VehicleTypeFerry
	case "SUBURBAN", "RAIL", "REGIONAL_RAIL", "REGIONAL_FAST_RAIL",
		"LONG_DISTANCE", "NIGHT_RAIL", "HIGHSPEED_RAIL":
		return models.VehicleTypeTrain
	default:
		return models.VehicleTypeBus
	}
}

// GetStops fetches all non-bus transit stops in the Kiel region.
// Only stops with at least one confirmed non-bus mode are returned,
// to avoid duplicating bus stops already covered by the KVG collector.
func GetStops() (map[string]*StopWithProviderID, error) {
	var places []place
	params := url.Values{}
	params.Set("min", kielerBboxMin)
	params.Set("max", kielerBboxMax)
	params.Set("language", "de")

	if err := getJSON("/api/v1/map/stops", params, &places); err != nil {
		return nil, err
	}

	result := make(map[string]*StopWithProviderID)
	for _, p := range places {
		if p.StopID == "" {
			continue
		}

		if skipStop(p) {
			continue
		}

		id := FormatMotisID(p.StopID)
		result[id] = &StopWithProviderID{
			Stop: &models.Stop{
				ID:       id,
				Provider: "motis",
				Name:     p.Name,
				Type:     modeToStopType(p.Modes),
				Location: &models.Location{
					Latitude:  int(p.Lat * 3600000),
					Longitude: int(p.Lon * 3600000),
				},
			},
			ProviderIDs: []string{p.StopID},
		}
	}

	return mergeNearbyStops(result), nil
}

// mergeNearbyStops combines stops that share the same name and are within
// 200 m of each other into a single entry. The first-seen stop's ID and
// location are kept; duplicate provider IDs are appended so departures can
// be fetched from all underlying MOTIS feeds.
func mergeNearbyStops(stops map[string]*StopWithProviderID) map[string]*StopWithProviderID {
	const mergeRadiusMeters = 200.0

	// canonical[name] = the representative stop for that name
	canonical := make(map[string]*StopWithProviderID)
	result := make(map[string]*StopWithProviderID)

	for id, s := range stops {
		rep, exists := canonical[s.Name]
		if exists && rep.Location.DistanceToMeters(s.Location) <= mergeRadiusMeters {
			// Merge: append provider IDs to the representative stop.
			rep.ProviderIDs = append(rep.ProviderIDs, s.ProviderIDs...)
		} else {
			// New canonical entry (either new name or too far from existing one).
			result[id] = s
			canonical[s.Name] = s
		}
	}

	return result
}

// stopTime mirrors the MOTIS StopTime schema (subset of fields we use)
type stopTime struct {
	Place          place  `json:"place"`
	Mode           string `json:"mode"`
	RealTime       bool   `json:"realTime"`
	Headsign       string `json:"headsign"`
	TripID         string `json:"tripId"`
	RouteShortName string `json:"routeShortName"`
	DisplayName    string `json:"displayName"`
	Cancelled      bool   `json:"cancelled"`
	TripCancelled  bool   `json:"tripCancelled"`
}

type stopTimesResponse struct {
	StopTimes []stopTime `json:"stopTimes"`
}

// GetStopDepartures fetches the next departures for one or more MOTIS stop IDs
// (merged stops may map to multiple provider IDs). Results are deduplicated by
// trip ID so the same service is not shown twice. Only non-bus modes are requested.
func GetStopDepartures(providerIDs []string) ([]*models.StopArrival, error) {
	seen := make(map[string]bool)
	now := time.Now()
	var arrivals []*models.StopArrival

	for _, stopID := range providerIDs {
		params := url.Values{}
		params.Set("stopId", stopID)
		params.Set("n", "10")
		params.Set("mode", strings.Join(nonBusModes, ","))
		params.Set("language", "de")

		var resp stopTimesResponse
		if err := getJSON("/api/v5/stoptimes", params, &resp); err != nil {
			return nil, err
		}

		for _, st := range resp.StopTimes {
			if seen[st.TripID] || st.TripCancelled {
				continue
			}
			seen[st.TripID] = true

			// Prefer realtime departure time, fall back to scheduled
			var depTime *time.Time
			if st.RealTime && st.Place.Departure != nil {
				depTime = st.Place.Departure
			} else if st.Place.ScheduledDeparture != nil {
				depTime = st.Place.ScheduledDeparture
			}

			eta := 0
			planned := ""
			if depTime != nil {
				eta = int(depTime.Sub(now).Seconds())
				planned = depTime.Format("15:04")
			}

			state := models.Planned
			if st.RealTime {
				state = models.Predicted
			}
			if st.Cancelled {
				state = models.Departed
			}

			platform := st.Place.Track
			if platform == "" {
				platform = st.Place.ScheduledTrack
			}

			routeName := st.DisplayName
			if routeName == "" {
				routeName = st.RouteShortName
			}

			arrivals = append(arrivals, &models.StopArrival{
				Name:      st.Headsign,
				Type:      modeToVehicleType(st.Mode),
				TripID:    FormatMotisID(st.TripID),
				VehicleID: "", // map to some vehicle / trip?
				RouteName: routeName,
				Direction: st.Headsign,
				State:     state,
				Eta:       eta,
				Planned:   planned,
				Platform:  platform,
			})
		}
	}

	return arrivals, nil
}
