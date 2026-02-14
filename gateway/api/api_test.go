package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kiel-live/kiel-live/gateway/database"
	"github.com/kiel-live/kiel-live/gateway/hub"
	"github.com/kiel-live/kiel-live/gateway/search"
	"github.com/kiel-live/kiel-live/pkg/models"
	"github.com/stretchr/testify/assert"
)

const testToken = "test-secret-token"

func toDegreesInt(value float64) int {
	return int(value * 3600000.0)
}

func setupTestServer(t *testing.T) (*Server, *http.ServeMux) {
	db := database.NewMemoryDatabase()
	err := db.Open()
	assert.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	searchIndex := search.NewMemorySearch()
	h := hub.NewHub(db)
	go h.Run()

	mux := http.NewServeMux()
	server := NewAPIServer(db, searchIndex, h, mux, testToken)

	return server, mux
}

func makeRequest(t *testing.T, mux *http.ServeMux, method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonData, err := json.Marshal(body)
		assert.NoError(t, err)
		reqBody = bytes.NewBuffer(jsonData)
	} else {
		reqBody = bytes.NewBuffer([]byte{})
	}

	req := httptest.NewRequest(method, path, reqBody)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	return w
}

// Stop Tests

func TestHandleGetStop_NotFound(t *testing.T) {
	_, mux := setupTestServer(t)

	w := makeRequest(t, mux, "GET", "/stops/non-existent", nil, "")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestHandleGetStop_Success(t *testing.T) {
	server, mux := setupTestServer(t)

	// Create a stop first
	stop := &models.Stop{
		ID:       "stop-1",
		Name:     "Test Stop",
		Provider: "test",
		Location: &models.Location{
			Latitude:  toDegreesInt(54.32),
			Longitude: toDegreesInt(10.14),
		},
	}
	err := server.db.SetStop(context.Background(), stop)
	assert.NoError(t, err)

	w := makeRequest(t, mux, "GET", "/stops/stop-1", nil, "")

	assert.Equal(t, http.StatusOK, w.Code)

	var retrieved models.Stop
	err = json.NewDecoder(w.Body).Decode(&retrieved)
	assert.NoError(t, err)
	assert.Equal(t, "stop-1", retrieved.ID)
	assert.Equal(t, "Test Stop", retrieved.Name)
}

func TestHandleGetStops_WithBounds(t *testing.T) {
	server, mux := setupTestServer(t)

	// Create multiple stops
	stops := []*models.Stop{
		{
			ID:       "stop-1",
			Name:     "Stop 1",
			Provider: "test",
			Location: &models.Location{
				Latitude:  toDegreesInt(54.32),
				Longitude: toDegreesInt(10.14),
			},
		},
		{
			ID:       "stop-2",
			Name:     "Stop 2",
			Provider: "test",
			Location: &models.Location{
				Latitude:  toDegreesInt(54.33),
				Longitude: toDegreesInt(10.15),
			},
		},
	}

	for _, stop := range stops {
		err := server.db.SetStop(context.Background(), stop)
		assert.NoError(t, err)
	}

	w := makeRequest(t, mux, "GET", "/stops?north=54.35&east=10.20&south=54.30&west=10.10", nil, "")

	assert.Equal(t, http.StatusOK, w.Code)

	var retrieved []*models.Stop
	err := json.NewDecoder(w.Body).Decode(&retrieved)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(retrieved), 2)
}

func TestHandleUpdateStop_WithAuth(t *testing.T) {
	_, mux := setupTestServer(t)

	stop := &models.Stop{
		ID:       "stop-1",
		Name:     "Updated Stop",
		Provider: "kvg",
		Location: &models.Location{
			Latitude:  toDegreesInt(54.32),
			Longitude: toDegreesInt(10.14),
		},
	}

	w := makeRequest(t, mux, "PUT", "/stops/stop-1", stop, testToken)

	assert.Equal(t, http.StatusOK, w.Code)

	var retrieved models.Stop
	err := json.NewDecoder(w.Body).Decode(&retrieved)
	assert.NoError(t, err)
	assert.Equal(t, "stop-1", retrieved.ID)
	assert.Equal(t, "Updated Stop", retrieved.Name)
}

