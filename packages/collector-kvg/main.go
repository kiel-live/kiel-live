package main

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/kiel-live/kiel-live/packages/client"
	protocol "github.com/kiel-live/kiel-live/packages/pub-sub-proto"
)

func main() {
	client := client.NewWebSocketClient("localhost:4000", func(msg protocol.ClientMessage) {
		fmt.Println("huhu > " + msg.Data())
	})

	client.Connect()
	defer client.Disconnect()

	client.Subscribe(protocol.SubscribedChannelName)

	s := gocron.NewScheduler(time.UTC)
	s.SetMaxConcurrentJobs(1, gocron.RescheduleMode) // prevent parallel execution and skip if last run hasn't finished yet
	s.Every(30).Seconds().Do(func() {
		fetchData()
	})

	s.StartBlocking()
}

func fetchData() {
	fmt.Println("now")
	// stops := getStops()
	// fmt.Println(stops.Stops[0])

	// vehicles := getVehicles()
	// fmt.Println(vehicles.Vehicles[0])

	// stop := getStop(stops.Stops[0].ShortName)
	// fmt.Println(stop)

	// trip := getTrip(stop.Departures[0].TripID)
	// fmt.Println(trip)
}
