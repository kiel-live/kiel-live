package websocket

import (
	"encoding/binary"
	"errors"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	proto "github.com/kiel-live/kiel-live/packages/pub-sub-proto"
	"github.com/pborman/uuid"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("origin", "WebsocketConnection")

type websocketConnection struct {
	Token    string
	Conn     *websocket.Conn
	Server   *Server
	AuthData proto.ClientMessage

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
			c.SendMessage(proto.NewErrorMessage(proto.ServerErrorMessage, err))
			c.Conn.Close()
		} else {
			http.Error(w, err.Error(), 500)
		}
	}
}

func (c *websocketConnection) writeConn(msg proto.ClientMessage) error {
	c.writeLock.Lock()
	defer c.writeLock.Unlock()
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

	m := proto.ClientMessage{}
	for {
		err := c.readConn(&m)
		if err != nil {
			c.Close(4400, err.Error())
			break
		}

		switch m.Type() {
		case proto.AuthMessage:
			if c.Server.CanAuthenticate != nil && !c.Server.CanAuthenticate(m) {
				c.SendMessage(proto.NewErrorMessage(proto.AuthFailedMessage, errors.New("Unauthorized")))
				continue
			}

			c.SendMessage(proto.NewMessage(proto.AuthOKMessage))

			c.AuthData = m

		case proto.PublishMessage:
			channel := m.Channel()
			if c.Server.CanPublish != nil && !c.Server.CanPublish(c.AuthData, channel) {
				c.SendMessage(proto.NewChannelErrorMessage(proto.PublishErrorMessage, channel, errors.New("Access denied")))
				continue
			}

			data := m.Data()
			err := hub.Publish(c, channel, data)
			if err != nil {
				c.SendMessage(proto.NewChannelErrorMessage(proto.PublishErrorMessage, channel, err))
				continue
			}

			c.SendMessage(proto.NewChannelMessage(proto.PublishOKMessage, channel))

		case proto.SubscribeMessage:
			channel := m.Channel()
			if c.Server.CanSubscribe != nil && !c.Server.CanSubscribe(c.AuthData, channel) {
				c.SendMessage(proto.NewChannelErrorMessage(proto.SubscribeErrorMessage, channel, errors.New("channel refused")))
				continue
			}

			err := hub.Subscribe(c, channel)
			if err != nil {
				c.SendMessage(proto.NewChannelErrorMessage(proto.SubscribeErrorMessage, channel, err))
				continue
			}

			c.SendMessage(proto.NewChannelMessage(proto.SubscribeOKMessage, channel))

		case proto.UnsubscribeMessage:
			channel := m.Channel()

			err := hub.Unsubscribe(c, channel)
			if err != nil {
				c.SendMessage(proto.NewChannelErrorMessage(proto.UnsubscribeErrorMessage, channel, err))
				continue
			}

			c.SendMessage(proto.NewChannelMessage(proto.UnsubscribeOKMessage, channel))

		case proto.PingMessage:
			// Do nothing

		default:
			c.SendMessage(proto.NewMessage(proto.UnknownMessage))
		}
	}
}

func (c *websocketConnection) Cleanup() {
	hub := c.Server.hub

	err := hub.Disconnect(c)
	if err != nil {
		c.SendMessage(proto.NewErrorMessage(proto.ServerErrorMessage, err))
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
	err := c.writeConn(proto.NewBroadcastMessage(channel, message))
	if err != nil {
		log.Errorln("can't send message: %s", err)
	}
}

func (c *websocketConnection) SendMessage(message proto.ClientMessage) {
	err := c.writeConn(message)
	if err != nil {
		log.Errorln("can't send message: %s", err)
	}
}

func (c *websocketConnection) GetToken() string {
	return c.Token
}
