package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	funk "github.com/thoas/go-funk"
)

var baseUrl = "https://www.kvg-kiel.de"
var stopUrl = baseUrl + "/internetservice/services/passageInfo/stopPassages/stop"
var tripUrl = baseUrl + "/internetservice/services/tripInfo/tripPassages"
var vehiclesUrl = baseUrl + "/internetservice/geoserviceDispatcher/services/vehicleinfo/vehicles"
var stopsUrl = baseUrl + "/internetservice/geoserviceDispatcher/services/stopinfo/stops"

func main() {
	stops := getStops()
	fmt.Println(stops.Stops[0])

	vehicles := getVehicles()
	fmt.Println(vehicles.Vehicles[0])

	stop := getStop(stops.Stops[0].ShortName)
	fmt.Println(stop)

	trip := getTrip(stop.Departures[0].TripId)
	fmt.Println(trip)
}

func post(url string, data url.Values) ([]byte, error) {
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		println(err)
		return nil, err
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)
	body, err := io.ReadAll(resp.Body)
	return body, err
}

type stop struct {
	Id        string `json:"id"`
	ShortName string `json:"shortName"`
	Name      string `json:"name"`
	Latitude  int    `json:"latitude"`
	Longitude int    `json:"longitude"`
}

type stops struct {
	Stops []stop `json:"stops"`
}

func getStops() stops {
	data := url.Values{}
	data.Set("top", "324000000")
	data.Set("bottom", "-324000000")
	data.Set("left", "-648000000")
	data.Set("right", "648000000")

	body, _ := post(stopsUrl, data)
	fmt.Printf("Body: %s\n", body)
	var stops stops
	if err := json.Unmarshal(body, &stops); err != nil {
		log.Fatalf("Parse response failed, reason: %v \n", err)
	}
	return stops
}

type vehicle struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Heading   int    `json:"heading"`
	Latitude  int    `json:"latitude"`
	Longitude int    `json:"longitude"`
	TripId    string `json:"tripId"`
}

type vehicles struct {
	Vehicles []vehicle `json:"vehicles"`
}

func getVehicles() vehicles {
	body, _ := post(vehiclesUrl, nil)
	var vehicles vehicles
	if err := json.Unmarshal(body, &vehicles); err != nil {
		log.Fatalf("Parse response failed, reason: %v \n", err)
	}

	vehicles.Vehicles = funk.Filter(vehicles.Vehicles, func(vehicle vehicle) bool {
		return vehicle.Latitude != 0
	}).([]vehicle)

	return vehicles
}

type departure struct {
	TripId             string `json:"tripId"`
	Status             string `json:"status"`
	Stop               string `json:"plannedTime"`
	ActualTime         string `json:"actualTime"`
	ActualRelativeTime int    `json:"actualRelativeTime"`
}

type stopDepartures struct {
	Departures []departure `json:"actual"`
}

func getStop(stopShortName string) stopDepartures {
	data := url.Values{}
	data.Set("stop", stopShortName)

	resp, _ := post(stopUrl, data)
	var stop stopDepartures
	if err := json.Unmarshal(resp, &stop); err != nil {
		log.Fatalf("Parse response failed, reason: %v \n", err)
	}
	return stop
}

type tripStops struct {
	Stop       stop   `json:"stop"`
	Status     string `json:"status"`
	ActualTime string `json:"actualTime"`
}

type trip struct {
	Stops         []tripStops `json:"actual"`
	OldStops      []tripStops `json:"old"`
	DirectionText string      `json:"directionText"`
	RouteName     string      `json:"routeName"`
}

func getTrip(tripId string) trip {
	data := url.Values{}
	data.Set("tripId", tripId)

	resp, _ := post(tripUrl, data)
	var trip trip
	if err := json.Unmarshal(resp, &trip); err != nil {
		log.Fatalf("Parse response failed, reason: %v \n", err)
	}
	return trip
}
