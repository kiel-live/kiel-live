package client_test

import (
	"fmt"
	"testing"

	"github.com/kiel-live/kiel-live/pkg/client"
)

func TestClient(t *testing.T) {
	c := client.NewNatsClient("localhost")
	err := c.Connect()
	if err != nil {
		t.Error(err)
	}

	defer func() {
		err := c.Disconnect()
		if err != nil {
			t.Error(err)
		}
	}()

	err = c.Subscribe("data.>", func(msg *client.TopicMessage) {
		fmt.Println(">>>", msg.Data)
	})
	if err != nil {
		t.Error(err)
	}
}
