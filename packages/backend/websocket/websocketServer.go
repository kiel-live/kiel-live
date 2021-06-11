package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/kiel-live/kiel-live/packages/backend/hub"
	proto "github.com/kiel-live/kiel-live/packages/pub-sub-proto"
)

type Server struct {
	// Invoked upon connection, can be used to do pre-connect checks.
	CanConnect func(conn *websocketConnection) bool

	// Invoked upon authentication, can be used to enforce access control.
	CanAuthenticate func(authMessage proto.ClientMessage) bool

	// Invoked upon channel subscription, can be used to enforce access control
	// for channels.
	CanSubscribe func(authData map[string]interface{}, channel string) bool

	// Invoked upon channel publish, can be used to enforce access control
	// for channels.
	CanPublish func(authData map[string]interface{}, channel string) bool

	// Can be set to allow CORS requests.
	CheckOrigin func(r *http.Request) bool

	// Can be used to configure buffer sizes etc.
	// See http://godoc.org/github.com/gorilla/websocket#Upgrader
	Upgrader websocket.Upgrader

	hub *hub.Hub
}

func NewServer(hub *hub.Hub, token string) *Server {
	return &Server{
		hub: hub,
		CanAuthenticate: func(authMessage proto.ClientMessage) bool {
			return authMessage.Data() == token
		},
		CanPublish: func(authData map[string]interface{}, channel string) bool {
			return authData != nil // require authentication
		},
	}
}

func (server *Server) WebsocketEndpoint(w http.ResponseWriter, r *http.Request) {
	newWebsocketConnection(w, r, server)
}
