package server_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kiel-live/kiel-live/hub/server"
	"github.com/kiel-live/kiel-live/hub/store"
	"github.com/kiel-live/kiel-live/pkg/models"
	"github.com/kiel-live/kiel-live/pkg/protocol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── fixtures ──────────────────────────────────────────────────────────────────

const testToken = "test-secret"

var kielLoc = &models.Location{
	Latitude:  195516960, // 54.31°N * 3_600_000
	Longitude: 36539160,  // 10.15°E * 3_600_000
}

var kielViewport = protocol.ViewportData{North: 54.35, South: 54.27, East: 10.22, West: 10.08}

var hamburgLoc = &models.Location{
	Latitude:  192780000, // 53.55°N * 3_600_000
	Longitude: 36000000,  // 10.00°E * 3_600_000
}

var hamburgViewport = protocol.ViewportData{North: 53.60, South: 53.50, East: 10.05, West: 9.95}

var tokyoLoc = &models.Location{
	Latitude:  128480400, // 35.69°N * 3_600_000
	Longitude: 502891200, // 139.69°E * 3_600_000
}

// ── wsClient: single-reader pump ─────────────────────────────────────────────

// wsClient wraps a WebSocket and runs a background goroutine that feeds all
// inbound messages into a channel. This avoids the gorilla/websocket bug where
// a SetReadDeadline timeout permanently poisons future ReadMessage calls.
type wsClient struct {
	ws    *websocket.Conn
	msgCh chan protocol.Envelope
}

func newWSClient(t *testing.T, ws *websocket.Conn) *wsClient {
	t.Helper()
	c := &wsClient{
		ws:    ws,
		msgCh: make(chan protocol.Envelope, 128),
	}
	go func() {
		for {
			_, data, err := ws.ReadMessage()
			if err != nil {
				close(c.msgCh)
				return
			}
			var env protocol.Envelope
			if jsonErr := json.Unmarshal(data, &env); jsonErr == nil {
				c.msgCh <- env
			}
		}
	}()
	t.Cleanup(func() { _ = ws.Close() })
	return c
}

// recv reads the next envelope, failing after 2 seconds.
func (c *wsClient) recv(t *testing.T) protocol.Envelope {
	t.Helper()
	select {
	case env, ok := <-c.msgCh:
		require.True(t, ok, "connection closed unexpectedly")
		return env
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for message")
		return protocol.Envelope{}
	}
}

// noRecv asserts no message arrives within 200 ms.
func (c *wsClient) noRecv(t *testing.T) {
	t.Helper()
	select {
	case env, ok := <-c.msgCh:
		if ok {
			t.Fatalf("expected no message but got method=%s topic=%s", env.Method, env.Topic)
		}
	case <-time.After(200 * time.Millisecond):
		// expected
	}
}

// drainSnapshots reads snapshot messages until 300 ms of silence.
// The new protocol sends one snapshot per cell per kind.
func (c *wsClient) drainSnapshots(t *testing.T) []protocol.Envelope {
	t.Helper()
	var snaps []protocol.Envelope
	for {
		select {
		case env, ok := <-c.msgCh:
			if !ok {
				return snaps
			}
			require.Equal(t, protocol.MethodSnapshot, env.Method,
				"expected snapshot, got method=%s topic=%s", env.Method, env.Topic)
			snaps = append(snaps, env)
		case <-time.After(300 * time.Millisecond):
			return snaps
		}
	}
}

// ── helpers ───────────────────────────────────────────────────────────────────

type testServer struct {
	*httptest.Server
	store *store.Store
}

func newTestServer(t *testing.T) *testServer {
	t.Helper()
	st := store.New()
	srv := server.New(st, testToken)
	mux := http.NewServeMux()
	srv.RegisterRoutes(mux)
	hs := httptest.NewServer(mux)
	t.Cleanup(hs.Close)
	return &testServer{Server: hs, store: st}
}

func wsURL(ts *testServer, path string) string {
	return "ws" + strings.TrimPrefix(ts.URL, "http") + path
}

