package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/kiel-live/kiel-live/collectors/motis/version"
)

const IDPrefix = "motis-"

const defaultBaseURL = "https://api.transitous.org"

func getBaseURL() string {
	if u := os.Getenv("MOTIS_URL"); u != "" {
		return u
	}
	return defaultBaseURL
}

func get(path string, params url.Values) ([]byte, error) {
	u := getBaseURL() + path
	if len(params) > 0 {
		u += "?" + params.Encode()
	}
	req, err := http.NewRequest(http.MethodGet, u, nil) //nolint:gosec
	if err != nil {
		return nil, err
	}

	userAgent := fmt.Sprintf("flott-live/%s (%s)", version.Version, version.Repo)

	req.Header.Set("User-Agent", userAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d for %s: %s", resp.StatusCode, u, string(body))
	}
	return body, nil
}

func getJSON(path string, params url.Values, v any) error {
	body, err := get(path, params)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}
