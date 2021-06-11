package webserver

import (
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/kiel-live/kiel-live/packages/backend/websocket"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("origin", "WebServer")

type WebServer struct {
	server          *http.Server
	websocketServer *websocket.Server
}

func NewWebServer(websocketServer *websocket.Server) *WebServer {
	return &WebServer{
		websocketServer: websocketServer,
	}
}

func (webServer *WebServer) serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	http.ServeFile(w, r, "home.html")
}

func (webServer *WebServer) router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", webServer.serveHome)

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		webServer.websocketServer.WebsocketEndpoint(w, r)
	})

	return mux
}

func (webServer *WebServer) Listen(port int) error {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	webServer.server = &http.Server{Addr: addr, Handler: webServer.router()}

	log.Infoln("🚀 Backend running at:")
	log.Infof("> Local: http://localhost:%d\n", port)

	return webServer.server.Serve(listener)
}

func (webServer *WebServer) Close() error {
	if webServer.server != nil {
		err := webServer.server.Close()
		if err == nil {
			webServer.server = nil
		}
		return err
	}

	return errors.New("can't close server as no one was open")
}
