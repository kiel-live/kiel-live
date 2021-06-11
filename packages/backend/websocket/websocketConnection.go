package websocket

import (
	"encoding/binary"
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	proto "github.com/kiel-live/kiel-live/packages/pub-sub-proto"
	"github.com/pborman/uuid"
)

type websocketConnection struct {
	Token    string
	Conn     *websocket.Conn
	Server   *WebsocketServer
	AuthData proto.ClientMessage

	writeLock sync.Mutex
	readLock  sync.Mutex
}

func newWebsocketConnection(w http.ResponseWriter, r *http.Request, s *WebsocketServer) {
	conn := &websocketConnection{
		Server: s,
		Token:  uuid.New(),
	}

	err := conn.handshake(w, r)
	if err != nil {
		if conn.Conn != nil {
			err = conn.Conn.WriteJSON(proto.NewErrorMessage(proto.ServerErrorMessage, err))
			if err != nil {
				log.Printf("can't send message: %s", err)
			}
			conn.Conn.Close()
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

	log.Printf("Client connecting ...")

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

	log.Printf("Client connected")

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
				err := c.writeConn(proto.NewErrorMessage(proto.AuthFailedMessage, errors.New("Unauthorized")))
				if err != nil {
					log.Printf("can't send message: %s", err)
				}
				continue
			}

			err := c.writeConn(proto.NewMessage(proto.AuthOKMessage))
			if err != nil {
				log.Printf("can't send message: %s", err)
			}

			c.AuthData = m

		case proto.PublishMessage:
			channel := m.Channel()
			if c.Server.CanPublish != nil && !c.Server.CanPublish(c.AuthData, channel) {
				err := c.writeConn(proto.NewChannelErrorMessage(proto.PublishErrorMessage, channel, errors.New("channel refused")))
				if err != nil {
					log.Printf("can't send message: %s", err)
				}
				continue
			}

			data := m.Data()
			err := hub.Publish(c, channel, data)
			if err != nil {
				err := c.writeConn(proto.NewChannelErrorMessage(proto.PublishErrorMessage, channel, err))
				if err != nil {
					log.Printf("can't send message: %s", err)
				}
				continue
			}

			err = c.writeConn(proto.NewChannelMessage(proto.PublishOKMessage, channel))
			if err != nil {
				log.Printf("can't send message: %s", err)
			}

		case proto.SubscribeMessage:
			channel := m.Channel()
			if c.Server.CanSubscribe != nil && !c.Server.CanSubscribe(c.AuthData, channel) {
				err := c.writeConn(proto.NewChannelErrorMessage(proto.SubscribeErrorMessage, channel, errors.New("channel refused")))
				if err != nil {
					log.Printf("can't send message: %s", err)
				}
				continue
			}

			err := hub.Subscribe(c, channel)
			if err != nil {
				err := c.writeConn(proto.NewChannelErrorMessage(proto.SubscribeErrorMessage, channel, err))
				if err != nil {
					log.Printf("can't send message: %s", err)
				}
				continue
			}

			err = c.writeConn(proto.NewChannelMessage(proto.SubscribeOKMessage, channel))
			if err != nil {
				log.Printf("can't send message: %s", err)
			}

		case proto.UnsubscribeMessage:
			channel := m.Channel()

			err := hub.Unsubscribe(c, channel)
			if err != nil {
				err := c.writeConn(proto.NewChannelErrorMessage(proto.UnsubscribeErrorMessage, channel, err))
				if err != nil {
					log.Printf("can't send message: %s", err)
				}
				continue
			}

			err = c.writeConn(proto.NewChannelMessage(proto.UnsubscribeOKMessage, channel))
			if err != nil {
				log.Printf("can't send message: %s", err)
			}

		case proto.PingMessage:
			// Do nothing

		default:
			err := c.writeConn(proto.NewMessage(proto.UnknownMessage))
			if err != nil {
				log.Printf("can't send message: %s", err)
			}
		}
	}
}

func (c *websocketConnection) Cleanup() {
	hub := c.Server.hub

	err := hub.Disconnect(c)
	if err != nil {
		err := c.writeConn(proto.NewErrorMessage(proto.ServerErrorMessage, err))
		if err != nil {
			log.Printf("can't send message: %s", err)
		}
	}

	c.Conn.Close()
}

func (c *websocketConnection) Close(code uint16, msg string) {
	payload := make([]byte, 2)
	binary.BigEndian.PutUint16(payload, code)
	payload = append(payload, []byte(msg)...)
	err := c.Conn.WriteMessage(websocket.CloseMessage, payload)
	if err != nil {
		log.Printf("can't send message: %s", err)
	}

	err = c.Conn.Close()
	if err != nil {
		log.Printf("can't close socket: %s", err)
	}

}

func (c *websocketConnection) Send(channel, message string) {
	err := c.writeConn(proto.NewBroadcastMessage(channel, message))
	if err != nil {
		log.Printf("can't send message: %s", err)
	}
}

func (c *websocketConnection) GetToken() string {
	return c.Token
}
