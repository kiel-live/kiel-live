package rt

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	gtfsRt "github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

type GTFSRTCollector struct {
	url    string
	client *http.Client
}

func NewGTFSRTCollector(ctx context.Context) (*GTFSRTCollector, error) {
	url := os.Getenv("GTFS_RT_URL")
	if url == "" {
		log.Infoln("gtfs-rt disabled! Please provide an url with GTFS_RT_URL")
		return nil, nil
	}

	client, err := setupClient(ctx)
	if err != nil {
		return nil, err
	}

	return &GTFSRTCollector{
		url:    url,
		client: client,
	}, nil
}

func setupClient(ctx context.Context) (*http.Client, error) {
	switch os.Getenv("GTFS_RT_AUTH") {
	case "azure":
		credentials := &AzureClientCredentials{
			TenantID:     os.Getenv("GTFS_RT_AZURE_TENANT_ID"),
			ClientID:     os.Getenv("GTFS_RT_AZURE_CLIENT_ID"),
			ClientSecret: os.Getenv("GTFS_RT_AZURE_CLIENT_SECRET"),
			Resource:     os.Getenv("GTFS_RT_AZURE_RESOURCE"),
		}
		return credentials.GetClient(ctx)
	default:
		return http.DefaultClient, nil
	}
}

func (g *GTFSRTCollector) FetchTripUpdates() (*gtfsRt.FeedMessage, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/tripupdates", g.url), nil)
	if err != nil {
		return nil, err
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected status code")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	feed := gtfsRt.FeedMessage{}
	err = proto.Unmarshal(body, &feed)
	if err != nil {
		return nil, err
	}

	return &feed, nil
}
