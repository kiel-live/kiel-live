package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/kiel-live/kiel-live/client"
	"github.com/kiel-live/kiel-live/protocol"
)

const ProviderID = "kvg"

var collectors map[string]*collector

func main() {
	// TODO load client address from env
	client := client.NewWebSocketClient("localhost:4000", func(msg protocol.ClientMessage) {
		if msg.Type() == protocol.AuthFailedMessage {
			log.Fatalln("Authentication failed")
			return
		}

		if msg.Type() == protocol.ChannelMessage && msg.Channel() == protocol.ChannelNameSubscribedChannels {
			fmt.Println("subscribed-channels > " + msg.Data())
			// TODO check which channel is new and needs a collector
			// TODO remove collectors of channels not being subscribed anymore
			return
		}

		log.Println(">>>", msg)
	})

	collectors = make(map[string]*collector)

	// auto load following collectors
	collector, _ := newCollector(client, protocol.ChannelNameVehicles)
	collectors[protocol.ChannelNameVehicles] = collector
	collector, _ = newCollector(client, protocol.ChannelNameStops)
	collectors[protocol.ChannelNameStops] = collector

	client.Connect()
	defer client.Disconnect()

	// TODO get token from env
	client.Authenticate("123")

	client.Subscribe(protocol.ChannelNameSubscribedChannels)

	s := gocron.NewScheduler(time.UTC)
	s.SetMaxConcurrentJobs(1, gocron.RescheduleMode)
	s.Every(30).Seconds().Do(func() {
		if !client.IsConnected() {
			return
		}

		for _, c := range collectors {
			// TODO maybe run in go routine
			fmt.Println("Collector for", c.channel, "running ...")
			c.run()
		}
	})

	s.StartBlocking()
}
