package api

import (
	"io"
	"net/http"
	"net/url"
	"strings"
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

func post(url string, data url.Values) ([]byte, error) {
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