func connectClient(t *testing.T, ts *testServer) *wsClient {
	t.Helper()
	ws, _, err := websocket.DefaultDialer.Dial(wsURL(ts, "/ws/client"), nil)
	require.NoError(t, err)
	return newWSClient(t, ws)
}

func connectCollector(t *testing.T, ts *testServer) *wsClient {
	t.Helper()
	hdr := http.Header{"Authorization": {"Bearer " + testToken}}
	ws, _, err := websocket.DefaultDialer.Dial(wsURL(ts, "/ws/collector"), hdr)
	require.NoError(t, err)
	return newWSClient(t, ws)
}

func sendRaw(t *testing.T, ws *websocket.Conn, env protocol.Envelope) {
	t.Helper()
	data, err := json.Marshal(env)
	require.NoError(t, err)
	require.NoError(t, ws.WriteMessage(websocket.TextMessage, data))
}

func sendData(t *testing.T, col *wsClient, method protocol.Method, topic string, v any) {
	t.Helper()
	b, _ := json.Marshal(v)
	raw := json.RawMessage(b)
	sendRaw(t, col.ws, protocol.Envelope{Method: method, Topic: topic, Data: &raw})
}

func sendViewport(t *testing.T, cli *wsClient, vp protocol.ViewportData) {
	t.Helper()
	b, _ := json.Marshal(vp)
	raw := json.RawMessage(b)
	sendRaw(t, cli.ws, protocol.Envelope{Method: protocol.MethodViewport, Data: &raw})
}

func sendSubscribe(t *testing.T, cli *wsClient, topic string) {
	t.Helper()
	sendRaw(t, cli.ws, protocol.Envelope{Method: protocol.MethodSubscribe, Topic: topic})
}

func sendUnsubscribe(t *testing.T, cli *wsClient, topic string) {
	t.Helper()
	sendRaw(t, cli.ws, protocol.Envelope{Method: protocol.MethodUnsubscribe, Topic: topic})
}

// subscribeMap sends subscribe for both map.stop and map.vehicle then the viewport.
func subscribeMap(t *testing.T, cli *wsClient, vp protocol.ViewportData) {
	t.Helper()
	sendSubscribe(t, cli, "map.stop")
	sendSubscribe(t, cli, "map.vehicle")
	sendViewport(t, cli, vp)
}

// ── tests ─────────────────────────────────────────────────────────────────────

