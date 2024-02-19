package api

import (
	"context"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/kiel-live/kiel-live/protocol"
	"golang.org/x/oauth2"
)

const token = ""
const endpoint = "https://graphql-sandbox-dds.rnv-online.de/"

const IDPrefix = "rnv:"

func getClient(ctx context.Context) graphql.Client {
	src := oauth2.StaticTokenSource(
		// &oauth2.Token{AccessToken: "Bearer: " + os.Getenv("GRAPHQL_TOKEN")},
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(ctx, src)

	client := graphql.NewClient(endpoint, httpClient)

	return client
}

func GetStops(_ctx context.Context) ([]protocol.Stop, error) {
	ctx, cancel := context.WithTimeout(_ctx, 3*time.Minute)
	defer cancel()

	client := getClient(ctx)

	stops := make([]protocol.Stop, 0)
	cursor := ""

	for {
		resp, err := getStationsPage(ctx, client, cursor)
		if err != nil {
			return nil, err
		}

		for _, element := range resp.Stations.Elements {
			if stop, ok := element.(*getStationsPageStationsSearchResultElementsStation); ok {
				stops = append(stops, protocol.Stop{
					ID:       stop.GlobalID,
					Name:     stop.LongName,
					Provider: "rnv",
					Location: protocol.Location{
						Longitude: int(stop.Location.Long),
						Latitude:  int(stop.Location.Lat),
					},
				})
			}
		}

		cursor = resp.Stations.Cursor
		if cursor == "" {
			break
		}
	}

	return stops, nil
}