func TestHandleUpdateStop_WithoutAuth(t *testing.T) {
	_, mux := setupTestServer(t)

	stop := &models.Stop{
		ID:       "stop-1",
		Name:     "Updated Stop",
		Provider: "kvg",
		Location: &models.Location{
			Latitude:  toDegreesInt(54.32),
			Longitude: toDegreesInt(10.14),
		},
	}

	w := makeRequest(t, mux, "PUT", "/stops/stop-1", stop, "")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestHandleUpdateStop_IDMismatch(t *testing.T) {
	_, mux := setupTestServer(t)

	stop := &models.Stop{
		ID:       "stop-wrong",
		Name:     "Updated Stop",
		Provider: "kvg",
		Location: &models.Location{
			Latitude:  toDegreesInt(54.32),
			Longitude: toDegreesInt(10.14),
		},
	}

	w := makeRequest(t, mux, "PUT", "/stops/stop-1", stop, testToken)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleDeleteStop_WithAuth(t *testing.T) {
	server, mux := setupTestServer(t)

	// Create a stop first
	stop := &models.Stop{
		ID:       "stop-1",
		Name:     "Test Stop",
		Provider: "test",
		Location: &models.Location{
			Latitude:  toDegreesInt(54.32),
			Longitude: toDegreesInt(10.14),
		},
	}
	err := server.db.SetStop(context.Background(), stop)
	assert.NoError(t, err)

	w := makeRequest(t, mux, "DELETE", "/stops/stop-1", nil, testToken)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify it's deleted
	_, err = server.db.GetStop(context.Background(), "stop-1")
	assert.Error(t, err)
}

// Vehicle Tests

func TestHandleGetVehicle_Success(t *testing.T) {
	server, mux := setupTestServer(t)

	vehicle := &models.Vehicle{
		ID:       "vehicle-1",
		Name:     "Bus 11",
		Provider: "kvg",
		Type:     models.VehicleTypeBus,
		Location: &models.Location{
			Latitude:  toDegreesInt(54.32),
			Longitude: toDegreesInt(10.14),
		},
	}
	err := server.db.SetVehicle(context.Background(), vehicle)
	assert.NoError(t, err)

	w := makeRequest(t, mux, "GET", "/vehicles/vehicle-1", nil, "")

	assert.Equal(t, http.StatusOK, w.Code)

	var retrieved models.Vehicle
	err = json.NewDecoder(w.Body).Decode(&retrieved)
	assert.NoError(t, err)
	assert.Equal(t, "vehicle-1", retrieved.ID)
	assert.Equal(t, "Bus 11", retrieved.Name)
}

func TestHandleGetVehicles_WithBounds(t *testing.T) {
	server, mux := setupTestServer(t)

	vehicles := []*models.Vehicle{
		{
			ID:       "vehicle-1",
			Name:     "Bus 11",
			Provider: "kvg",
			Type:     models.VehicleTypeBus,
			Location: &models.Location{
				Latitude:  toDegreesInt(54.32),
				Longitude: toDegreesInt(10.14),
			},
		},
		{
			ID:       "vehicle-2",
			Name:     "Bus 12",
			Provider: "kvg",
			Type:     models.VehicleTypeBus,
			Location: &models.Location{
				Latitude:  toDegreesInt(54.33),
				Longitude: toDegreesInt(10.15),
			},
		},
	}

	for _, vehicle := range vehicles {
		err := server.db.SetVehicle(context.Background(), vehicle)
		assert.NoError(t, err)
	}

	w := makeRequest(t, mux, "GET", "/vehicles?north=54.35&east=10.20&south=54.30&west=10.10", nil, "")

	assert.Equal(t, http.StatusOK, w.Code)

	var retrieved []*models.Vehicle
	err := json.NewDecoder(w.Body).Decode(&retrieved)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(retrieved), 2)
}

func TestHandleUpdateVehicle_WithAuth(t *testing.T) {
	_, mux := setupTestServer(t)

	vehicle := &models.Vehicle{
		ID:       "vehicle-1",
		Name:     "Updated Bus",
		Provider: "kvg",
		Type:     models.VehicleTypeBus,
		Location: &models.Location{
			Latitude:  toDegreesInt(54.32),
			Longitude: toDegreesInt(10.14),
		},
	}

	w := makeRequest(t, mux, "PUT", "/vehicles/vehicle-1", vehicle, testToken)

	assert.Equal(t, http.StatusOK, w.Code)

	var retrieved models.Vehicle
	err := json.NewDecoder(w.Body).Decode(&retrieved)
	assert.NoError(t, err)
	assert.Equal(t, "vehicle-1", retrieved.ID)
	assert.Equal(t, "Updated Bus", retrieved.Name)
}

func TestHandleUpdateVehicle_LocationRequired(t *testing.T) {
	_, mux := setupTestServer(t)

	vehicle := &models.Vehicle{
		ID:       "vehicle-1",
		Name:     "Bus without location",
		Provider: "kvg",
		Type:     models.VehicleTypeBus,
		Location: nil, // No location
	}

	w := makeRequest(t, mux, "PUT", "/vehicles/vehicle-1", vehicle, testToken)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleDeleteVehicle_WithAuth(t *testing.T) {
	server, mux := setupTestServer(t)

	vehicle := &models.Vehicle{
		ID:       "vehicle-1",
		Name:     "Bus 11",
		Provider: "kvg",
		Type:     models.VehicleTypeBus,
		Location: &models.Location{
			Latitude:  toDegreesInt(54.32),
			Longitude: toDegreesInt(10.14),
		},
	}
	err := server.db.SetVehicle(context.Background(), vehicle)
	assert.NoError(t, err)

	w := makeRequest(t, mux, "DELETE", "/vehicles/vehicle-1", nil, testToken)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify it's deleted
	_, err = server.db.GetVehicle(context.Background(), "vehicle-1")
	assert.Error(t, err)
}

// Trip Tests

func TestHandleGetTrip_Success(t *testing.T) {
	server, mux := setupTestServer(t)

	trip := &models.Trip{
		ID:        "trip-1",
		Provider:  "kvg",
		VehicleID: "bus-123",
		Direction: "Innenstadt",
	}
	err := server.db.SetTrip(context.Background(), trip)
	assert.NoError(t, err)

	w := makeRequest(t, mux, "GET", "/trips/trip-1", nil, "")

	assert.Equal(t, http.StatusOK, w.Code)

	var retrieved models.Trip
	err = json.NewDecoder(w.Body).Decode(&retrieved)
	assert.NoError(t, err)
	assert.Equal(t, "trip-1", retrieved.ID)
	assert.Equal(t, "Innenstadt", retrieved.Direction)
}

func TestHandleUpdateTrip_WithAuth(t *testing.T) {
	_, mux := setupTestServer(t)

	trip := &models.Trip{
		ID:        "trip-1",
		Provider:  "kvg",
		VehicleID: "bus-123",
		Direction: "Updated Direction",
	}

	w := makeRequest(t, mux, "PUT", "/trips/trip-1", trip, testToken)

	assert.Equal(t, http.StatusOK, w.Code)

	var retrieved models.Trip
	err := json.NewDecoder(w.Body).Decode(&retrieved)
	assert.NoError(t, err)
	assert.Equal(t, "trip-1", retrieved.ID)
	assert.Equal(t, "Updated Direction", retrieved.Direction)
}

func TestHandleDeleteTrip_WithAuth(t *testing.T) {
	server, mux := setupTestServer(t)

	trip := &models.Trip{
		ID:        "trip-1",
		Provider:  "kvg",
		VehicleID: "bus-123",
		Direction: "Innenstadt",
	}
	err := server.db.SetTrip(context.Background(), trip)
	assert.NoError(t, err)

	w := makeRequest(t, mux, "DELETE", "/trips/trip-1", nil, testToken)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify it's deleted
	_, err = server.db.GetTrip(context.Background(), "trip-1")
	assert.Error(t, err)
}

// Authorization Tests

func TestWithAuth_ValidToken(t *testing.T) {
	_, mux := setupTestServer(t)

	vehicle := &models.Vehicle{
		ID:       "vehicle-1",
		Name:     "Test Vehicle",
		Provider: "test",
		Type:     models.VehicleTypeBus,
		Location: &models.Location{
			Latitude:  toDegreesInt(54.32),
			Longitude: toDegreesInt(10.14),
		},
	}

	w := makeRequest(t, mux, "PUT", "/vehicles/vehicle-1", vehicle, testToken)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestWithAuth_InvalidToken(t *testing.T) {
	_, mux := setupTestServer(t)

	vehicle := &models.Vehicle{
		ID:       "vehicle-1",
		Name:     "Test Vehicle",
		Provider: "test",
		Type:     models.VehicleTypeBus,
		Location: &models.Location{
			Latitude:  toDegreesInt(54.32),
			Longitude: toDegreesInt(10.14),
		},
	}

	w := makeRequest(t, mux, "PUT", "/vehicles/vehicle-1", vehicle, "wrong-token")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestWithAuth_MissingToken(t *testing.T) {
	_, mux := setupTestServer(t)

	vehicle := &models.Vehicle{
		ID:       "vehicle-1",
		Name:     "Test Vehicle",
		Provider: "test",
		Type:     models.VehicleTypeBus,
		Location: &models.Location{
			Latitude:  toDegreesInt(54.32),
			Longitude: toDegreesInt(10.14),
		},
	}

	w := makeRequest(t, mux, "PUT", "/vehicles/vehicle-1", vehicle, "")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// CORS Tests

func TestCORSHeader(t *testing.T) {
	_, mux := setupTestServer(t)

	w := makeRequest(t, mux, "GET", "/vehicles?north=54.35&east=10.20&south=54.30&west=10.10", nil, "")

	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
}
