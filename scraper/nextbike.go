package main

import (
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"

	"github.com/kiel-live/kiel-live/shared/models"
)

func nextbike() ([]*models.Stop, error) {
	stops := make([]*models.Stop, 0)

	cityIDs := "195" // Mannheim

	resp, err := http.Get("https://api.nextbike.net/maps/nextbike-live.json?city=" + cityIDs)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	nextbikeResp := NextbikeResponse{}
	err = json.Unmarshal(data, &nextbikeResp)
	if err != nil {
		return nil, err
	}

	for _, country := range nextbikeResp.Countries {
		for _, city := range country.Cities {
			for _, place := range city.Places {
				vehicles := make([]*models.Vehicle, 0)

				for _, bike := range place.BikeList {
					vehicle := &models.Vehicle{
						ID:       fmt.Sprintf("nextbike-%s", bike.Number),
						Name:     bike.Number,
						Type:     "bike",
						Provider: "nextbike",
						State:    bike.State,
						Location: &models.Location{
							Latitude:  place.Lat,
							Longitude: place.Lng,
						},
					}
					vehicles = append(vehicles, vehicle)
				}

				slices.SortFunc(vehicles, func(a, b *models.Vehicle) int {
					return cmp.Compare(strings.ToLower(a.ID), strings.ToLower(b.ID))
				})

				stop := &models.Stop{
					ID:       fmt.Sprintf("nextbike-%d", place.UID),
					Provider: "nextbike",
					Name:     place.Name,
					Type:     "bike-stop",
					Location: &models.Location{
						Latitude:  place.Lat,
						Longitude: place.Lng,
					},
					Vehicles: vehicles,
				}

				stops = append(stops, stop)
			}
		}
	}

	return stops, nil
}