// TestCollectorAuth rejects connections without a valid token.
func TestCollectorAuth(t *testing.T) {
	ts := newTestServer(t)

	_, resp, _ := websocket.DefaultDialer.Dial(wsURL(ts, "/ws/collector"),
		http.Header{"Authorization": {"Bearer wrong"}})
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	_, resp, _ = websocket.DefaultDialer.Dial(wsURL(ts, "/ws/collector"), nil)
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

// TestNoMapDataWithoutSubscribe verifies that a viewport alone does not trigger
// any data delivery — the client must explicitly subscribe to map.stop / map.vehicle.
func TestNoMapDataWithoutSubscribe(t *testing.T) {
	ts := newTestServer(t)
	cli := connectClient(t, ts)
	col := connectCollector(t, ts)

	sendViewport(t, cli, kielViewport)
	time.Sleep(30 * time.Millisecond)

	stop := &models.Stop{
		ID: "kvg-1", Provider: "kvg", Name: "Hauptbahnhof",
		Type: models.StopTypeBusStop, Location: kielLoc,
	}
	sendData(t, col, protocol.MethodUpdate, "stop.kvg-1", stop)

	cli.noRecv(t)
}

// TestClientReceivesEmptySnapshots verifies that subscribing to map.stop and
// map.vehicle with a viewport yields per-cell snapshots (empty when store is empty).
func TestClientReceivesEmptySnapshots(t *testing.T) {
	ts := newTestServer(t)
	cli := connectClient(t, ts)

	subscribeMap(t, cli, kielViewport)

	snaps := cli.drainSnapshots(t)
	require.NotEmpty(t, snaps, "must receive at least one snapshot")

	for _, s := range snaps {
		assert.True(t, s.Topic == "map.stop" || s.Topic == "map.vehicle",
			"snapshot topic must be map.stop or map.vehicle, got %s", s.Topic)
		var items []json.RawMessage
		require.NoError(t, json.Unmarshal(*s.Data, &items))
		assert.Empty(t, items, "store is empty, snapshot must be empty")
	}
}

// TestMarkerDeliveredToViewport publishes a stop inside the viewport and
// verifies the client receives a sparse marker (no departures).
func TestMarkerDeliveredToViewport(t *testing.T) {
	ts := newTestServer(t)
	cli := connectClient(t, ts)
	col := connectCollector(t, ts)

	subscribeMap(t, cli, kielViewport)
	cli.drainSnapshots(t)

	stop := &models.Stop{
		ID: "kvg-1", Provider: "kvg", Name: "Hauptbahnhof",
		Type: models.StopTypeBusStop, Location: kielLoc,
		Departures: []*models.StopDepartures{{Name: "Bus 1", Direction: "North"}},
	}
	sendData(t, col, protocol.MethodUpdate, "stop.kvg-1", stop)

	e := cli.recv(t)
	assert.Equal(t, protocol.MethodUpdate, e.Method)
	assert.Equal(t, "map.stop.kvg-1", e.Topic)

	var marker models.Stop
	require.NoError(t, json.Unmarshal(*e.Data, &marker))
	assert.Equal(t, "kvg-1", marker.ID)
	assert.Equal(t, "Hauptbahnhof", marker.Name)
	assert.Nil(t, marker.Departures, "marker must not include departures")
}

// TestMarkerNotDeliveredWithoutMapSubscribe verifies that a client without
// map.stop subscribed does not receive any map data.
func TestMarkerNotDeliveredWithoutMapSubscribe(t *testing.T) {
	ts := newTestServer(t)
	cli := connectClient(t, ts)
	col := connectCollector(t, ts)

	sendViewport(t, cli, kielViewport)

	stop := &models.Stop{
		ID: "kvg-1", Provider: "kvg", Name: "Hauptbahnhof",
		Type: models.StopTypeBusStop, Location: kielLoc,
	}
	sendData(t, col, protocol.MethodUpdate, "stop.kvg-1", stop)

	cli.noRecv(t)
}

// TestMarkerNotDeliveredOutsideViewport verifies viewport scoping.
func TestMarkerNotDeliveredOutsideViewport(t *testing.T) {
	ts := newTestServer(t)
	cli := connectClient(t, ts)
	col := connectCollector(t, ts)

	subscribeMap(t, cli, kielViewport)
	cli.drainSnapshots(t)

	tokyoStop := &models.Stop{
		ID: "gtfs-tokyo-1", Provider: "gtfs", Name: "Tokyo Station",
		Type: models.StopTypeBusStop, Location: tokyoLoc,
	}
	sendData(t, col, protocol.MethodUpdate, "stop.gtfs-tokyo-1", tokyoStop)

	cli.noRecv(t)
}

// TestUnsubscribeMapStopStopsDelivery verifies that unsubscribing from map.stop
// stops future map data delivery.
func TestUnsubscribeMapStopStopsDelivery(t *testing.T) {
	ts := newTestServer(t)
	cli := connectClient(t, ts)
	col := connectCollector(t, ts)

	subscribeMap(t, cli, kielViewport)
	cli.drainSnapshots(t)

	stop := &models.Stop{
		ID: "kvg-1", Provider: "kvg", Name: "A",
		Type: models.StopTypeBusStop, Location: kielLoc,
	}
	sendData(t, col, protocol.MethodUpdate, "stop.kvg-1", stop)
	markerEnv := cli.recv(t)
	assert.Equal(t, "map.stop.kvg-1", markerEnv.Topic)

	// Unsubscribe — server sends deletes for all visible stops.
	sendUnsubscribe(t, cli, "map.stop")
	deleteEnv := cli.recv(t)
	assert.Equal(t, protocol.MethodDelete, deleteEnv.Method)
	assert.Equal(t, "map.stop.kvg-1", deleteEnv.Topic)

	// Further updates must not reach the client.
	stop2 := *stop
	stop2.Name = "B"
	sendData(t, col, protocol.MethodUpdate, "stop.kvg-1", &stop2)
	cli.noRecv(t)
}

// TestDetailSubscription checks that:
//   - subscribing to stop.<id> sends the full entity immediately
//   - the collector is notified via the subscriptions topic
//   - the marker (map.*) stays sparse
func TestDetailSubscription(t *testing.T) {
	ts := newTestServer(t)
	col := connectCollector(t, ts)

	sendRaw(t, col.ws, protocol.Envelope{Method: protocol.MethodSubscribe, Topic: protocol.TopicSubscriptions})

	cli := connectClient(t, ts)
	subscribeMap(t, cli, kielViewport)
	cli.drainSnapshots(t)

	stop := &models.Stop{
		ID: "kvg-2", Provider: "kvg", Name: "Wik", Type: models.StopTypeFerryStop,
		Location:   kielLoc,
		Departures: []*models.StopDepartures{{Name: "Fähre", Direction: "Süd"}},
	}
	sendData(t, col, protocol.MethodUpdate, "stop.kvg-2", stop)

	markerEnv := cli.recv(t)
	assert.Equal(t, "map.stop.kvg-2", markerEnv.Topic)
	var markerStop models.Stop
	require.NoError(t, json.Unmarshal(*markerEnv.Data, &markerStop))
	assert.Nil(t, markerStop.Departures)

	sendSubscribe(t, cli, "stop.kvg-2")

	subEnv := col.recv(t)
	assert.Equal(t, protocol.TopicSubscriptions, subEnv.Topic)
	var se protocol.SubscriptionEvent
	require.NoError(t, json.Unmarshal(*subEnv.Data, &se))
	assert.Equal(t, "stop.kvg-2", se.Topic)
	assert.True(t, se.Subscribed)

	detailEnv := cli.recv(t)
	assert.Equal(t, "stop.kvg-2", detailEnv.Topic)
	var fullStop models.Stop
	require.NoError(t, json.Unmarshal(*detailEnv.Data, &fullStop))
	require.NotNil(t, fullStop.Departures)
	assert.Len(t, fullStop.Departures, 1)
}

// TestCollectorSubscriptionUnsubscribeEvent verifies the full sub/unsub cycle.
func TestCollectorSubscriptionUnsubscribeEvent(t *testing.T) {
	ts := newTestServer(t)
	col := connectCollector(t, ts)

	sendRaw(t, col.ws, protocol.Envelope{Method: protocol.MethodSubscribe, Topic: protocol.TopicSubscriptions})

	cli := connectClient(t, ts)
	subscribeMap(t, cli, kielViewport)
	cli.drainSnapshots(t)

	stop := &models.Stop{ID: "kvg-3", Provider: "kvg", Name: "Gaarden", Type: models.StopTypeBusStop, Location: kielLoc}
	sendData(t, col, protocol.MethodUpdate, "stop.kvg-3", stop)

	cli.recv(t)             // map.stop.kvg-3 update
	sendSubscribe(t, cli, "stop.kvg-3")

	subOn := col.recv(t)   // subscriptions: subscribed=true
	cli.recv(t)             // detail entity sent to client

	var seOn protocol.SubscriptionEvent
	require.NoError(t, json.Unmarshal(*subOn.Data, &seOn))
	assert.True(t, seOn.Subscribed)

	sendUnsubscribe(t, cli, "stop.kvg-3")

	subOff := col.recv(t)
	var seOff protocol.SubscriptionEvent
	require.NoError(t, json.Unmarshal(*subOff.Data, &seOff))
	assert.Equal(t, "stop.kvg-3", seOff.Topic)
	assert.False(t, seOff.Subscribed)
}

// TestViewportChangeDelta verifies that panning produces enter/leave deltas.
func TestViewportChangeDelta(t *testing.T) {
	ts := newTestServer(t)
	cli := connectClient(t, ts)
	col := connectCollector(t, ts)

	subscribeMap(t, cli, kielViewport)
	cli.drainSnapshots(t)

	kielStop := &models.Stop{
		ID: "kvg-kiel", Provider: "kvg", Name: "Kiel HBF",
		Type: models.StopTypeBusStop, Location: kielLoc,
	}
	hamburgStop := &models.Stop{
		ID: "kvg-hamburg", Provider: "kvg", Name: "Hamburg HBF",
		Type: models.StopTypeBusStop, Location: hamburgLoc,
	}

	sendData(t, col, protocol.MethodUpdate, "stop.kvg-kiel", kielStop)
	e := cli.recv(t)
	assert.Equal(t, "map.stop.kvg-kiel", e.Topic)

	sendData(t, col, protocol.MethodUpdate, "stop.kvg-hamburg", hamburgStop)
	time.Sleep(30 * time.Millisecond)

	// Pan to Hamburg.
	sendViewport(t, cli, hamburgViewport)

	// Collect all messages for 500 ms. Kiel stop should be deleted (explicit delete
	// message). Hamburg stop was already in the store, so it arrives in a snapshot.
	deletedTopics := map[string]bool{}
	snapshotIDs := map[string]bool{}
	deadline := time.After(500 * time.Millisecond)
loop:
	for {
		select {
		case env, ok := <-cli.msgCh:
			if !ok {
				break loop
			}
			switch env.Method {
			case protocol.MethodDelete:
				deletedTopics[env.Topic] = true
			case protocol.MethodSnapshot:
				var items []json.RawMessage
				if env.Data != nil {
					_ = json.Unmarshal(*env.Data, &items)
					for _, item := range items {
						var stop models.Stop
						if json.Unmarshal(item, &stop) == nil && stop.ID != "" {
							snapshotIDs[stop.ID] = true
						}
					}
				}
			}
		case <-deadline:
			break loop
		}
	}

	assert.True(t, deletedTopics["map.stop.kvg-kiel"], "Kiel stop should be deleted on pan away")
	assert.True(t, snapshotIDs["kvg-hamburg"], "Hamburg stop should appear in snapshot on pan into its viewport")
}

// TestCollectorDisconnectKeepsEntities verifies entities survive a collector disconnect.
func TestCollectorDisconnectKeepsEntities(t *testing.T) {
	ts := newTestServer(t)

	cli := connectClient(t, ts)
	col := connectCollector(t, ts)

	subscribeMap(t, cli, kielViewport)
	cli.drainSnapshots(t)

	stop := &models.Stop{
		ID: "kvg-temp", Provider: "kvg", Name: "Temp Stop",
		Type: models.StopTypeBusStop, Location: kielLoc,
	}
	sendData(t, col, protocol.MethodUpdate, "stop.kvg-temp", stop)

	markerEnv := cli.recv(t)
	assert.Equal(t, "map.stop.kvg-temp", markerEnv.Topic)

	_ = col.ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	col.ws.Close()

	cli.noRecv(t)
	assert.NotNil(t, ts.store.GetStop("kvg-temp"))
}

// TestNoOpSuppression verifies that identical consecutive vehicle positions
// don't produce duplicate marker updates.
func TestNoOpSuppression(t *testing.T) {
	ts := newTestServer(t)
	cli := connectClient(t, ts)
	col := connectCollector(t, ts)

	subscribeMap(t, cli, kielViewport)
	cli.drainSnapshots(t)

	vehicle := &models.Vehicle{
		ID: "ais-1", Provider: "ais", Name: "Fähre", Type: models.VehicleTypeFerry,
		State: "running", Location: kielLoc,
	}

	sendData(t, col, protocol.MethodUpdate, "vehicle.ais-1", vehicle)
	e := cli.recv(t)
	assert.Equal(t, protocol.MethodUpdate, e.Method)
	assert.Equal(t, "map.vehicle.ais-1", e.Topic)

	sendData(t, col, protocol.MethodUpdate, "vehicle.ais-1", vehicle)
	cli.noRecv(t)
}

// TestGTFSIdWithDots verifies dotted IDs survive round-trips.
func TestGTFSIdWithDots(t *testing.T) {
	ts := newTestServer(t)
	col := connectCollector(t, ts)

	sendRaw(t, col.ws, protocol.Envelope{Method: protocol.MethodSubscribe, Topic: protocol.TopicSubscriptions})

	cli := connectClient(t, ts)
	subscribeMap(t, cli, kielViewport)
	cli.drainSnapshots(t)

	const dottedID = "gtfs-agency:route.stop-1"
	stop := &models.Stop{
		ID: dottedID, Provider: "gtfs", Name: "GTFS Stop",
		Type: models.StopTypeBusStop, Location: kielLoc,
	}
	sendData(t, col, protocol.MethodUpdate, "stop."+dottedID, stop)

	e := cli.recv(t)
	assert.Equal(t, "map.stop."+dottedID, e.Topic, "dotted GTFS id must survive in topic")

	detailTopic := "stop." + dottedID
	sendSubscribe(t, cli, detailTopic)

	subEnv := col.recv(t)
	var se protocol.SubscriptionEvent
	require.NoError(t, json.Unmarshal(*subEnv.Data, &se))
	assert.Equal(t, detailTopic, se.Topic)
}

// TestVehicleDetailSubscription verifies vehicles work the same as stops.
func TestVehicleDetailSubscription(t *testing.T) {
	ts := newTestServer(t)
	col := connectCollector(t, ts)
	cli := connectClient(t, ts)

	subscribeMap(t, cli, kielViewport)
	cli.drainSnapshots(t)

	heading := 90
	v := &models.Vehicle{
		ID: "ais-99", Provider: "ais", Name: "Ferry X",
		Type: models.VehicleTypeFerry, State: "running",
		Location:    &models.Location{Latitude: kielLoc.Latitude, Longitude: kielLoc.Longitude, Heading: &heading},
		Description: "Full description only in detail",
	}
	sendData(t, col, protocol.MethodUpdate, "vehicle.ais-99", v)

	markerEnv := cli.recv(t)
	assert.Equal(t, "map.vehicle.ais-99", markerEnv.Topic)
	var markerVehicle models.Vehicle
	require.NoError(t, json.Unmarshal(*markerEnv.Data, &markerVehicle))
	assert.Empty(t, markerVehicle.Description, "description must not be in marker")

	sendSubscribe(t, cli, "vehicle.ais-99")

	detailEnv := cli.recv(t)
	assert.Equal(t, "vehicle.ais-99", detailEnv.Topic)
	var fullVehicle models.Vehicle
	require.NoError(t, json.Unmarshal(*detailEnv.Data, &fullVehicle))
	assert.Equal(t, "Full description only in detail", fullVehicle.Description)
}

// TestTripDetailSubscription verifies trips are routed to detail subscribers.
func TestTripDetailSubscription(t *testing.T) {
	ts := newTestServer(t)
	col := connectCollector(t, ts)
	cli := connectClient(t, ts)

	subscribeMap(t, cli, kielViewport)
	cli.drainSnapshots(t)

	sendSubscribe(t, cli, "trip.kvg-t1")
	time.Sleep(20 * time.Millisecond)

	trip := &models.Trip{
		ID: "kvg-t1", Provider: "kvg", Direction: "North",
		Departures: []*models.TripDeparture{{ID: "kvg-t1-stop1", Name: "A", State: models.Planned}},
	}
	sendData(t, col, protocol.MethodUpdate, "trip.kvg-t1", trip)

	e := cli.recv(t)
	assert.Equal(t, "trip.kvg-t1", e.Topic)
	var rt models.Trip
	require.NoError(t, json.Unmarshal(*e.Data, &rt))
	assert.Len(t, rt.Departures, 1)
}
