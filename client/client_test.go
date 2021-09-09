package client_test

import (
	"fmt"

	"github.com/kiel-live/kiel-live/client"
)

func TestClient() {
	client := client.NewClient("localhost")
	client.Connect()
	defer client.Disconnect()

	client.Subscribe("data.>", func(msg string) {
		fmt.Println(">>>", string(msg))
	})
}
