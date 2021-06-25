package api

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	baseURL     = "https://www.kvg-kiel.de"
	stopURL     = baseURL + "/internetservice/services/passageInfo/stopPassages/stop"
	tripURL     = baseURL + "/internetservice/services/tripInfo/tripPassages"
	vehiclesURL = baseURL + "/internetservice/geoserviceDispatcher/services/vehicleinfo/vehicles"
	stopsURL    = baseURL + "/internetservice/geoserviceDispatcher/services/stopinfo/stops"
)

func post(url string, data url.Values) ([]byte, error) {
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		println(err)
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
