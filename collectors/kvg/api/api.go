package api

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	baseURL      = "https://kvg-internetservice-proxy.p.networkteam.com"
	stopURL      = baseURL + "/internetservice/services/passageInfo/stopPassages/stop"
	platformURL  = baseURL + "/internetservice/services/passageInfo/stopPassages/stopPoint"
	tripURL      = baseURL + "/internetservice/services/tripInfo/tripPassages"
	tripPathURL  = baseURL + "/internetservice/geoserviceDispatcher/services/pathinfo/trip"
	vehiclesURL  = baseURL + "/internetservice/geoserviceDispatcher/services/vehicleinfo/vehicles"
	stopsURL     = baseURL + "/internetservice/geoserviceDispatcher/services/stopinfo/stops"
	platformsURL = baseURL + "/internetservice/geoserviceDispatcher/services/stopinfo/stopPoints"
)

const IDPrefix = "kvg-"

// The KVG API computes relative times against its own clock, so the Date header is the right anchor for converting them to absolute timestamps.
func postWithServerTime(url string, data url.Values) ([]byte, time.Time, error) {
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, time.Time{}, err
	}
	defer resp.Body.Close()

	serverTime := time.Now()
	if t, err := http.ParseTime(resp.Header.Get("Date")); err == nil {
		serverTime = t
	}

	body, err := io.ReadAll(resp.Body)
	return body, serverTime, err
}

func post(url string, data url.Values) ([]byte, error) {
	body, _, err := postWithServerTime(url, data)
	return body, err
}
