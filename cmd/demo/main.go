package main

import (
	"fmt"

	"github.com/kiel-live/kiel-live/client"
	log "github.com/sirupsen/logrus"
)

func main() {

	c := client.NewClient("localhost")
	c.Connect()
	defer c.Disconnect()

	log.Infoln("⚡ Demo consumer connected")

	err := c.Subscribe("data.map.>", func(msg *client.SubjectMessage) {
		fmt.Println("new data", msg.Data)
	}, c.WithAck(), c.WithCache())
	if err != nil {
		log.Fatalln(err)
	}

	// don't kill main process
	select {}
}
