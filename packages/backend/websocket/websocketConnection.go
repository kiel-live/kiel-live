package websocket

import (
	"encoding/binary"
	"errors"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/kiel-live/kiel-live/backend/proto"
	"github.com/pborman/uuid"
)

type websocketConnection struct {
	Token    string
	Conn     *websocket.Conn
	Server   *WebsocketServer
	AuthData proto.ClientMessage

	write_lock sync.Mutex
	read_lock  sync.Mutex
}

func newWebsocketConnection(w http.ResponseWriter, r *http.Request, s *WebsocketServer) {
	conn := &websocketConnection{
		Server: s,
		Token:  uuid.New(),
	}

	err := conn.handshake(w, r)
	if err != nil {
		if conn.Conn != nil {
			conn.Conn.WriteJSON(proto.NewErrorMessage(proto.ServerErrorMessage, err))
			conn.Conn.Close()
		} else {
			http.Error(w, err.Error(), 500)
		}
	}
}

func (c *websocketConnection) writeConn(msg proto.ClientMessage) error {
	c.write_lock.Lock()
	defer c.write_lock.Unlock()
	return c.Conn.WriteJSON(msg)
}

func (c *websocketConnection) readConn(v interface{}) error {
	c.read_lock.Lock()
	defer c.read_lock.Unlock()
	return c.Conn.ReadJSON(v)
}

func (c *websocketConnection) handshake(w http.ResponseWriter, r *http.Request) error {
	conn, err := c.Server.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		// websocket library already sends error message, nothing to do here
		return nil
	}
	c.Conn = conn

	err = c.readConn(&c.AuthData)
	if err != nil {
		c.Close(4400, err.Error())
		return nil
	}

	// Expect auth packet first.
	if c.AuthData.Type() != proto.AuthMessage {
		c.writeConn(proto.NewErrorMessage(proto.AuthFailedMessage, errors.New("auth expected")))
		c.Close(4401, "Auth expected")
		return nil
	}

	if c.Server.CanConnect != nil && !c.Server.CanConnect(c.AuthData) {
		c.writeConn(proto.NewErrorMessage(proto.AuthFailedMessage, errors.New("Unauthorized")))
		c.Close(4401, "Unauthorized")
		return nil
	}

	defer c.Cleanup()

	err = c.writeConn(proto.NewMessage(proto.AuthOKMessage))
	if err != nil {
		return err
	}

	hub := c.Server.hub
	err = hub.Connect(c)
	if err != nil {
		return err
	}

	c.Run()

	return nil
}

func (c *websocketConnection) Run() {
	hub := c.Server.hub

	m := proto.ClientMessage{}
	for {
		err := c.readConn(&m)
		if err != nil {
			c.Close(4400, err.Error())
			break
		}

		switch m.Type() {
		case proto.SubscribeMessage:
			channel := m.Channel()
			if c.Server.CanSubscribe != nil && !c.Server.CanSubscribe(c.AuthData, channel) {
				c.writeConn(proto.NewChannelErrorMessage(proto.SubscribeErrorMessage, channel, errors.New("channel refused")))
				continue
			}

			err := hub.Subscribe(c, channel)
			if err != nil {
				c.writeConn(proto.NewChannelErrorMessage(proto.SubscribeErrorMessage, channel, err))
			} else {
				c.writeConn(proto.NewChannelMessage(proto.SubscribeOKMessage, channel))
			}

		case proto.UnsubscribeMessage:
			channel := m.Channel()

			err := hub.Unsubscribe(c, channel)
			if err != nil {
				c.writeConn(proto.NewChannelErrorMessage(proto.UnsubscribeErrorMessage, channel, err))
				continue
			}
			c.writeConn(proto.NewChannelMessage(proto.UnsubscribeOKMessage, channel))

		case proto.PingMessage:
			// Do nothing

		default:
			c.writeConn(proto.NewMessage(proto.UnknownMessage))
		}
	}
}

func (c *websocketConnection) Cleanup() {
	hub := c.Server.hub

	err := hub.Disconnect(c)
	if err != nil {
		c.writeConn(proto.NewErrorMessage(proto.ServerErrorMessage, err))
	}

	c.Conn.Close()
}

func (c *websocketConnection) Close(code uint16, msg string) {
	payload := make([]byte, 2)
	binary.BigEndian.PutUint16(payload, code)
	payload = append(payload, []byte(msg)...)
	c.Conn.WriteMessage(websocket.CloseMessage, payload)
	c.Conn.Close()
}

func (c *websocketConnection) Send(channel, message string) {
	c.writeConn(proto.NewBroadcastMessage(channel, message))
}

func (c *websocketConnection) Process(t string, args []string) {
	panic("Websocket connections don't use control messages!")
}

func (c *websocketConnection) GetToken() string {
	return c.Token
}
