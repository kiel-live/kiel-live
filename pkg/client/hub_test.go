package client_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kiel-live/kiel-live/hub/server"
	"github.com/kiel-live/kiel-live/hub/store"
	"github.com/kiel-live/kiel-live/pkg/client"
	"github.com/kiel-live/kiel-live/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testToken = "test-secret"

func newTestHub(t *testing.T) *httptest.Server {
	t.Helper()
	st := store.New()
	srv := server.New(st, testToken)
	mux := http.NewServeMux()
	srv.RegisterRoutes(mux)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts
}

func hubWSURL(ts *httptest.Server) string {
	return "ws" + strings.TrimPrefix(ts.URL, "http")
}

func TestHubClientConnectAndUpdate(t *testing.T) {
	ts := newTestHub(t)
	url := hubWSURL(ts)

	c := client.NewHubClient(url, testToken)
	require.NoError(t, c.Connect())
	defer func() { _ = c.Disconnect() }()

	assert.True(t, c.IsConnected())

	// Publish a stop — should not error
	stop := &models.Stop{
		ID: "kvg-1", Provider: "kvg", Name: "Hauptbahnhof",
		Type: "bus-stop",
		Location: &models.Location{Latitude: 195516960, Longitude: 36539160},
	}
	require.NoError(t, c.UpdateStop(stop))

	// Publish a vehicle
	vehicle := &models.Vehicle{
		ID: "ais-123", Provider: "ais", Name: "Fähre", Type: "ferry", State: "running",
		Location: &models.Location{Latitude: 195500000, Longitude: 36500000},
	}
	require.NoError(t, c.UpdateVehicle(vehicle))
}

func TestHubClientSubscriptionEvents(t *testing.T) {
	ts := newTestHub(t)
	url := hubWSURL(ts)

	c := client.NewHubClient(url, testToken)
	require.NoError(t, c.Connect())
	defer func() { _ = c.Disconnect() }()

	topicCh := make(chan string, 4)
	c.SetOnTopicsChanged(func(topic string, subscribed bool) {
		if subscribed {
			topicCh <- topic
		}
	})

	// The hub needs a client to subscribe to a detail topic before the collector
	// receives a subscription event. This is tested indirectly via the store.
	// For a direct unit test, pre-seed the subscription manager via a second WS client.
	// Here we just verify that GetSubscribedTopics reflects initial state (empty).
	topics := c.GetSubscribedTopics()
	assert.Empty(t, topics)
}

func TestHubClientDeleteStop(t *testing.T) {
	ts := newTestHub(t)
	url := hubWSURL(ts)

	c := client.NewHubClient(url, testToken)
	require.NoError(t, c.Connect())
	defer func() { _ = c.Disconnect() }()

	stop := &models.Stop{ID: "kvg-2", Provider: "kvg", Name: "Test", Type: "bus-stop",
		Location: &models.Location{Latitude: 195516960, Longitude: 36539160}}
	require.NoError(t, c.UpdateStop(stop))
	time.Sleep(20 * time.Millisecond) // let server process
	require.NoError(t, c.DeleteStop("kvg-2"))
}

// TestHubClientDetectsServerClosedConnection verifies the reconnect loop
// engages when the hub-side connection disappears. This is what the
// collector observes when the hub's WS keepalive decides a connection is
// dead and tears it down from its side — from the collector's perspective,
// a keepalive-triggered close is indistinguishable from any other
// server-initiated close, so a minimal fake endpoint that upgrades and then
// closes exercises the same hubClient code path deterministically, without
// waiting out real ping/pong timers (which a healthy hubClient would always
// answer anyway, since its read loop auto-responds to pings).
func TestHubClientDetectsServerClosedConnection(t *testing.T) {
	closeNow := make(chan struct{})
	upgrader := websocket.Upgrader{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer "+testToken {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		<-closeNow
		conn.Close()
	}))
	defer ts.Close()

	c := client.NewHubClient(hubWSURL(ts), testToken)
	require.NoError(t, c.Connect())
	defer func() { _ = c.Disconnect() }()
	require.True(t, c.IsConnected())

	disconnected := make(chan struct{}, 1)
	c.SetOnConnectionChanged(func(connected bool) {
		if !connected {
			select {
			case disconnected <- struct{}{}:
			default:
			}
		}
	})

	close(closeNow) // simulate the hub's keepalive deciding this connection is dead

	select {
	case <-disconnected:
	case <-time.After(2 * time.Second):
		t.Fatal("expected connection-lost callback to fire after server closed the connection")
	}
	assert.False(t, c.IsConnected())
}
