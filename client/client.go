package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kiel-live/kiel-live/protocol"
)

// Send pings to peer with this period
const pingPeriod = 30 * time.Second

// WebSocketClient return websocket client connection
type WebSocketClient struct {
	Listen func(msg protocol.ClientMessage)

	configStr string
	sendBuf   chan []byte
	ctx       context.Context
	ctxCancel context.CancelFunc

	mu  sync.RWMutex
	wsc *websocket.Conn
}

// NewWebSocketClient create new websocket connection
func NewWebSocketClient(host string, listen func(msg protocol.ClientMessage)) *WebSocketClient {
	c := WebSocketClient{
		Listen:  listen,
		sendBuf: make(chan []byte, 10),
	}
	c.ctx, c.ctxCancel = context.WithCancel(context.Background())

	u := url.URL{Scheme: "ws", Host: host, Path: "/ws"}
	c.configStr = u.String()

	go c.listen()
	go c.listenWrite()
	go c.ping()
	return &c
}

func (c *WebSocketClient) Connect() *websocket.Conn {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.wsc != nil {
		return c.wsc
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for ; ; <-ticker.C {
		select {
		case <-c.ctx.Done():
			return nil
		default:
			ws, _, err := websocket.DefaultDialer.Dial(c.configStr, nil)
			if err != nil {
				c.log("connect", err, fmt.Sprintf("Cannot connect to websocket: %s", c.configStr))
				continue
			}
			c.log("connect", nil, fmt.Sprintf("connect to websocket to %s", c.configStr))
			c.wsc = ws
			return c.wsc
		}
	}
}

func (c *WebSocketClient) IsConnected() bool {
	// TODO
	return true
}

func (c *WebSocketClient) listen() {
	c.log("listen", nil, fmt.Sprintf("listen for the messages: %s", c.configStr))
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	m := protocol.ClientMessage{}
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			for {
				ws := c.Connect()
				if ws == nil {
					return
				}

				err := ws.ReadJSON(&m)
				if err != nil {
					c.log("listen", err, "Cannot read websocket message")
					c.closeWs()
					break
				}
				if c.Listen != nil {
					c.Listen(m)
				}
			}
		}
	}
}

func (c *WebSocketClient) listenWrite() {
	for data := range c.sendBuf {
		ws := c.Connect()
		if ws == nil {
			err := fmt.Errorf("c.ws is nil")
			c.log("listenWrite", err, "No websocket connection")
			continue
		}

		err := ws.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			c.log("listenWrite", nil, "WebSocket Write Error")
		}
	}
}

// Close will send close message and shutdown websocket connection
func (c *WebSocketClient) Disconnect() {
	c.ctxCancel()
	c.closeWs()
}

func (c *WebSocketClient) Subscribe(channel string) error {
	return c.write(protocol.NewSubscribeMessage(channel))
}

func (c *WebSocketClient) Unsubscribe(channel string) error {
	return c.write(protocol.NewUnsubscribeMessage(channel))
}

func (c *WebSocketClient) Authenticate(token string) error {
	return c.write(protocol.NewAuthenticateMessage(token))
}

func (c *WebSocketClient) Publish(channel string, data string) error {
	return c.write(protocol.NewPublishMessage(channel, data))
}

// Write data to the websocket server
func (c *WebSocketClient) write(msg protocol.ClientMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
	defer cancel()

	for {
		select {
		case c.sendBuf <- data:
			return nil
		case <-ctx.Done():
			return fmt.Errorf("context canceled")
		}
	}
}

// Close will send close message and shutdown websocket connection
func (c *WebSocketClient) closeWs() {
	c.mu.Lock()
	if c.wsc != nil {
		c.wsc.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		c.wsc.Close()
		c.wsc = nil
	}
	c.mu.Unlock()
}

func (c *WebSocketClient) ping() {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			ws := c.Connect()
			if ws == nil {
				continue
			}
			if err := c.wsc.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(pingPeriod/2)); err != nil {
				c.closeWs()
			}
		case <-c.ctx.Done():
			return
		}
	}
}

// Log print log statement
// In real word I would recommend to use zerolog or any other solution
func (c *WebSocketClient) log(f string, err error, msg string) {
	if err != nil {
		fmt.Printf("Error in func: %s, err: %v, msg: %s\n", f, err, msg)
	} else {
		fmt.Printf("Log in func: %s, %s\n", f, msg)
	}
}
