package api

import (
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoundsFromQuery_Success(t *testing.T) {
	values := url.Values{}
	values.Set("north", "54.35")
	values.Set("east", "10.20")
	values.Set("south", "54.30")
	values.Set("west", "10.10")

	bounds, err := boundsFromQuery(values)

	assert.NoError(t, err)
	assert.NotNil(t, bounds)
	assert.Equal(t, 54.35, bounds.North)
	assert.Equal(t, 10.20, bounds.East)
	assert.Equal(t, 54.30, bounds.South)
	assert.Equal(t, 10.10, bounds.West)
}

func TestBoundsFromQuery_MissingNorth(t *testing.T) {
	values := url.Values{}
	values.Set("east", "10.20")
	values.Set("south", "54.30")
	values.Set("west", "10.10")

	_, err := boundsFromQuery(values)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "north")
}

func TestBoundsFromQuery_MissingEast(t *testing.T) {
	values := url.Values{}
	values.Set("north", "54.35")
	values.Set("south", "54.30")
	values.Set("west", "10.10")

	_, err := boundsFromQuery(values)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "east")
}

func TestBoundsFromQuery_MissingSouth(t *testing.T) {
	values := url.Values{}
	values.Set("north", "54.35")
	values.Set("east", "10.20")
	values.Set("west", "10.10")

	_, err := boundsFromQuery(values)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "south")
}

func TestBoundsFromQuery_MissingWest(t *testing.T) {
	values := url.Values{}
	values.Set("north", "54.35")
	values.Set("east", "10.20")
	values.Set("south", "54.30")

	_, err := boundsFromQuery(values)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "west")
}

func TestBoundsFromQuery_InvalidFloat(t *testing.T) {
	values := url.Values{}
	values.Set("north", "invalid")
	values.Set("east", "10.20")
	values.Set("south", "54.30")
	values.Set("west", "10.10")

	_, err := boundsFromQuery(values)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "north")
}

func TestRespondWithJSON(t *testing.T) {
	w := httptest.NewRecorder()

	data := map[string]string{
		"message": "Hello, World!",
	}

	respondWithJSON(w, 200, data)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Body.String(), "Hello, World!")
}

func TestRespondWithError(t *testing.T) {
	w := httptest.NewRecorder()

	respondWithError(w, 400, "Bad Request")

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Body.String(), "Bad Request")
	assert.Contains(t, w.Body.String(), "error")
}
