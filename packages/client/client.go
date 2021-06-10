package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	proto "github.com/kiel-live/kiel-live/packages/pub-sub-proto"
)

// Send pings to peer with this period
const pingPeriod = 30 * time.Second

// WebSocketClient return websocket client connection
type WebSocketClient struct {
	Listen func(msg proto.ClientMessage)

	configStr string
	sendBuf   chan []byte
	ctx       context.Context
	ctxCancel context.CancelFunc

	mu     sync.RWMutex
	wsconn *websocket.Conn
}

// NewWebSocketClient create new websocket connection
func NewWebSocketClient(host string, listen func(msg proto.ClientMessage)) *WebSocketClient {
	conn := WebSocketClient{
		Listen:  listen,
		sendBuf: make(chan []byte, 10),
	}
	conn.ctx, conn.ctxCancel = context.WithCancel(context.Background())

	u := url.URL{Scheme: "ws", Host: host, Path: "/ws"}
	conn.configStr = u.String()

	go conn.listen()
	go conn.listenWrite()
	go conn.ping()
	return &conn
}

func (conn *WebSocketClient) Connect() *websocket.Conn {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	if conn.wsconn != nil {
		return conn.wsconn
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for ; ; <-ticker.C {
		select {
		case <-conn.ctx.Done():
			return nil
		default:
			ws, _, err := websocket.DefaultDialer.Dial(conn.configStr, nil)
			if err != nil {
				conn.log("connect", err, fmt.Sprintf("Cannot connect to websocket: %s", conn.configStr))
				continue
			}
			conn.log("connect", nil, fmt.Sprintf("connected to websocket to %s", conn.configStr))
			conn.wsconn = ws
			return conn.wsconn
		}
	}
}

func (conn *WebSocketClient) listen() {
	conn.log("listen", nil, fmt.Sprintf("listen for the messages: %s", conn.configStr))
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	m := proto.ClientMessage{}
	for {
		select {
		case <-conn.ctx.Done():
			return
		case <-ticker.C:
			for {
				ws := conn.Connect()
				if ws == nil {
					return
				}

				err := ws.ReadJSON(&m)
				if err != nil {
					conn.log("listen", err, "Cannot read websocket message")
					conn.closeWs()
					break
				}
				if conn.Listen != nil {
					conn.Listen(m)
				}
			}
		}
	}
}

func (conn *WebSocketClient) listenWrite() {
	for data := range conn.sendBuf {
		ws := conn.Connect()
		if ws == nil {
			err := fmt.Errorf("conn.ws is nil")
			conn.log("listenWrite", err, "No websocket connection")
			continue
		}

		err := ws.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			conn.log("listenWrite", nil, "WebSocket Write Error")
		}
	}
}

// Close will send close message and shutdown websocket connection
func (conn *WebSocketClient) Disconnect() {
	conn.ctxCancel()
	conn.closeWs()
}

func (c *WebSocketClient) Subscribe(channel string) error {
	return c.write(proto.NewSubscribeMessage(channel))
}

func (c *WebSocketClient) Unsubscribe(channel string) error {
	return c.write(proto.NewUnsubscribeMessage(channel))
}

func (c *WebSocketClient) Authenticate(token string) error {
	return c.write(proto.NewAuthenticateMessage(token))
}

func (c *WebSocketClient) Publish(channel string, data string) error {
	return c.write(proto.NewPublishMessage(channel, data))
}

// Write data to the websocket server
func (conn *WebSocketClient) write(msg proto.ClientMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
	defer cancel()

	for {
		select {
		case conn.sendBuf <- data:
			return nil
		case <-ctx.Done():
			return fmt.Errorf("context canceled")
		}
	}
}

// Close will send close message and shutdown websocket connection
func (conn *WebSocketClient) closeWs() {
	conn.mu.Lock()
	if conn.wsconn != nil {
		conn.wsconn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		conn.wsconn.Close()
		conn.wsconn = nil
	}
	conn.mu.Unlock()
}

func (conn *WebSocketClient) ping() {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			ws := conn.Connect()
			if ws == nil {
				continue
			}
			if err := conn.wsconn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(pingPeriod/2)); err != nil {
				conn.closeWs()
			}
		case <-conn.ctx.Done():
			return
		}
	}
}

// Log print log statement
// In real word I would recommend to use zerolog or any other solution
func (conn *WebSocketClient) log(f string, err error, msg string) {
	if err != nil {
		fmt.Printf("Error in func: %s, err: %v, msg: %s\n", f, err, msg)
	} else {
		fmt.Printf("Log in func: %s, %s\n", f, msg)
	}
}
