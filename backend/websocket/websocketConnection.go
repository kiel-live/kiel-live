package websocket

import (
	"encoding/binary"
	"errors"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/kiel-live/kiel-live/protocol"
	"github.com/pborman/uuid"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("origin", "WebsocketConnection")

type websocketConnection struct {
	Token    string
	Conn     *websocket.Conn
	Server   *Server
	AuthData protocol.ClientMessage

	writeLock sync.Mutex
	readLock  sync.Mutex
}

func newWebsocketConnection(w http.ResponseWriter, r *http.Request, s *Server) {
	c := &websocketConnection{
		Server: s,
		Token:  uuid.New(),
	}

	err := c.handshake(w, r)
	if err != nil {
		if c.Conn != nil {
			c.SendMessage(protocol.NewErrorMessage(protocol.ServerErrorMessage, err))
			c.Conn.Close()
		} else {
			http.Error(w, err.Error(), 500)
		}
	}
}

func (c *websocketConnection) writeConn(msg protocol.ClientMessage) error {
	c.writeLock.Lock()
	defer c.writeLock.Unlock()
	log.Debug(">>> ", msg)
	return c.Conn.WriteJSON(msg)
}

func (c *websocketConnection) readConn(v interface{}) error {
	c.readLock.Lock()
	defer c.readLock.Unlock()
	return c.Conn.ReadJSON(v)
}

func (c *websocketConnection) handshake(w http.ResponseWriter, r *http.Request) error {
	conn, err := c.Server.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		// websocket library already sends error message, nothing to do here
		return nil
	}
	c.Conn = conn

	log.Debugln("Client connecting ...")

	if c.Server.CanConnect != nil && !c.Server.CanConnect(c) {
		c.Close(4400, err.Error())
		return nil
	}

	defer c.Cleanup()

	hub := c.Server.hub
	err = hub.Connect(c)
	if err != nil {
		return err
	}

	log.Debugln("Client connected")

	c.Run()

	return nil
}

func (c *websocketConnection) Run() {
	hub := c.Server.hub

	m := protocol.ClientMessage{}
	for {
		err := c.readConn(&m)
		if err != nil {
			c.Close(4400, err.Error())
			break
		}

		switch m.Type() {
		case protocol.AuthMessage:
			if c.Server.CanAuthenticate != nil && !c.Server.CanAuthenticate(m) {
				c.SendMessage(protocol.NewErrorMessage(protocol.AuthFailedMessage, errors.New("Unauthorized")))
				continue
			}

			c.SendMessage(protocol.NewMessage(protocol.AuthOKMessage))

			c.AuthData = m

		case protocol.PublishMessage:
			channel := m.Channel()
			if c.Server.CanPublish != nil && !c.Server.CanPublish(c.AuthData, channel) {
				c.SendMessage(protocol.NewChannelErrorMessage(protocol.PublishErrorMessage, channel, errors.New("Access denied")))
				continue
			}

			data := m.Data()
			err := hub.Publish(c, channel, data)
			if err != nil {
				c.SendMessage(protocol.NewChannelErrorMessage(protocol.PublishErrorMessage, channel, err))
				continue
			}

			c.SendMessage(protocol.NewChannelMessage(protocol.PublishOKMessage, channel))

		case protocol.SubscribeMessage:
			channel := m.Channel()
			if c.Server.CanSubscribe != nil && !c.Server.CanSubscribe(c.AuthData, channel) {
				c.SendMessage(protocol.NewChannelErrorMessage(protocol.SubscribeErrorMessage, channel, errors.New("channel refused")))
				continue
			}

			err := hub.Subscribe(c, channel)
			if err != nil {
				c.SendMessage(protocol.NewChannelErrorMessage(protocol.SubscribeErrorMessage, channel, err))
				continue
			}

			c.SendMessage(protocol.NewChannelMessage(protocol.SubscribeOKMessage, channel))

		case protocol.UnsubscribeMessage:
			channel := m.Channel()

			err := hub.Unsubscribe(c, channel)
			if err != nil {
				c.SendMessage(protocol.NewChannelErrorMessage(protocol.UnsubscribeErrorMessage, channel, err))
				continue
			}

			c.SendMessage(protocol.NewChannelMessage(protocol.UnsubscribeOKMessage, channel))

		case protocol.PingMessage:
			// Do nothing

		default:
			c.SendMessage(protocol.NewMessage(protocol.UnknownMessage))
		}
	}
}

func (c *websocketConnection) Cleanup() {
	hub := c.Server.hub

	err := hub.Disconnect(c)
	if err != nil {
		c.SendMessage(protocol.NewErrorMessage(protocol.ServerErrorMessage, err))
	}

	c.Conn.Close()
}

func (c *websocketConnection) Close(code uint16, msg string) {
	payload := make([]byte, 2)
	binary.BigEndian.PutUint16(payload, code)
	payload = append(payload, []byte(msg)...)
	err := c.Conn.WriteMessage(websocket.CloseMessage, payload)
	if err != nil {
		log.Errorln("can't send message: %s", err)
	}

	err = c.Conn.Close()
	if err != nil {
		log.Errorln("can't close socket: %s", err)
	}

}

func (c *websocketConnection) Send(channel, message string) {
	err := c.writeConn(protocol.NewBroadcastMessage(channel, message))
	if err != nil {
		log.Errorln("can't send message: %s", err)
	}
}

func (c *websocketConnection) SendMessage(message protocol.ClientMessage) {
	err := c.writeConn(message)
	if err != nil {
		log.Errorln("can't send message: %s", err)
	}
}

func (c *websocketConnection) GetToken() string {
	return c.Token
}
