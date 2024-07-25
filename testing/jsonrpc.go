package testing

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	websocketjsonrpc2 "github.com/sourcegraph/jsonrpc2/websocket"

	"github.com/kiel-live/kiel-live/jsonrpc/rpc"
)

type jsonrpc struct{}

func (g *jsonrpc) Name() string {
	return "jsonrpc"
}

func (g *jsonrpc) SendData(testSet *TestSet) error {
	url := "ws://localhost:4568/ws"
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	defer c.Close()

	ctx := context.Background()
	rpc := rpc.NewClientPeer(ctx, websocketjsonrpc2.NewObjectStream(c))
	defer rpc.Close()

	testSet.StartTime = time.Now()
	fmt.Printf("%s: %d\n", "sending", time.Now().UnixMicro())
	return rpc.Publish(ctx, "Hello", testSet.ID)
}

func (g *jsonrpc) WaitForMessage(_ []*TestSet, connectingWG *sync.WaitGroup, done func(s string)) error {
	url := "ws://localhost:4568/ws"
	c, resp, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	defer c.Close()

	ctx := context.Background()
	rpc := rpc.NewClientPeer(ctx, websocketjsonrpc2.NewObjectStream(c))
	defer rpc.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)

	err = rpc.Subscribe(ctx, "Hello", func(rm *json.RawMessage) {
		done("last")

		reply := ""
		err := json.Unmarshal(*rm, &reply)
		if err != nil {
			return
		}

		if reply == "last" {
			wg.Done()
		}
	})
	if err != nil {
		return err
	}

	fmt.Printf("%s: %d\n", "connected", time.Now().UnixMicro())

	connectingWG.Done()

	wg.Wait()

	fmt.Printf("%s: %d\n", "client done", time.Now().UnixMicro())

	return rpc.Unsubscribe(ctx, "Hello")
}
